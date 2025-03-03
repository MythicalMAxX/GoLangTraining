package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"gorm.io/gorm"

	"orderservice/config"
	"orderservice/models"
)

type OrderMessage struct {
	Type      string       `json:"type"`
	Order     models.Order `json:"order"`
	Timestamp time.Time    `json:"timestamp"`
}

type OrderHandler struct {
	db      *gorm.DB
	channel *amqp.Channel
}

func NewOrderHandler(db *gorm.DB, ch *amqp.Channel) *OrderHandler {
	return &OrderHandler{
		db:      db,
		channel: ch,
	}
}

func (h *OrderHandler) publishMessage(msgType string, order models.Order) error {
	message := OrderMessage{
		Type:      msgType,
		Order:     order,
		Timestamp: time.Now(),
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return h.channel.Publish(
		config.OrdersExchange,
		msgType,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if order.RestaurantID == "" {
		http.Error(w, "restaurant_id is required", http.StatusBadRequest)
		return
	}

	if order.CustomerID == "" {
		http.Error(w, "customer_id is required", http.StatusBadRequest)
		return
	}

	// Set order metadata
	order.ID = uuid.New().String()
	order.Status = "PENDING"

	// Set IDs for order items
	for i := range order.Items {
		order.Items[i].ID = uuid.New().String()
		order.Items[i].OrderID = order.ID
	}

	// Calculate total amount
	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.TotalAmount = total

	// Store in database
	if err := h.db.Create(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Publish to RabbitMQ
	if err := h.publishMessage(config.MessageTypeNewOrder, order); err != nil {
		// Log the error but don't fail the request
		log.Printf("Failed to publish order: %v", err)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var order models.Order
	if err := h.db.Preload("Items").First(&order, "id = ?", vars["id"]).Error; err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.ID = vars["id"]
	if err := h.db.Save(&order).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := h.db.Delete(&models.Order{}, "id = ?", vars["id"]).Error; err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	var orders []models.Order
	if err := h.db.Preload("Items").Find(&orders).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}

// Add new method to get orders by restaurant
func (h *OrderHandler) GetOrdersByRestaurant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restaurantID := vars["restaurant_id"]

	var orders []models.Order
	if err := h.db.Preload("Items").Where("restaurant_id = ?", restaurantID).Find(&orders).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
