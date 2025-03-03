package discovery

import (
    "fmt"
    "github.com/hashicorp/consul/api"
)

type ServiceDiscovery struct {
    client *api.Client
}

func NewServiceDiscovery(address string) (*ServiceDiscovery, error) {
    config := api.DefaultConfig()
    config.Address = address

    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }

    return &ServiceDiscovery{client: client}, nil
}

func (sd *ServiceDiscovery) Register(name string, id string, address string, port int) error {
    registration := &api.AgentServiceRegistration{
        ID:      id,
        Name:    name,
        Address: address,
        Port:    port,
        Check: &api.AgentServiceCheck{
            HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
            Interval: "10s",
            Timeout:  "5s",
        },
    }

    return sd.client.Agent().ServiceRegister(registration)
}

func (sd *ServiceDiscovery) GetService(name string) (string, error) {
    services, _, err := sd.client.Health().Service(name, "", true, nil)
    if err != nil {
        return "", err
    }

    if len(services) == 0 {
        return "", fmt.Errorf("service %s not found", name)
    }

    // Get the first healthy instance
    service := services[0].Service
    return fmt.Sprintf("http://%s:%d", service.Address, service.Port), nil
}