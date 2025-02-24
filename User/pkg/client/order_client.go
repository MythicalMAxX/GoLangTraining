package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type OrderClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewOrderClient() *OrderClient {
	return &OrderClient{
		baseURL:    os.Getenv("ORDER_SERVICE_URL"),
		httpClient: &http.Client{},
	}
}

// Update Order struct in pkg/client/order_client.go
type Order struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	InventoryID string  `json:"inventory_id"`
	Quantity    int     `json:"quantity"`
	Amount      float64 `json:"amount"` // Added Amount
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"` // Added UpdatedAt
}

type OrderResponse struct {
	Orders []Order `json:"orders"`
}

func (c *OrderClient) GetUserOrders(userID string) ([]Order, error) {
	// Create request
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/v1/users/%s/orders", c.baseURL, userID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Add Authorization header with user token
	token := c.GetAuthToken(userID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to get orders, status: %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("failed to get orders: %s", errorResponse.Error)
	}

	// Try to decode as array directly first
	var orders []Order
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return orders, nil
}

// Update the CreateOrderRequest struct in pkg/client/order_client.go
type CreateOrderRequest struct {
	InventoryID string  `json:"inventory_id"` // Changed from ProductID
	Quantity    int     `json:"quantity"`
	Amount      float64 `json:"amount"` // Added Amount field
	UserID      string  `json:"user_id,omitempty"`
}

// Update the CreateOrder method to include authorization for user_id
func (c *OrderClient) CreateOrder(order CreateOrderRequest) error {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/orders", c.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	// Add Authorization header with user token including user ID
	token := c.GetAuthToken(order.UserID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// Read error response
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to create order, status: %d", resp.StatusCode)
		}
		return fmt.Errorf("failed to create order: %s", errorResponse.Error)
	}

	return nil
}

// Update CancelOrder in OrderClient
func (c *OrderClient) CancelOrder(orderID string, userID string) error {
	// Create request
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/v1/orders/%s/cancel", c.baseURL, orderID),
		nil,
	)
	if err != nil {
		return err
	}

	// Add Authorization header with user token
	token := c.GetAuthToken(userID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to cancel order, status: %d", resp.StatusCode)
		}
		return fmt.Errorf("failed to cancel order: %s", errorResponse.Error)
	}

	return nil
}

// Add this to pkg/client/order_client.go
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Last function should be GetAuthToken
func (c *OrderClient) GetAuthToken(userID string) string {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "user-service",
		},
	})

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return ""
	}
	return token
}

type UserHandler struct {
	orderClient *OrderClient
}
