package impl

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/podengo-project/idmsvc-backend/internal/config"
	"github.com/podengo-project/idmsvc-backend/internal/handler"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/router"
	"github.com/podengo-project/idmsvc-backend/internal/infrastructure/service"
	"golang.org/x/exp/slog"
)

type metricsService struct {
	context   context.Context
	cancel    context.CancelFunc
	waitGroup *sync.WaitGroup
	config    *config.Config

	echo *echo.Echo
}

func NewMetrics(ctx context.Context, wg *sync.WaitGroup, config *config.Config, app handler.Application) service.ApplicationService {
	if config == nil {
		panic("config is nil")
	}
	if wg == nil {
		panic("wg is nil")
	}

	routerConfig := router.RouterConfig{
		Version:     "v1.0",
		MetricsPath: config.Metrics.Path,
		PrivatePath: "/private",
		Handlers:    app,
	}

	result := &metricsService{}
	result.context, result.cancel = context.WithCancel(ctx)
	result.waitGroup = wg
	result.config = config

	result.echo = router.NewRouterForMetrics(
		echo.New(),
		routerConfig,
	)
	result.echo.HideBanner = true

	if result.config.Logging.Level == "debug" {
		routes := result.echo.Routes()
		slog.Debug("Printing metrics routes")
		for idx, route := range routes {
			slog.Debug("routing",
				slog.Int("index", idx),
				slog.Any("route", route),
			)
		}
	}

	return result
}

func (srv *metricsService) Start() error {
	srv.waitGroup.Add(2)
	go func() {
		defer srv.waitGroup.Done()
		srvAddress := fmt.Sprintf(":%d", srv.config.Metrics.Port)
		slog.Debug("metrics", slog.String("srvAddress", srvAddress))
		if err := srv.echo.Start(srvAddress); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start metricsService", slog.Any("error", err))
		}
	}()

	go func() {
		defer srv.waitGroup.Done()
		defer srv.cancel()
		<-srv.context.Done()
		slog.Info("Shutting down metricsService")
		if err := srv.echo.Shutdown(context.Background()); err != nil {
			slog.Error(
				"error shuttingdown metricsService",
				slog.Any("error", err),
			)
		}
	}()

	return nil
}

func (srv *metricsService) Stop() error {
	srv.cancel()
	return nil
}
