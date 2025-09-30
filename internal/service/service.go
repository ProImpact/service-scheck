package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/ProImpact/passboult/internal/repository"
	"github.com/ProImpact/passboult/pkg/model"
)

type Service struct {
	Name string
	Cmd  *exec.Cmd
	PID  int
}

// ServiceManager check for the manage all the background process
type ServiceManager struct {
	mu       sync.RWMutex
	Services map[string]*BackgroundProcessCheker
	repo     *repository.ServiceRepository
}

func NewServiceManager(db *sql.DB) *ServiceManager {
	return &ServiceManager{
		Services: make(map[string]*BackgroundProcessCheker),
		repo:     repository.NewServiceRepository(db),
	}
}

// CreateService registrer a new service in the process manager and save it in the database
// also creates a new sistem service with the cmd
func (sm *ServiceManager) CreateService(name, healtCheckEndpoint, command string, pingTime time.Duration, stdout, stderr *os.File, args []string) (int, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, exists := sm.Services[name]; exists {
		return 0, fmt.Errorf("the service '%s' already exists", name)
	}

	cmd := exec.Command(command, args...)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()
	if err != nil {
		slog.Error("Can not execute the command", "service", name, "error", err)
		return 0, err
	}

	pid := cmd.Process.Pid

	cleaner := make([]func() error, 0)
	cleaner = append(cleaner, stderr.Close, stdout.Close)

	sm.Services[name] = NewBackgroundProcessChecker(&model.Service{
		ServiceName:        name,
		StartupTime:        time.Now(),
		Status:             model.Bootstraping,
		Command:            command,
		HealtCheckEndpoint: healtCheckEndpoint,
		PingTime:           &pingTime,
	}, sm.repo, pid, cleaner)

	// Run the service as a command
	sm.Services[name].Run()

	go func() {
		for errs := range sm.Services[name].Errors {
			if errs != nil {
				stderr.Write([]byte(errs.Error()))
				stderr.Write([]byte("\n"))
			}
		}
	}()

	go func() {
		time.Sleep(time.Second * 2)
		err := cmd.Wait()
		if err != nil {
			slog.Warn("Error creating the service", "service", name, "pid", pid, "exit_error", err)
		}
		time.Sleep(time.Second * 2)
		if cmd.Err != nil {
			slog.Warn("Service error", "service", name, "pid", pid, "err", cmd.Err)
		}
		sm.mu.Lock()
		delete(sm.Services, name)
		sm.mu.Unlock()
	}()

	return pid, nil
}

func (sm *ServiceManager) StopService(name string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	dbService, err := sm.repo.GetService(name)
	if err != nil {
		if errors.Is(err, repository.ErrServiceNotFound) {
			return fmt.Errorf("el servicio '%s' no esta registrado", name)
		}
		slog.Error(err.Error())
		return fmt.Errorf("internal server error")
	}

	service, exists := sm.Services[name]
	if !exists {
		slog.Error(fmt.Sprintf("el servicio '%s' no es manejado por el process manager", name))
	}
	if service != nil {
		service.Close()
	}
	slog.Info("Deteniendo el servicio...", "service", name, "pid", dbService.Pid)
	cmd := exec.Command("kill", fmt.Sprintf("%d", dbService.Pid))
	err = cmd.Start()
	if err != nil {
		slog.Warn("SIGTERM falló, intentando con SIGKILL...", "service", name, "error", err.Error())
		return fmt.Errorf("falló al detener el servicio '%s' con SIGTERM y SIGKILL: %v", name, err.Error())
	}
	go func() {
		err = cmd.Wait()
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	err = sm.repo.DeleteService(name)
	if err != nil {
		return err
	}
	delete(sm.Services, name)
	slog.Info("Service stopped", "service", name)
	return nil
}
