package model

import (
	"time"
)

type ServiceStatus string

const (
	Down            ServiceStatus = "DOWN"
	Ready           ServiceStatus = "READY"
	Bootstraping    ServiceStatus = "BOOTSTRAPING"
	InvalidEndpoint ServiceStatus = "INVALID_ENDPOINT"
)

type Service struct {
	ServiceName        string         `json:"service_name,omitempty"`
	StartupTime        time.Time      `json:"startup_time,omitempty"`
	UpTime             time.Duration  `json:"up_time,omitempty"`
	Status             ServiceStatus  `json:"status,omitempty"`
	Command            string         `json:"command,omitempty"`
	HealtCheckEndpoint string         `json:"healt_check_endpoint,omitempty"`
	PingTime           *time.Duration `json:"ping_time,omitempty"`
	Pid                int            `json:"pid,omitempty"`
}

type ServiceCreate struct {
	ServiceName        string `json:"service_name,omitempty"`
	Command            string `json:"command,omitempty"`
	HealtCheckEndpoint string `json:"healt_check_endpoint,omitempty"`
	PingTime           string `json:"ping_time,omitempty"`
}

type UpdateServiceParams struct {
	Status             *ServiceStatus `json:"status,omitempty"`
	ServiceName        *string        `json:"service_name,omitempty"`
	Command            *string        `json:"command,omitempty"`
	HealtCheckEndpoint *string        `json:"healt_check_endpoint,omitempty"`
	PingTime           *time.Duration `json:"ping_time,omitempty"`
}
