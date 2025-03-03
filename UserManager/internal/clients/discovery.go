package clients

import (
    "fmt"
    "log"
    "github.com/hashicorp/consul/api"
)

func getServiceURL(serviceName string) (string, error) {
    config := api.DefaultConfig()
    client, err := api.NewClient(config)
    if err != nil {
        return "", fmt.Errorf("failed to create consul client: %v", err)
    }

    services, _, err := client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return "", fmt.Errorf("failed to get service from consul: %v", err)
    }

    if len(services) == 0 {
        return "", fmt.Errorf("no healthy instances found for service: %s", serviceName)
    }

    serviceURL := fmt.Sprintf("http://%s", services[0].Service.Address)
    if services[0].Service.Port != 0 {
        serviceURL = fmt.Sprintf("%s:%d", serviceURL, services[0].Service.Port)
    }

    log.Printf("Found %s at: %s", serviceName, serviceURL)
    return serviceURL, nil
}