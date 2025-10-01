package model

import (
	"time"
)

type ServiceStatus string

const (
	Down                ServiceStatus = "DOWN"
	Ready               ServiceStatus = "READY"
	Bootstraping        ServiceStatus = "BOOTSTRAPING"
	InvalidEndpoint     ServiceStatus = "INVALID_ENDPOINT"
	InvalidCheckCommand ServiceStatus = "CHECK_COMMAND_ERROR"
)

type Service struct {
	ServiceName string         `json:"service_name,omitempty"`
	StartupTime time.Time      `json:"startup_time,omitempty"`
	UpTime      time.Duration  `json:"up_time,omitempty"`
	Status      ServiceStatus  `json:"status,omitempty"`
	Check       Check          `json:"command,omitempty"`
	ExecCommand string         `json:"exec_command,omitempty"`
	PingTime    *time.Duration `json:"ping_time,omitempty"`
	Pid         int            `json:"pid,omitempty"`
}

type CheckType = string

const (
	CMD  CheckType = "CMD"
	REST CheckType = "REST"
)

type Check struct {
	Type         CheckType `json:"type"`
	CheckCommand string    `json:"check_command"`
}

type ServiceCreate struct {
	ServiceName string `json:"service_name"`
	Check       Check  `json:"check"`
	PingTime    string `json:"ping_time"`
	Command     string `json:"command"`
}

type UpdateServiceParams struct {
	Status             *ServiceStatus `json:"status,omitempty"`
	ServiceName        *string        `json:"service_name,omitempty"`
	ExecuteCommand     *string        `json:"execute_command,omitempty"`
	HealtCheckEndpoint *string        `json:"healt_check_endpoint,omitempty"`
	PingTime           *time.Duration `json:"ping_time,omitempty"`
}
