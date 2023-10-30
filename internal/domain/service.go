package domain

import (
	"inttest-runtime/pkg/utils"
	"time"
)

type ServiceID string

func (id ServiceID) IsEmpty() bool {
	return id == ""
}

func ParseServiceID(sID string) (ServiceID, error) {
	return ServiceID(sID), nil
}

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

	CurrentStatus ServiceUptimeStatus `json:"status"`
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

func (t RpcType) String() string {
	return string(t)
}

const (
	RpcTypeRestApi RpcType = "rest"
	RpcTypeGrpc    RpcType = "grpc"
)

type ServiceUptimeStatus string

func (s ServiceUptimeStatus) String() string {
	return string(s)
}

const (
	ServiceUptimeStatus_DOWN     ServiceUptimeStatus = "Down"
	ServiceUptimeStatus_UP       ServiceUptimeStatus = "Up"
	ServiceUptimeStatus_STARTING ServiceUptimeStatus = "Starting"
	ServiceUptimeStatus_SHUTTING ServiceUptimeStatus = "Shutting"
)
