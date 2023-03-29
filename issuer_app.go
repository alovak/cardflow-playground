package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/alovak/cardflow-playground/issuer"
	"github.com/go-chi/chi/v5"
)

type IssuerApp struct {
	srv  *http.Server
	wg   *sync.WaitGroup
	Addr string
}

func NewIssuerApp() *IssuerApp {
	return &IssuerApp{
		wg: &sync.WaitGroup{},
	}
}

func (a *IssuerApp) Run() {
	a.Start()

	// Wait for interrupt signal to gracefully shutdown the app with all services
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	a.Shutdown()
}

func (a *IssuerApp) Start() error {
	fmt.Println("Starting issuer app...")

	// setup the issuer
	router := chi.NewRouter()
	repository := issuer.NewRepository()
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
		fmt.Printf("Issuer http server started on port %s\n", a.Addr)

		if err := a.srv.Serve(l); err != nil {
			if err != http.ErrServerClosed {
				fmt.Printf("Error starting issuer http server: %v\n", err)
			}

			fmt.Println("Issuer http server stopped")
		}

		a.wg.Done()
	}()

	return nil
}

func (a *IssuerApp) Shutdown() {
	fmt.Println("Shutting down issuer app...")

	a.srv.Shutdown(context.Background())

	a.wg.Wait()

	fmt.Println("Issuer app stopped")
}
