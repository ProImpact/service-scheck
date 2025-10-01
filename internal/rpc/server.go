package rpc

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/rpc"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ProImpact/service-check/internal/event"
	"github.com/ProImpact/service-check/internal/repository"
	"github.com/ProImpact/service-check/internal/service"
	"github.com/ProImpact/service-check/pkg"
	model "github.com/ProImpact/service-check/pkg/model"
)

type Server struct {
	Repo    *repository.ServiceRepository
	Serv    *service.ServiceManager
	LogsDir string
}

func CreateRPCServer(repo *repository.ServiceRepository, port int, logsDir string) net.Listener {
	srvManager := service.NewServiceManager(repo.DB)
	service := &Server{
		Repo: repo,
		Serv: srvManager,
	}
	err := rpc.Register(service)
	if err != nil {
		log.Fatal(err)
	}
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	event.CheckServices(repo.DB, srvManager, logsDir)
	return listener
}

type ServiceArgs struct {
	Service model.ServiceCreate `json:"service"`
}

type ServiceArgsUpdate struct {
	Service model.UpdateServiceParams
}

type ServiceNameArgs struct {
	ServiceName string
}

type ServiceListReply struct {
	Services []model.Service
}

type ServiceLogsReply struct {
	Data string
}

func (s *Server) CreateService(args *ServiceArgs, reply *bool) error {
	pkg.MustPrint(pkg.Event{
		Name: "CreateService",
		Data: args,
	})
	err := os.MkdirAll(path.Join(s.LogsDir, args.Service.ServiceName), 0750)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	fstdout, fstderr, err := pkg.CreateLogsServicesFolderStructure(s.LogsDir, args.Service.ServiceName)
	if err != nil {
		return err
	}
	comand := strings.Split(args.Service.Check.CheckCommand, " ")
	duration, err := time.ParseDuration(args.Service.PingTime)
	if err != nil {
		return err
	}
	pid, err := s.Serv.CreateService(
		args.Service.ServiceName,
		args.Service.Check.CheckCommand,
		args.Service.Check.Type,
		args.Service.Command,
		args.Service.Check.CheckCommand,
		duration,
		fstdout,
		fstderr,
		comand[1:],
	)
	if err != nil {
		*reply = false
		return err
	}
	err = s.Repo.CreateService(model.Service{
		ServiceName: args.Service.ServiceName,
		StartupTime: time.Now(),
		Status:      model.Bootstraping,
		Check:       args.Service.Check,
		PingTime:    &duration,
		Pid:         pid,
		ExecCommand: args.Service.Command,
	})
	if err != nil {
		*reply = false
		return err
	}
	*reply = true
	return nil
}

func (s *Server) GetService(args *ServiceNameArgs, reply *model.Service) error {
	result, err := s.Repo.GetService(args.ServiceName)
	if err != nil {
		return err
	}
	*reply = *result
	return nil
}

func (s *Server) Logs(args *ServiceNameArgs, reply *ServiceLogsReply) error {
	f, err := os.Open(filepath.Join(s.LogsDir, args.ServiceName, "stderr.log"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("service not found")
		}
		slog.Error(err.Error())
		return fmt.Errorf("server error handling the request")
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	*reply = ServiceLogsReply{
		Data: string(data),
	}
	return nil
}

func (s *Server) DeleteService(args *ServiceNameArgs, reply *bool) error {
	err := s.Serv.StopService(args.ServiceName)
	if err != nil {
		*reply = false
		return err
	}
	*reply = true
	return nil
}

func (s *Server) GetAllServices(_ struct{}, reply *ServiceListReply) error {
	services, err := s.Repo.GetAllServices()
	if err != nil {
		return err
	}
	reply.Services = services
	return nil
}
