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

type ServiceUptimeStatus int

func (s ServiceUptimeStatus) String() string {
	if s < 0 || int(s) >= len(serviceStatusStrings) {
		return ""
	}
	return serviceStatusStrings[s]
}

const (
	ServiceUptimeStatus_DOWN ServiceUptimeStatus = iota
	ServiceUptimeStatus_UP
	ServiceUptimeStatus_STARTING
	ServiceUptimeStatus_SHUTTING
)

var serviceStatusStrings = [...]string{
	ServiceUptimeStatus_DOWN:     "Down",
	ServiceUptimeStatus_UP:       "Up",
	ServiceUptimeStatus_STARTING: "Starting",
	ServiceUptimeStatus_SHUTTING: "Shutting",
}
