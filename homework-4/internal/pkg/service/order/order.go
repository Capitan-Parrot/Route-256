package order

import (
	"fmt"
	"sync"

	"homework/internal/model"
)

func New() *Server {
	return &Server{
		cache: map[model.OrderID]model.Order{},
	}
}

// Server для хранения и обработки заказов
type Server struct {
	mu    sync.RWMutex
	cache map[model.OrderID]model.Order
}

func (s *Server) Get(orderID model.OrderID) (model.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.cache[orderID]
	if !ok {
		return model.Order{}, fmt.Errorf("cannot find order by id: [%d]", orderID)
	}

	return order, nil
}

func (s *Server) Create(order model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[order.ID] = order
	return nil
}

func (s *Server) Update(order model.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.cache[order.ID]
	if !ok {
		return fmt.Errorf("cannot find order by id: [%d]", order.ID)
	}
	s.cache[order.ID] = order
	return nil
}
