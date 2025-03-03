package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
    "log"
    "io"
	"net/http"
	"net/url"
	"userservices/internal/models"

	"github.com/hashicorp/consul/api"
)


type InventoryClient struct {
	consulClient *api.Client
	baseURL      string
}

func NewInventoryClient() (*InventoryClient, error) {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %v", err)
	}

	baseURL, err := getServiceURL("inventory-service")
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory service URL: %v", err)
	}

	return &InventoryClient{
		consulClient: client,
		baseURL:      baseURL,
	}, nil
}

func (c *InventoryClient) CreateInventory(req *models.InventoryRequest) (*models.InventoryResponse, error) {
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %v", err)
    }

    resp, err := http.Post(c.baseURL+"/api/v1/inventory", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Read the body first
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %v", err)
    }

    if resp.StatusCode != http.StatusCreated {
        // Log the actual response for debugging
        log.Printf("Inventory service error response: %s", string(body))
        
        // Try to parse error response in different formats
        var errorResp struct {
            Error   interface{} `json:"error"`
            Message string      `json:"message"`
            Details interface{} `json:"details"`
        }
        
        if err := json.Unmarshal(body, &errorResp); err != nil {
            // If can't parse, return the raw response
            return nil, fmt.Errorf("inventory service error: %s", string(body))
        }
        
        // Return formatted error message
        return nil, fmt.Errorf("inventory service error: %v - %s", errorResp.Error, errorResp.Message)
    }

    var response models.InventoryResponse
    if err := json.Unmarshal(body, &response); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    return &response, nil
}

func (c *InventoryClient) GetInventoryByID(id string) (*models.InventoryResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/inventory/%s", c.baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("inventory not found")
	}

	var response models.InventoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

func (c *InventoryClient) GetAllInventory(page, pageSize int, status, productName string) (*models.InventoryListResponse, error) {
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", page))
	query.Set("page_size", fmt.Sprintf("%d", pageSize))
	if status != "" {
		query.Set("status", status)
	}
	if productName != "" {
		query.Set("product_name", productName)
	}

	url := fmt.Sprintf("%s/api/v1/inventory?%s", c.baseURL, query.Encode())
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response models.InventoryListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

func (c *InventoryClient) UpdateInventory(id string, updates map[string]interface{}) error {
	jsonData, err := json.Marshal(updates)
	if err != nil {
		return fmt.Errorf("failed to marshal updates: %v", err)
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/v1/inventory/%s", c.baseURL, id), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error   string `json:"error"`
			Details string `json:"details"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return fmt.Errorf("failed to decode error response: %v", err)
		}
		return fmt.Errorf("%s: %s", errorResp.Error, errorResp.Details)
	}

	return nil
}
