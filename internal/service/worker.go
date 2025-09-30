package service

import (
	"errors"
	"log/slog"
	"time"

	"github.com/ProImpact/passboult/internal/repository"
	"github.com/ProImpact/passboult/pkg/httpclient"
	"github.com/ProImpact/passboult/pkg/model"
)

// BackgroundProcessCheker check for the healtyness of a service
type BackgroundProcessCheker struct {
	PID int
	model.Service
	done         chan struct{}
	Errors       chan error
	repo         *repository.ServiceRepository
	cleanupFuncs []func() error
}

var ErrServiceUnabailable = errors.New("service unavailable")

func NewBackgroundProcessChecker(serv *model.Service, repo *repository.ServiceRepository, pid int, cleaner []func() error) *BackgroundProcessCheker {
	return &BackgroundProcessCheker{
		Service: *serv,
		done:    make(chan struct{}),
		Errors:  make(chan error),
		PID:     pid,
		repo:    repo,
	}
}

func (b *BackgroundProcessCheker) Close() {
	b.done <- struct{}{}
	close(b.done)
	close(b.Errors)
	for _, clean := range b.cleanupFuncs {
		err := clean()
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func (b *BackgroundProcessCheker) Run() {
	tick := time.NewTicker(*b.PingTime)
	go func() {
		down := model.Down
	FOR:
		for {
			select {
			case <-tick.C:
				resp := httpclient.MakeRequest(b.HealtCheckEndpoint)
				if resp == nil {
					err := b.repo.UpdateService(model.UpdateServiceParams{
						Status:      &down,
						ServiceName: &b.ServiceName,
					})
					if err != nil {
						b.Errors <- err
					}
					b.Errors <- err
					continue
				}
				if resp.StatusCode >= 200 || resp.StatusCode <= 299 {
					ready := model.Ready
					err := b.repo.UpdateService(model.UpdateServiceParams{
						Status:      &ready,
						ServiceName: &b.ServiceName,
					})
					if err != nil {
						b.Errors <- err
					}
				} else {
					ready := model.InvalidEndpoint
					err := b.repo.UpdateService(model.UpdateServiceParams{
						Status:      &ready,
						ServiceName: &b.ServiceName,
					})
					if err != nil {
						b.Errors <- err
					}
					b.Errors <- ErrServiceUnabailable
				}
			case <-b.done:
				tick.Stop()
				break FOR
			}
		}
	}()
}
