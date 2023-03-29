package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/alovak/cardflow-playground/issuer"
	issuer8583 "github.com/alovak/cardflow-playground/issuer/iso8583"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

type IssuerApp struct {
	srv               *http.Server
	wg                *sync.WaitGroup
	Addr              string
	ISO8583ServerAddr string
	logger            *slog.Logger
	iso8583Server     io.Closer
}

func NewIssuerApp(logger *slog.Logger) *IssuerApp {
	logger = logger.With(slog.String("app", "issuer"))

	return &IssuerApp{
		wg:     &sync.WaitGroup{},
		logger: logger,
	}
}

func (a *IssuerApp) Run() {
	if err := a.Start(); err != nil {
		a.logger.Error("Error starting app", "err", err)
		os.Exit(1)
	}

	// Wait for interrupt signal to gracefully shutdown the app with all services
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	a.Shutdown()
}

func (a *IssuerApp) Start() error {
	a.logger.Info("starting app...")

	// setup the issuer
	router := chi.NewRouter()
	repository := issuer.NewRepository()

	iso8583Server := issuer8583.NewServer(a.logger, "127.0.0.1:0")
	err := iso8583Server.Start()
	if err != nil {
		return fmt.Errorf("starting iso8583 server: %w", err)
	}
	a.ISO8583ServerAddr = iso8583Server.Addr

	a.iso8583Server = iso8583Server

	iss := issuer.NewIssuer(repository)
	api := issuer.NewAPI(iss)
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
				a.logger.Error("starting http server", "err", err)
			}

			a.logger.Info("http server stopped")
		}

		a.wg.Done()
	}()

	return nil
}

func (a *IssuerApp) Shutdown() {
	a.logger.Info("shutting down app...")

	a.srv.Shutdown(context.Background())

	err := a.iso8583Server.Close()
	if err != nil {
		a.logger.Error("closing iso8583 server", "err", err)
	}

	a.wg.Wait()

	a.logger.Info("app stopped")
}
