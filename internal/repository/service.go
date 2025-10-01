package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"github.com/ProImpact/passboult/internal/db"
	"github.com/ProImpact/passboult/pkg/model"
)

var ErrServiceAlreadyCreated = errors.New("service already created")
var ErrServiceNotFound = errors.New("service not found")

type ServiceRepository struct {
	DB      *sql.DB
	queries db.Queries
}

func NewServiceRepository(database *sql.DB) *ServiceRepository {
	return &ServiceRepository{
		DB:      database,
		queries: *db.New(database),
	}
}

func (repo *ServiceRepository) CreateService(service model.Service) error {
	_, err := repo.queries.ServiceGetByName(context.Background(), service.ServiceName)
	if errors.Is(err, sql.ErrNoRows) {
		return repo.queries.ServiceCreate(context.Background(), db.ServiceCreateParams{
			ID:                 hex.EncodeToString(sha256.New().Sum([]byte(service.ServiceName))),
			ServiceName:        service.ServiceName,
			StartupTime:        time.Now(),
			Status:             string(model.Ready),
			Command:            service.Command,
			HealtcheckEndpoint: service.HealtCheckEndpoint,
			PingTime:           service.PingTime.String(),
			Pid:                int64(service.Pid),
		})
	}
	if err == nil {
		return ErrServiceAlreadyCreated
	}
	return err
}

func (repo *ServiceRepository) GetService(serviceName string) (*model.Service, error) {
	dbService, err := repo.queries.ServiceGetByName(context.Background(), serviceName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrServiceNotFound
	}
	if err != nil {
		return nil, err
	}
	startupTime := dbService.StartupTime.(time.Time)
	pingTime, err := time.ParseDuration(dbService.PingTime)
	if err != nil {
		slog.Error(err.Error())
	}
	return &model.Service{
		ServiceName:        dbService.ServiceName,
		StartupTime:        startupTime,
		UpTime:             time.Since(startupTime),
		Status:             model.ServiceStatus(dbService.Status),
		Command:            dbService.Command,
		HealtCheckEndpoint: dbService.HealtcheckEndpoint,
		PingTime:           &pingTime,
		Pid:                int(dbService.Pid),
	}, nil
}

func (repo *ServiceRepository) UpdateService(params model.UpdateServiceParams) error {
	if params.ServiceName == nil {
		return errors.New("service name not specified in the update")
	}
	srv, err := repo.queries.ServiceGetByName(context.Background(), *params.ServiceName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrServiceNotFound
		}
		return err
	}
	if params.Status != nil {
		srv.Status = string(*params.Status)
	}
	if params.ServiceName != nil {
		srv.ServiceName = *params.ServiceName
	}
	if params.Command != nil {
		srv.Command = *params.Command
	}
	if params.HealtCheckEndpoint != nil {
		srv.HealtcheckEndpoint = *params.HealtCheckEndpoint
	}
	if params.PingTime != nil {
		srv.PingTime = params.PingTime.String()
	}
	return repo.queries.ServiceFullUpdate(context.Background(), db.ServiceFullUpdateParams{
		ServiceName:        srv.ServiceName,
		Status:             srv.Status,
		Command:            srv.Command,
		HealtcheckEndpoint: srv.HealtcheckEndpoint,
		PingTime:           srv.PingTime,
		ServiceName_2:      *params.ServiceName,
		Pid:                srv.Pid,
	})
}

func (repo *ServiceRepository) DeleteService(serviceName string) error {
	ctx := context.TODO()
	_, err := repo.queries.ServiceGetByName(ctx, serviceName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrServiceNotFound
		}
		return err
	}
	return repo.queries.ServiceDeleteByName(ctx, serviceName)
}

func (repo *ServiceRepository) GetAllServices() ([]model.Service, error) {
	services, err := repo.queries.ServiceGetAll(context.Background())
	if err != nil {
		return nil, err
	}
	var m []model.Service
	for _, srv := range services {
		startupTime := srv.StartupTime.(time.Time)
		pingTime, err := time.ParseDuration(srv.PingTime)
		if err != nil {
			return nil, err
		}
		m = append(m, model.Service{
			ServiceName:        srv.ServiceName,
			StartupTime:        startupTime,
			UpTime:             time.Since(startupTime),
			Status:             model.ServiceStatus(srv.Status),
			Command:            srv.Command,
			HealtCheckEndpoint: srv.HealtcheckEndpoint,
			PingTime:           &pingTime,
			Pid:                int(srv.Pid),
		})
	}
	return m, nil
}
