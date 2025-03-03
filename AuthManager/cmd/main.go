package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"authmanager/internal/handlers"

	"github.com/hashicorp/consul/api"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func registerService() {
	log.Printf("Configuring Consul client...")
	config := api.DefaultConfig()
	config.Address = "localhost:8500" // Explicitly set Consul address

	consul, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	registration := &api.AgentServiceRegistration{
		ID:      "auth-manager-1", // Unique ID for this instance
		Name:    "auth-manager",
		Port:    8082,
		Address: "localhost",
		Tags:    []string{"auth", "jwt"},
		Check: &api.AgentServiceCheck{
			HTTP:                           "http://localhost:8082/health",
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	log.Printf("Attempting to register service with Consul...")
	if err := consul.Agent().ServiceRegister(registration); err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}
	log.Printf("Successfully registered service with Consul")

	// Verify registration
	services, err := consul.Agent().Services()
	if err != nil {
		log.Printf("WARNING: Could not verify service registration: %v", err)
		return
	}

	if _, exists := services["auth-manager-1"]; exists {
		log.Printf("Verified service registration in Consul")
	} else {
		log.Printf("WARNING: Service registration not found in Consul")
	}
}

func main() {
	log.Printf("Starting Auth Manager service...")
	log.Printf("Loading environment variables...")

	// Check required environment variables
	requiredEnvVars := []string{"SECRET_KEY", "SERVICE_PORT"}
	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			log.Fatalf("Required environment variable %s is not set", env)
		}
		log.Printf("Environment variable %s is properly set", env)
	}

	// Register service with Consul
	log.Printf("Registering service with Consul...")
	registerService()
	log.Printf("Successfully registered with Consul")

	// Initialize handlers
	handler := handlers.NewAuthHandler()
	log.Printf("Handlers initialized")

	// Setup routes
	http.HandleFunc("/generate", handler.GenerateToken)
	http.HandleFunc("/validate", handler.ValidateToken)
	// Setup routes with enhanced health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Health check received from %s", r.RemoteAddr)

		// Check if required env vars are set
		secretKey := os.Getenv("SECRET_KEY")
		if secretKey == "" {
			log.Printf("ERROR: Health check failed - SECRET_KEY not set")
			http.Error(w, "Service unhealthy - configuration missing", http.StatusServiceUnavailable)
			return
		}

		// Return detailed health status
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "auth-manager",
			"version": "1.0",
		})
		log.Printf("Health check successful")
	})
	log.Printf("Routes configured")

	// Setup graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan
		log.Printf("Received shutdown signal: %v", sig)

		// Deregister service from Consul
		config := api.DefaultConfig()
		consul, err := api.NewClient(config)
		if err != nil {
			log.Printf("ERROR: Failed to create Consul client for deregistration: %v", err)
			return
		}

		if err := consul.Agent().ServiceDeregister("auth-manager-1"); err != nil {
			log.Printf("ERROR: Failed to deregister service from Consul: %v", err)
		} else {
			log.Printf("Successfully deregistered service from Consul")
		}

		os.Exit(0)
	}()

	// Start HTTP server
	port := os.Getenv("SERVICE_PORT")
	log.Printf("Starting HTTP server on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
