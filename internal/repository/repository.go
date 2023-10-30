package repository

import (
	"inttest-runtime/internal/repository/service_repo"
	"sync"
)

type Repository struct {
	mu *sync.Mutex

	*service_repo.ServiceRepository
}

func NewRepository() *Repository {
	return &Repository{
		mu:                &sync.Mutex{},
		ServiceRepository: service_repo.NewServiceRepository(),
	}
}

func (r *Repository) WithLock(fn func(r IRepository) error) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return fn(r)
}

type IRepository interface {
	service_repo.IServiceRepository
}

type Transactioner interface {
	WithLock(fn func(r IRepository) error) error
}
