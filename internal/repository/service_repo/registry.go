package service_repo

import (
	"errors"
	"inttest-runtime/internal/domain"
	"sync"
)

type ServiceRegistry struct {
	byId map[domain.ServiceID]*domain.RpcService
	mu   *sync.Mutex
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		mu:   &sync.Mutex{},
		byId: map[domain.ServiceID]*domain.RpcService{},
	}
}

func (r *ServiceRegistry) WithLock(fn func(r *ServiceRegistry) error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return fn(r)
}

func (r ServiceRegistry) Register(s *domain.RpcService) error {
	if _, ok := r.byId[s.ID]; ok {
		return errors.New("service non-unque id registered")
	}
	r.byId[s.ID] = s.PtrCopy()
	return nil
}

func (r ServiceRegistry) Get(id domain.ServiceID) (*domain.RpcService, error) {
	res, ok := r.byId[id]
	if !ok {
		return nil, errors.New("no session found")
	}
	return res.PtrCopy(), nil
}
