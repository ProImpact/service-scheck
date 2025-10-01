package service

import (
	"context"

	"github.com/ProImpact/service-check/pkg/model"
)

type ServiceServicer interface {
	CreateService(ctx context.Context, serviceName string) (*model.Service, error)
	GetServices(ctx context.Context) ([]model.Service, error)
	KillService(ctx context.Context, serviceName string) error
	GetServiceInfo(ctx context.Context, serviceName string) (model.Service, error)
}

// Run a process
type Runner interface {
	Run(ctx context.Context) error
	Check(ctx context.Context) error
	GetInfo(ctx context.Context) map[string]any
}
