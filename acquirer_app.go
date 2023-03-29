package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/alovak/cardflow-playground/acquirer"
	"github.com/alovak/cardflow-playground/acquirer/iso8583"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

type AcquirerApp struct {
	srv               *http.Server
	wg                *sync.WaitGroup
	Addr              string
	ISO8583ServerAddr string
	logger            *slog.Logger
}

func NewAcquirerApp(logger *slog.Logger, iso8583ServerAddr string) *AcquirerApp {
	logger = logger.With(slog.String("app", "acquirer"))

	return &AcquirerApp{
		logger:            logger,
		wg:                &sync.WaitGroup{},
		ISO8583ServerAddr: iso8583ServerAddr,
	}
}

func (a *AcquirerApp) Run() {
	err := a.Start()
	if err != nil {
		a.logger.Error("Error starting app", "err", err)
		os.Exit(1)
	}

	// Wait for interrupt signal to gracefully shutdown the app with all services
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	a.Shutdown()
}

func (a *AcquirerApp) Start() error {
	a.logger.Info("starting app...")

	// setup the acquirer
	router := chi.NewRouter()
	repository := acquirer.NewRepository()

	// setup iso8583Client
	iso8583Client, err := iso8583.NewClient(a.logger, a.ISO8583ServerAddr)
	if err != nil {
		return fmt.Errorf("creating iso8583 client: %w", err)
	}

	// connect to iso8583 server
	if err := iso8583Client.Connect(); err != nil {
		return fmt.Errorf("connecting to iso8583 server: %w", err)
	}

	acq := acquirer.NewAcquirer(repository, iso8583Client)
	api := acquirer.NewAPI(a.logger, acq)
	api.AppendRoutes(router)

	l, err := net.Listen("tcp", "127.0.0.1:0")
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

func (a *AcquirerApp) Shutdown() {
	a.logger.Info("shutting down app...")

	a.srv.Shutdown(context.Background())

	a.wg.Wait()

	a.logger.Info("app stopped")
}
