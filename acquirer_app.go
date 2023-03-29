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
)

type AcquirerApp struct {
	srv  *http.Server
	wg   *sync.WaitGroup
	Addr string
}

func NewAcquirerApp() *AcquirerApp {
	return &AcquirerApp{
		wg: &sync.WaitGroup{},
	}
}

func (a *AcquirerApp) Run() {
	a.Start()

	// Wait for interrupt signal to gracefully shutdown the app with all services
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	a.Shutdown()
}

func (a *AcquirerApp) Start() error {
	fmt.Println("Starting acquirer app...")

	// setup the acquirer
	router := chi.NewRouter()
	repository := acquirer.NewRepository()

	// setup iso8583Client
	iso8583Client := iso8583.NewClient()

	acq := acquirer.NewAcquirer(repository, iso8583Client)
	api := acquirer.NewAPI(acq)
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
		fmt.Printf("Acquirer http server started on port %s\n", a.Addr)

		if err := a.srv.Serve(l); err != nil {
			if err != http.ErrServerClosed {
				fmt.Printf("Error starting acquirer http server: %v\n", err)
			}

			fmt.Println("Acquirer http server stopped")
		}

		a.wg.Done()
	}()

	return nil
}

func (a *AcquirerApp) Shutdown() {
	fmt.Println("Shutting down acquirer app...")

	a.srv.Shutdown(context.Background())

	a.wg.Wait()

	fmt.Println("Acquirer app stopped")
}
