package main

import (
    "fmt"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "strings"
    "time"

    "github.com/hashicorp/consul/api"
)

// APIConfig holds API configuration
type APIConfig struct {
    Version     string
    BaseURL     string
    ServicePort int
}

type ServiceRouter struct {
    consulClient *api.Client
    config       *APIConfig
}

// Custom error types
type RouterError struct {
    Code    int
    Message string
}

func (e *RouterError) Error() string {
    return e.Message
}

func NewServiceRouter(config *APIConfig) (*ServiceRouter, error) {
    consulConfig := api.DefaultConfig()
    client, err := api.NewClient(consulConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create consul client: %w", err)
    }
    return &ServiceRouter{
        consulClient: client,
        config:      config,
    }, nil
}

func (sr *ServiceRouter) getServiceAddress(serviceName string) (string, error) {
    services, _, err := sr.consulClient.Health().Service(serviceName, "", true, &api.QueryOptions{
        WaitTime: time.Second * 5,
    })
    if err != nil {
        return "", &RouterError{
            Code:    http.StatusServiceUnavailable,
            Message: fmt.Sprintf("service lookup failed: %v", err),
        }
    }
    if len(services) == 0 {
        return "", &RouterError{
            Code:    http.StatusNotFound,
            Message: fmt.Sprintf("service '%s' not found", serviceName),
        }
    }
    
    // Basic load balancing - round robin could be implemented here
    service := services[0].Service
    return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}

func (sr *ServiceRouter) handleRequest(w http.ResponseWriter, r *http.Request) {
    // Add common headers
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "DENY")
    w.Header().Set("X-XSS-Protection", "1; mode=block")

    // Parse path: /v1/service-name/actual/path
    parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 3)
    if len(parts) < 2 {
        http.Error(w, "invalid path", http.StatusBadRequest)
        return
    }

    // Validate API version
    version := parts[0]
    if version != sr.config.Version {
        http.Error(w, "unsupported API version", http.StatusNotFound)
        return
    }

    serviceName := parts[1]
    serviceURL, err := sr.getServiceAddress(serviceName)
    if err != nil {
        if routerErr, ok := err.(*RouterError); ok {
            http.Error(w, routerErr.Message, routerErr.Code)
            return
        }
        http.Error(w, "internal server error", http.StatusInternalServerError)
        return
    }

    targetURL, err := url.Parse(serviceURL)
    if err != nil {
        http.Error(w, "invalid service URL", http.StatusInternalServerError)
        return
    }

    // Create reverse proxy
    proxy := httputil.NewSingleHostReverseProxy(targetURL)
    
    // Modify request path to remove version and service name
    if len(parts) == 3 {
        r.URL.Path = "/" + parts[2]
    } else {
        r.URL.Path = "/"
    }

    // Add headers for debugging and tracing
    r.Header.Set("X-Forwarded-Host", r.Host)
    r.Header.Set("X-Origin-Host", targetURL.Host)
    
    proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
        log.Printf("proxy error: %v", err)
        http.Error(w, "service unavailable", http.StatusServiceUnavailable)
    }

    proxy.ServeHTTP(w, r)
}

// Display all the servces in ServiceRegistry
func (sr *ServiceRouter) listRegisteredServices() {
    services, _, err := sr.consulClient.Catalog().Services(&api.QueryOptions{})
    if err != nil {
        log.Printf("Error fetching services: %v", err)
        return
    }

    log.Printf("\nRegistered Services in Consul:")
    log.Printf("==============================")
    for serviceName, tags := range services {
        health, _, err := sr.consulClient.Health().Service(serviceName, "", true, nil)
        if err != nil {
            continue
        }

        for _, entry := range health {
            log.Printf("Service: %s", serviceName)
            log.Printf("  Tags: %v", tags)
            log.Printf("  ID: %s", entry.Service.ID)
            log.Printf("  Address: %s:%d", entry.Service.Address, entry.Service.Port)
            log.Printf("  Health: %s", entry.Checks.AggregatedStatus())
            log.Printf("------------------------------")
        }
    }
}

func main() {
    config := &APIConfig{
        Version:     "v1",
        BaseURL:     "/api",
        ServicePort: 8080,
    }

    router, err := NewServiceRouter(config)
    if err != nil {
        log.Fatal(err)
    }

    // Display Consul connection info
    consulConfig := api.DefaultConfig()
    log.Printf("Consul Connection Info:")
    log.Printf("  Address: %s", consulConfig.Address) // Default: http://127.0.0.1:8500
    log.Printf("  Datacenter: %s", consulConfig.Datacenter)
    log.Printf("  Scheme: %s", consulConfig.Scheme)

    // List registered services
    router.listRegisteredServices()

    server := &http.Server{
        Addr:         fmt.Sprintf(":%d", config.ServicePort),
        Handler:      http.HandlerFunc(router.handleRequest),
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Printf("\nService Router starting on %s...", server.Addr)
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}