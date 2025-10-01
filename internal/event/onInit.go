package event

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/ProImpact/service-check/internal/db"
	"github.com/ProImpact/service-check/internal/repository"
	"github.com/ProImpact/service-check/internal/service"
	"github.com/ProImpact/service-check/pkg"
	"github.com/ProImpact/service-check/pkg/httpclient"
	"github.com/ProImpact/service-check/pkg/model"
)

func CheckServices(database *sql.DB, process *service.ServiceManager, logsDir string) {
	queries := db.New(database)
	slog.Info("Ckecking service health")
	services, err := queries.ServiceGetAll(context.Background())
	if err != nil {
		panic(err)
	}
	for _, srv := range services {
		resp := httpclient.MakeRequest(srv.HealtcheckEndpoint)
		// the service is down or not  created
		if resp == nil {
			slog.Error(
				fmt.Sprintf(
					"the service %s appears to be down becouse the endpoint is not avaliable or the endpoint is incorrect", srv.ServiceName,
				),
				"pid", srv.Pid,
			)
			fout, ferr, err := pkg.CreateLogsServicesFolderStructure(logsDir, srv.ServiceName)
			if err != nil {
				panic(err)
			}
			args := strings.Split(srv.Command, " ")
			pingTime, err := time.ParseDuration(srv.PingTime)
			if err != nil {
				panic(err)
			}
			cmd := exec.Command("kill", fmt.Sprintf("%d", srv.Pid))
			err = cmd.Start()
			if err != nil {
				slog.Error("SIGTERM falló, intentando con SIGKILL...", "service", srv.ServiceName, "error", err.Error())
				log.Fatal(err)
			}
			go func() {
				err = cmd.Wait()
				if err != nil {
					slog.Error(err.Error())
				}
			}()
			pid, err := process.CreateService(
				srv.ServiceName,
				srv.HealtcheckEndpoint,
				srv.Command,
				pingTime,
				fout,
				ferr,
				args[1:],
			)
			if err != nil {
				panic(err.Error())
			}
			err = queries.ServiceFullUpdate(context.Background(), db.ServiceFullUpdateParams{
				ServiceName:        srv.ServiceName,
				Status:             srv.Status,
				Command:            srv.Command,
				HealtcheckEndpoint: srv.HealtcheckEndpoint,
				PingTime:           srv.PingTime,
				ServiceName_2:      srv.ServiceName,
				Pid:                int64(pid),
			})
			if err != nil {
				panic(err.Error())
			}
			slog.Info("service recreated", "name", srv.ServiceName, "pid", pid)
			continue
		}
		startupTime := srv.StartupTime.(time.Time)
		pingTime, err := time.ParseDuration(srv.PingTime)
		if err != nil {
			slog.Error(err.Error())
		}
		process.Services[srv.ServiceName] = service.NewBackgroundProcessChecker(&model.Service{
			ServiceName:        srv.ServiceName,
			StartupTime:        startupTime,
			Status:             model.ServiceStatus(srv.Status),
			Command:            srv.Command,
			HealtCheckEndpoint: srv.HealtcheckEndpoint,
			PingTime:           &pingTime,
		}, repository.NewServiceRepository(database), int(srv.Pid), make([]func() error, 0))

		// Run the service as a command
		process.Services[srv.ServiceName].Run()
		slog.Info(fmt.Sprintf("Service %s is already running", srv.ServiceName))
	}
}
