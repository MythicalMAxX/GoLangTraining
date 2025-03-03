package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"userservices/internal/models"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

type OrderClient struct {
	consulClient *api.Client
	baseURL      string
}

func NewOrderClient() (*OrderClient, error) {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %v", err)
	}

	baseURL, err := getServiceURL("order-service")
	if err != nil {
		return nil, fmt.Errorf("failed to get order service URL: %v", err)
	}

	return &OrderClient{
		consulClient: client,
		baseURL:      baseURL,
	}, nil
}

func (c *OrderClient) CreateOrder(userID string, inventory *models.InventoryResponse, orderCount int) (*models.OrderResponse, error) {
	totalPrice := inventory.Price * float64(orderCount)

	request := struct {
		UserID        string  `json:"user_id"`
		InventoryID   string  `json:"inventory_id"`
		OrderCount    int     `json:"order_count"`
		Price         float64 `json:"price"`
		TotalPrice    float64 `json:"total_price"`
		ProductName   string  `json:"product_name"`
		PaymentStatus string  `json:"payment_status"`
		OrderStatus   string  `json:"order_status"`
	}{
		UserID:        userID,
		InventoryID:   inventory.ID.String(),
		OrderCount:    orderCount,
		Price:         inventory.Price,
		TotalPrice:    totalPrice,
		ProductName:   inventory.ProductName,
		PaymentStatus: "pending",
		OrderStatus:   "pending",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	log.Printf("Sending create order request to %s/orders with payload: %s", c.baseURL, string(jsonData))

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/orders", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	log.Printf("Received response from order service: %s", string(body))

	if resp.StatusCode != http.StatusCreated {
		var errorResp struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("failed to create order: %s - %s", errorResp.Error, errorResp.Message)
	}

	// Parse direct response without message wrapper
	var rawResponse struct {
		ID            string    `json:"id"`
		UserID        string    `json:"user_id"`
		InventoryID   string    `json:"inventory_id"`
		ProductName   string    `json:"product_name"`
		Price         float64   `json:"price"`
		TotalPrice    float64   `json:"total_price"`
		OrderCount    int       `json:"order_count"`
		PaymentStatus string    `json:"payment_status"`
		OrderStatus   string    `json:"order_status"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	if err := json.Unmarshal(body, &rawResponse); err != nil {
		log.Printf("Failed to decode response: %v. Response body: %s", err, string(body))
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Convert string UUIDs to UUID type
	orderID, err := uuid.Parse(rawResponse.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID in response: %v", err)
	}

	userUUID, err := uuid.Parse(rawResponse.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in response: %v", err)
	}

	inventoryUUID, err := uuid.Parse(rawResponse.InventoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid inventory ID in response: %v", err)
	}

	// Build response with proper types
	orderResponse := &models.OrderResponse{
		ID:            orderID,
		UserID:        userUUID,
		InventoryID:   inventoryUUID,
		ProductName:   rawResponse.ProductName,
		Price:         rawResponse.Price,
		TotalPrice:    rawResponse.TotalPrice,
		OrderCount:    rawResponse.OrderCount,
		PaymentStatus: rawResponse.PaymentStatus,
		OrderStatus:   rawResponse.OrderStatus,
		CreatedAt:     &rawResponse.CreatedAt,
		UpdatedAt:     &rawResponse.UpdatedAt,
	}

	log.Printf("Successfully created order with ID: %s", orderResponse.ID)
	return orderResponse, nil
}

func (c *OrderClient) GetOrder(id string) (*models.OrderResponse, error) {
	// Clean up the UUID string first
	cleanID := strings.ReplaceAll(strings.TrimSpace(id), " ", "")

	// Validate UUID format
	parsedUUID, err := uuid.Parse(cleanID)
	if err != nil {
		log.Printf("Invalid UUID format: %s (Error: %v)", id, err)
		return nil, fmt.Errorf("invalid order ID format: must be a valid UUID v4")
	}

	url := fmt.Sprintf("%s/api/v1/orders/%s", c.baseURL, parsedUUID.String())
	log.Printf("Sending GET request to: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	log.Printf("Received response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("order not found")
		case http.StatusBadRequest:
			return nil, fmt.Errorf("invalid order ID")
		default:
			return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
		}
	}

	var response models.OrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Failed to decode response: %v. Response body: %s", err, string(body))
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

func (c *OrderClient) CancelOrder(id string) error {
	// Clean up and validate UUID
	cleanID := strings.ReplaceAll(strings.TrimSpace(id), " ", "")

	// Validate UUID format
	parsedUUID, err := uuid.Parse(cleanID)
	if err != nil {
		log.Printf("Invalid UUID format: %s (Error: %v)", id, err)
		return fmt.Errorf("invalid order ID format: must be a valid UUID v4")
	}

	url := fmt.Sprintf("%s/api/v1/orders/%s/cancel", c.baseURL, parsedUUID.String())
	log.Printf("Sending PUT request to: %s", url)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Read response body for better error handling
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Server returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
		switch resp.StatusCode {
		case http.StatusNotFound:
			return fmt.Errorf("order not found")
		case http.StatusBadRequest:
			return fmt.Errorf("invalid order ID")
		case http.StatusForbidden:
			return fmt.Errorf("not authorized to cancel this order")
		default:
			return fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(body))
		}
	}

	log.Printf("Successfully cancelled order %s", parsedUUID.String())
	return nil
}
