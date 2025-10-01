package service

import (
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/ProImpact/service-check/internal/repository"
	"github.com/ProImpact/service-check/pkg/httpclient"
	"github.com/ProImpact/service-check/pkg/model"
)

// BackgroundProcessCheker check for the healtyness of a service
type BackgroundProcessCheker struct {
	PID          int
	Service      *model.Service
	done         chan struct{}
	Errors       chan error
	repo         *repository.ServiceRepository
	cleanupFuncs []func() error
}

var ErrServiceUnabailable = errors.New("service unavailable")

func NewBackgroundProcessChecker(serv *model.Service, repo *repository.ServiceRepository, pid int, cleaner []func() error) *BackgroundProcessCheker {
	return &BackgroundProcessCheker{
		Service: serv,
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
	tick := time.NewTicker(*b.Service.PingTime)
	slog.Info("running background processor for", "service", b.Service.ServiceName)
	go func() {
		down := model.Down
	FOR:
		for {
			select {
			case <-tick.C:
				if b.Service.Check.Type == model.REST {
					resp := httpclient.MakeRequest(b.Service.Check.CheckCommand)
					if resp == nil {
						err := b.repo.UpdateServiceStatus(model.UpdateServiceParams{
							Status:      &down,
							ServiceName: &b.Service.ServiceName,
						})
						if err != nil {
							b.Errors <- err
						}
						b.Errors <- err
						continue
					}
					if resp.StatusCode >= 200 || resp.StatusCode <= 299 {
						ready := model.Ready
						err := b.repo.UpdateServiceStatus(model.UpdateServiceParams{
							Status:      &ready,
							ServiceName: &b.Service.ServiceName,
						})
						if err != nil {
							b.Errors <- err
						}
					} else {
						invalid := model.InvalidEndpoint
						err := b.repo.UpdateServiceStatus(model.UpdateServiceParams{
							Status:      &invalid,
							ServiceName: &b.Service.ServiceName,
						})
						if err != nil {
							b.Errors <- err
						}
						b.Errors <- ErrServiceUnabailable
					}
				} else {
					comand := strings.Split(b.Service.Check.CheckCommand, " ")
					cmd := exec.Command(comand[0], comand[1:]...)
					output, err := cmd.CombinedOutput()
					statusCode := cmd.ProcessState.ExitCode()
					if err != nil {
						enhancedErr := fmt.Sprintf("failed to execute the comand: %s", output)
						switch statusCode {
						case 1:
							b.Errors <- fmt.Errorf("file or directory not found: %s", enhancedErr)
						case 2:
							b.Errors <- fmt.Errorf("bad use of the command: %s", enhancedErr)
						default:
							b.Errors <- fmt.Errorf("unexpected error code when executing the command: code %d: output: %s", statusCode, enhancedErr)
						}
						invalid := model.InvalidCheckCommand
						err := b.repo.UpdateServiceStatus(model.UpdateServiceParams{
							Status:      &invalid,
							ServiceName: &b.Service.ServiceName,
						})
						if err != nil {
							b.Errors <- err
						}
						b.Errors <- ErrServiceUnabailable
						continue
					}
					// service is ok
					ready := model.Ready
					err = b.repo.UpdateServiceStatus(model.UpdateServiceParams{
						Status:      &ready,
						ServiceName: &b.Service.ServiceName,
					})
					if err != nil {
						b.Errors <- err
					}
				}
			case <-b.done:
				tick.Stop()
				break FOR
			}
		}
	}()
}
