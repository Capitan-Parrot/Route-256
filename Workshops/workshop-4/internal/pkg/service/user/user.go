package user

import (
	"fmt"
	"sync"

	"workshop.4.2/internal/model"
)

func New() *Server {
	return &Server{
		cache: map[model.ClientID]model.Client{
			model.ClientUserID: {
				ID:   model.ClientUserID,
				Name: model.ClientUserName,
			},
		},
	}
}

type Server struct {
	mu    sync.RWMutex
	cache map[model.ClientID]model.Client
}

func (s *Server) Get(id model.ClientID) (model.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.cache[id]
	if !ok {
		return model.Client{}, fmt.Errorf("cannot find user by id: [%d]", id)
	}

	return user, nil
}
