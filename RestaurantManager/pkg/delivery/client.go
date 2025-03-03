package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DeliveryClient struct {
	baseURL    string
	httpClient *http.Client
}

type DeliveryRequest struct {
	OrderID         string    `json:"order_id"` // Changed from ID to OrderID
	CustomerID      string    `json:"customer_id"`
	RestaurantID    string    `json:"restaurant_id"`
	Status          string    `json:"status"`
	TotalAmount     float64   `json:"total_amount"`
	DeliveryAddress string    `json:"delivery_address"`
	Items           []string  `json:"items"` // Added items field
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

func NewDeliveryClient(port int) *DeliveryClient {
	return &DeliveryClient{
		baseURL: fmt.Sprintf("http://localhost:%d", port),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *DeliveryClient) RequestDelivery(req *DeliveryRequest) error {
	// Add validation
	if req.OrderID == "" {
		return fmt.Errorf("order ID cannot be empty")
	}
	if req.CustomerID == "" {
		return fmt.Errorf("customer ID cannot be empty")
	}
	if req.RestaurantID == "" {
		return fmt.Errorf("restaurant ID cannot be empty")
	}

	log.Printf("Sending delivery request to delivery service: %+v", req)
	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery request: %w", err)
	}

	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/api/v1/delivery/request", c.baseURL),
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return fmt.Errorf("failed to send delivery request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("delivery service returned non-success status: %d", resp.StatusCode)
	}

	return nil
}
