package storage

import (
	"errors"
	"sync"

	"orderservice/models"
)

type OrderStorage struct {
	orders map[string]models.Order
	mutex  sync.RWMutex
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]models.Order),
	}
}

func (s *OrderStorage) Create(order models.Order) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.orders[order.ID]; exists {
		return errors.New("order already exists")
	}

	s.orders[order.ID] = order
	return nil
}

func (s *OrderStorage) Get(id string) (models.Order, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if order, exists := s.orders[id]; exists {
		return order, nil
	}
	return models.Order{}, errors.New("order not found")
}

func (s *OrderStorage) Update(order models.Order) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.orders[order.ID]; !exists {
		return errors.New("order not found")
	}

	s.orders[order.ID] = order
	return nil
}

func (s *OrderStorage) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.orders[id]; !exists {
		return errors.New("order not found")
	}

	delete(s.orders, id)
	return nil
}

func (s *OrderStorage) List() []models.Order {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	orders := make([]models.Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders
}
