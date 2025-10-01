package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log/slog"
	"time"

	"github.com/ProImpact/service-check/internal/db"
	"github.com/ProImpact/service-check/pkg/model"
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
			Status:             string(model.Bootstraping),
			CommandType:        service.Check.Type,
			ExecuteCommand:     service.ExecCommand,
			HealtcheckEndpoint: service.Check.CheckCommand,
			PingTime:           service.PingTime.String(),
			Pid:                int64(service.Pid),
			CmdCheckCommand: sql.NullString{
				String: service.Check.CheckCommand,
				Valid:  true,
			},
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
	if dbService.CommandType == model.CMD {
		return &model.Service{
			ServiceName: dbService.ServiceName,
			StartupTime: startupTime,
			UpTime:      time.Since(startupTime),
			Status:      model.ServiceStatus(dbService.Status),
			Check: model.Check{
				Type:         model.CMD,
				CheckCommand: dbService.CmdCheckCommand.String,
			},
			PingTime: &pingTime,
			Pid:      int(dbService.Pid),
		}, nil
	}
	return &model.Service{
		ServiceName: dbService.ServiceName,
		StartupTime: startupTime,
		UpTime:      time.Since(startupTime),
		Status:      model.ServiceStatus(dbService.Status),
		Check: model.Check{
			Type:         model.REST,
			CheckCommand: dbService.HealtcheckEndpoint,
		},
		PingTime: &pingTime,
		Pid:      int(dbService.Pid),
	}, nil
}

func (repo *ServiceRepository) UpdateServiceStatus(params model.UpdateServiceParams) error {
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
	return repo.queries.ServiceUpdateStatus(context.Background(), db.ServiceUpdateStatusParams{
		Status:      srv.Status,
		ServiceName: srv.ServiceName,
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
		if srv.CommandType == model.REST {
			m = append(m, model.Service{
				ServiceName: srv.ServiceName,
				StartupTime: startupTime,
				UpTime:      time.Since(startupTime),
				Status:      model.ServiceStatus(srv.Status),
				Check: model.Check{
					Type:         model.REST,
					CheckCommand: srv.HealtcheckEndpoint,
				},
				PingTime: &pingTime,
				Pid:      int(srv.Pid),
			})
			continue
		}
		m = append(m, model.Service{
			ServiceName: srv.ServiceName,
			StartupTime: startupTime,
			UpTime:      time.Since(startupTime),
			Status:      model.ServiceStatus(srv.Status),
			Check: model.Check{
				Type:         model.CMD,
				CheckCommand: srv.CmdCheckCommand.String,
			},
			PingTime: &pingTime,
			Pid:      int(srv.Pid),
		})
	}
	return m, nil
}
