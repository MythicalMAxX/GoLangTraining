package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"
    "userservices/config"
    "userservices/internal/handlers"
    "userservices/internal/repositories"
    "userservices/internal/services"
    "userservices/routes"

    "github.com/gorilla/mux"
    "github.com/hashicorp/consul/api"
    "github.com/joho/godotenv"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func registerService(serviceName string, port int) error {
    config := api.DefaultConfig()
    client, err := api.NewClient(config)
    if err != nil {
        return err
    }

    registration := &api.AgentServiceRegistration{
        ID:      fmt.Sprintf("%s-%d", serviceName, port),
        Name:    serviceName,
        Port:    port,
        Address: "localhost",
        Tags:    []string{"v1"},
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://localhost:%d/health", port),
            Interval: "10s",
            Timeout:  "5s",
        },
    }

    return client.Agent().ServiceRegister(registration)
}

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Initialize database
    db, err := config.InitDB()
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

        // Initialize repositories, services, and handlers
    userRepo := repositories.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    userHandler, err := handlers.NewUserHandler(userService)
    if err != nil {
        log.Fatal("Failed to create user handler:", err)
    }

    // Get environment variables
    serviceName := os.Getenv("SERVICE_NAME")
    servicePort, _ := strconv.Atoi(os.Getenv("SERVICE_PORT"))

    // Register service with Consul
    if err := registerService(serviceName, servicePort); err != nil {
        log.Fatal("Failed to register service:", err)
    }

    router := mux.NewRouter()

    // Apply middleware
    router.Use(loggingMiddleware)

    // Setup routes
    routes.SetupRoutes(router, userHandler)

    // Health check endpoint
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    // Start server
    address := fmt.Sprintf(":%d", servicePort)
    log.Printf("Starting %s on port %d", serviceName, servicePort)
    log.Fatal(http.ListenAndServe(address, router))
}