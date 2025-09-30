package app

import (
	dsql "database/sql"
	"errors"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ProImpact/passboult/internal/config"
	"github.com/ProImpact/passboult/internal/db/sql"
	"github.com/ProImpact/passboult/internal/repository"
	"github.com/ProImpact/passboult/internal/rpc"
)

type ServerApp struct {
	db           *dsql.DB
	cleanupFuncs []func() error
	repo         *repository.ServiceRepository
	port         int
	listener     net.Listener
	cfg          *config.ServerConfiguration
}

func NewServerApp(cfg *config.ServerConfiguration) *ServerApp {
	dbDriver, err := config.ConnectToDatabase(cfg.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	err = sql.Migrate(dbDriver)
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewServiceRepository(dbDriver)
	srv := &ServerApp{
		db:           dbDriver,
		repo:         repo,
		cleanupFuncs: make([]func() error, 0),
		port:         cfg.Port,
	}
	srv.cleanupFuncs = append(srv.cleanupFuncs, repo.DB.Close)
	listener := rpc.CreateRPCServer(srv.repo, cfg.Port, cfg.LogsDir)
	srv.listener = listener
	srv.cleanupFuncs = append(srv.cleanupFuncs, listener.Close)
	err = os.Mkdir(cfg.LogsDir, 0750)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}
	}
	srv.cfg = cfg
	return srv
}

func (s *ServerApp) Run() {
	signTerm := make(chan os.Signal, 1)
	signal.Notify(signTerm, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{})
	slog.Info("Rpc server initializaded", "port", s.port)

	go http.Serve(s.listener, nil)

	go func() {
		<-signTerm
		slog.Warn("app shutdown,killing background process")
		for _, clean := range s.cleanupFuncs {
			err := clean()
			if err != nil {
				slog.Error(err.Error())
			}
		}
		done <- struct{}{}
	}()
	<-done
}
