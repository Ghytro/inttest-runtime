package domain

import (
	"inttest-runtime/pkg/utils"
	"time"
)

type ServiceID string

type baseDomainEntity struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (b baseDomainEntity) PtrCopy() *baseDomainEntity {
	return &baseDomainEntity{
		CreatedAt: b.CreatedAt,
		UpdatedAt: utils.ToPtr(*b.UpdatedAt),
		DeletedAt: utils.ToPtr(*b.DeletedAt),
	}
}

type RpcService struct {
	baseDomainEntity

	ID         ServiceID `json:"id"`
	Type       RpcType   `json:"type"`
	RemoteAddr string    `json:"remote_addr"`
}

func (s RpcService) PtrCopy() *RpcService {
	return &RpcService{
		baseDomainEntity: *s.baseDomainEntity.PtrCopy(),
		ID:               s.ID,
		Type:             s.Type,
		RemoteAddr:       s.RemoteAddr,
	}
}

type RpcType string

const (
	RpcTypeRestApi RpcType = "rest"
	RpcTypeGrpc    RpcType = "grpc"
)
