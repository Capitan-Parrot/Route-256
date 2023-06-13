package driver

import (
	"fmt"
	"sync"

	"workshop.4.2/internal/model"
)

func New() *Server {
	return &Server{
		cache: map[model.DriverID]model.Driver{
			model.DriverUserID: {
				ID:   model.DriverUserID,
				Name: model.DriverUserName,
				Car:  model.DriverUserCar,
			},
		},
	}
}

type Server struct {
	mu    sync.RWMutex
	cache map[model.DriverID]model.Driver
}

func (s *Server) Get(id model.DriverID) (model.Driver, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.cache[id]
	if !ok {
		return model.Driver{}, fmt.Errorf("cannot find user by id: [%d]", id)
	}

	return user, nil
}

func (s *Server) Free() (model.Driver, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.cache[model.DriverUserID], nil
}
