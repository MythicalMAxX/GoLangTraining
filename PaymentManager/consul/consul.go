package consul

import (
	"fmt"
	"log"

	consul "github.com/hashicorp/consul/api"
)

func RegisterService(serviceID string, serviceName string, servicePort string) (*consul.Client, error) {
	config := consul.DefaultConfig()
	config.Address = "localhost:8500"

	client, err := consul.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %v", err)
	}

	registration := &consul.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Port:    8085,
		Address: "localhost",
		Check: &consul.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://localhost:%s/health", servicePort),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		return nil, fmt.Errorf("failed to register service: %v", err)
	}

	log.Printf("Successfully registered service: %s", serviceID)
	return client, nil
}

func DeregisterService(client *consul.Client, serviceID string) error {
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		return fmt.Errorf("failed to deregister service: %v", err)
	}

	log.Printf("Successfully deregistered service: %s", serviceID)
	return nil
}
