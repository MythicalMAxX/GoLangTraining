package services

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"restaurantservice/internal/models"
	"restaurantservice/internal/repositories"
	"restaurantservice/pkg/messaging"
)

type RestaurantService struct {
	repo      *repositories.RestaurantRepository
	msgClient *messaging.RabbitMQClient
}

func NewRestaurantService(repo *repositories.RestaurantRepository, msgClient *messaging.RabbitMQClient) *RestaurantService {
	return &RestaurantService{
		repo:      repo,
		msgClient: msgClient,
	}
}

func (s *RestaurantService) CreateRestaurant(restaurant *models.Restaurant) error {
	err := s.repo.Create(restaurant)
	if err != nil {
		return err
	}

	// Notify about new restaurant
	s.msgClient.PublishMessage("restaurant.created", restaurant)
	return nil
}

func (s *RestaurantService) GetRestaurant(id uuid.UUID) (*models.Restaurant, error) {
	return s.repo.GetByID(id)
}

func (s *RestaurantService) UpdateRestaurant(restaurant *models.Restaurant) error {
	_, err := s.repo.GetByID(restaurant.ID)
	if err != nil {
		return errors.New("restaurant not found")
	}

	err = s.repo.Update(restaurant)
	if err != nil {
		return err
	}

	s.msgClient.PublishMessage("restaurant.updated", restaurant)
	return nil
}

func (s *RestaurantService) DeleteRestaurant(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *RestaurantService) ListRestaurants() ([]models.Restaurant, error) {
	return s.repo.List()
}

type Order struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	Items       []string  `json:"items"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type RestaurantAssignedEvent struct {
	OrderID      string    `json:"order_id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Status       string    `json:"status"`
	AssignedAt   time.Time `json:"assigned_at"`
}

func (s *RestaurantService) HandleNewOrder(orderData []byte) error {
	var order Order
	if err := json.Unmarshal(orderData, &order); err != nil {
		return err
	}

	// Find suitable restaurant
	restaurants, err := s.repo.List()
	if err != nil || len(restaurants) == 0 {
		return errors.New("no restaurants available")
	}

	// For demo, assign to first available restaurant
	assignedRestaurant := restaurants[0]

	// Create restaurant assigned event
	event := RestaurantAssignedEvent{
		OrderID:      order.OrderID,
		RestaurantID: assignedRestaurant.ID,
		Status:       "RESTAURANT_ASSIGNED",
		AssignedAt:   time.Now(),
	}

	// Publish restaurant assigned event
	err = s.msgClient.PublishRestaurantAssigned(event)
	if err != nil {
		return err
	}

	log.Printf("Restaurant %s assigned to order %s", assignedRestaurant.ID, order.OrderID)
	return nil
}
