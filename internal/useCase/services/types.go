package services

import (
	"inttest-runtime/internal/domain"
	"time"
)

type ServiceStatusResp struct {
	ID            domain.ServiceID `json:"id"`
	Type          string           `json:"type"`
	LaunchedAt    time.Time        `json:"launched_at"`
	CurrentStatus ServiceStatus    `json:"status"`
}

type ServiceStatus struct {
	Status    string     `json:"status"`
	UpdatedAt *time.Time `json:"updated_at"`
}
