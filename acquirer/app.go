package acquirer

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/alovak/cardflow-playground/acquirer/iso8583"
	"github.com/alovak/cardflow-playground/internal/middleware"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

// App is the main application, it contains all the components of the issuer service
// and is responsible for starting and stopping them.
type App struct {
	srv               *http.Server
	wg                *sync.WaitGroup
	Addr              string
	ISO8583ServerAddr string
	logger            *slog.Logger
	config            *Config
}

func NewApp(logger *slog.Logger, config *Config) *App {
	logger = logger.With(slog.String("app", "acquirer"))

	if config == nil {
		config = DefaultConfig()
	}

	return &App{
		logger: logger,
		wg:     &sync.WaitGroup{},
		config: config,
	}
}

func (a *App) Start() error {
	a.logger.Info("starting app...")

	// setup the acquirer
	router := chi.NewRouter()
	router.Use(middleware.NewStructuredLogger(a.logger))

	repository := NewRepository()

	// setup iso8583Client
	stanGenerator := iso8583.NewStanGenerator()
	iso8583Client, err := iso8583.NewClient(a.logger, a.config.ISO8583Addr, stanGenerator)
	if err != nil {
		return fmt.Errorf("creating iso8583 client: %w", err)
	}

	// connect to iso8583 server
	if err := iso8583Client.Connect(); err != nil {
		return fmt.Errorf("connecting to iso8583 server: %w", err)
	}

	acq := NewService(repository, iso8583Client)
	api := NewAPI(a.logger, acq)
	api.AppendRoutes(router)

	l, err := net.Listen("tcp", a.config.HTTPAddr)
	if err != nil {
		return fmt.Errorf("listening tcp port: %w", err)
	}

	a.Addr = l.Addr().String()

	a.srv = &http.Server{
		Handler: router,
	}

	a.wg.Add(1)
	go func() {
		a.logger.Info("http server started", slog.String("addr", a.Addr))

		if err := a.srv.Serve(l); err != nil {
			if err != http.ErrServerClosed {
				a.logger.Error("Error starting acquirer http server", "err", err)
			}

			a.logger.Info("http server stopped")
		}

		a.wg.Done()
	}()

	return nil
}

func (a *App) Shutdown() {
	a.logger.Info("shutting down app...")

	a.srv.Shutdown(context.Background())

	a.wg.Wait()

	a.logger.Info("app stopped")
}
