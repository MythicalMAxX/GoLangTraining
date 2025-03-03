package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"userservices/internal/models"

	"github.com/hashicorp/consul/api"
)

type PaymentClient struct {
	consulClient *api.Client
	baseURL      string
}

func NewPaymentClient() (*PaymentClient, error) {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %v", err)
	}

	baseURL, err := getServiceURL("payment-service")
	if err != nil {
		return nil, fmt.Errorf("failed to get payment service URL: %v", err)
	}

	return &PaymentClient{
		consulClient: client,
		baseURL:      baseURL,
	}, nil
}

func (c *PaymentClient) MakePayment(req *models.PaymentRequest) (*models.PaymentResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	log.Printf("Sending payment request to %s/payment", c.baseURL)
	resp, err := http.Post(c.baseURL+"/payment", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Payment request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Payment service returned non-201 status: %d, body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("payment service error: %s", string(body))
	}

	var response models.PaymentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Failed to decode response: %v. Response body: %s", err, string(body))
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

func (c *PaymentClient) GetPayment(id string) (*models.PaymentResponse, error) {
	url := fmt.Sprintf("%s/payment/%s", c.baseURL, id)
	log.Printf("Getting payment details from: %s", url)

	resp, err := http.Get(url)
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

	if resp.StatusCode != http.StatusOK {
		log.Printf("Payment service returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("payment not found")
		default:
			return nil, fmt.Errorf("payment service error: %s", string(body))
		}
	}

	var response models.PaymentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Failed to decode response: %v. Response body: %s", err, string(body))
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}
