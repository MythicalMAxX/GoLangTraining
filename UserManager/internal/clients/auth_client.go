package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	// "os"

	"github.com/hashicorp/consul/api"
)

type AuthClient struct {
	consulClient *api.Client
}

type TokenRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ValidationResponse struct {
	IsValid   int    `json:"isValid"` // Direct fields, not nested
	UUID      string `json:"uuid"`
	Role      string `json:"role"`
	Action    string `json:"action"`
	RequestID string `json:"request_id"`
}

type ValidationBody struct {
	IsValid   int    `json:"isValid"`
	UUID      string `json:"uuid"`
	Role      string `json:"role"`
	Action    string `json:"action"`
	RequestID string `json:"request_id"`
}

func NewAuthClient() (*AuthClient, error) {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &AuthClient{consulClient: client}, nil
}

func (c *AuthClient) getAuthServiceURL() (string, error) {
	// Try to get service from Consul with correct service name
	services, _, err := c.consulClient.Health().Service("auth-manager", "", true, nil)
	if err != nil {
		log.Printf("Consul error: Failed to get auth service: %v", err)
		return "", err
	}
	if len(services) == 0 {
		log.Printf("Consul error: No healthy auth service instances found")
		return "", fmt.Errorf("auth service not found")
	}
	serviceURL := fmt.Sprintf("http://%s:%d", services[0].Service.Address, services[0].Service.Port)
	log.Printf("Found auth service at: %s", serviceURL)
	return serviceURL, nil
}

func (c *AuthClient) GenerateToken(userID, role string) (string, error) {
	serviceURL, err := c.getAuthServiceURL()
	if err != nil {
		return "", err
	}

	request := struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
		Action string `json:"action"`
		Token  string `json:"token"`
	}{
		UserID: userID,
		Role:   role,
		Action: "generate",
		Token:  "",
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(serviceURL+"/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response struct {
		Token     string `json:"token"`
		Action    string `json:"action"`
		RequestID string `json:"request_id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Token, nil
}

func (c *AuthClient) ValidateToken(token string) (*ValidationResponse, error) {
	serviceURL, err := c.getAuthServiceURL()
	if err != nil {
		log.Printf("Failed to get auth service URL: %v", err)
		return nil, err
	}

	request := struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
		Action string `json:"action"`
		Token  string `json:"token"`
	}{
		UserID: "",
		Role:   "",
		Action: "validate",
		Token:  token,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal token validation request: %v", err)
		return nil, err
	}

	log.Printf("Sending token validation request to %s/validate", serviceURL)
	resp, err := http.Post(serviceURL+"/validate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body for logging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}
	log.Printf("Received response from auth service: %s", string(body))

	var validationResp ValidationResponse
	if err := json.Unmarshal(body, &validationResp); err != nil {
		log.Printf("Failed to decode validation response: %v", err)
		return nil, err
	}

	return &validationResp, nil
}
