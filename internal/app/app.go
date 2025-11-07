package app

import (
	"calendar-server/internal/config"
	handler "calendar-server/internal/delivery/http-server/handler/event_handler"
	"calendar-server/internal/delivery/http-server/router"
	repository "calendar-server/internal/repository/event_repository/inmemory"
	usecase "calendar-server/internal/usecase/event_usecase"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"calendar-server/pkg/logger/zappretty"

	"go.uber.org/zap"
)

// App представляет приложение
type App struct {
	config *config.Config
	server *http.Server
	logger *zap.Logger
}

// New создает новый экземпляр App
func New(logger *zap.Logger) *App {
	cfg := config.MustLoad()

	eventRepo := repository.NewEventRepository(logger)

	eventUseCase := usecase.NewEventUseCase(eventRepo, logger)

	eventHandler := handler.NewEventHandler(eventUseCase, logger)

	r := router.NewRouter(eventHandler, logger)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		config: cfg,
		server: server,
		logger: logger,
	}
}

// Run запускает сервер
func (a *App) Run() error {
	a.logger.Info("Starting server",
		zappretty.Field("port", a.config.Port),
		zappretty.Field("environment", a.config.Environment),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.handleSignals(cancel)

	serverErr := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		a.logger.Error("Server error", zappretty.Field("error", err))
		return err
	case <-ctx.Done():
		a.logger.Info("Shutdown signal received, stopping server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			a.logger.Error("Server shutdown failed", zappretty.Field("error", err))
			return err
		}

		a.logger.Info("Server stopped gracefully")
		return nil
	}
}

// handleSignals обрабатывает сигналы OS для graceful shutdown
func (a *App) handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	a.logger.Info("Received signal", zappretty.Field("signal", sig))
	cancel()
}
