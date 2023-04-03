package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/alovak/cardflow-playground/issuer"
	"github.com/alovak/cardflow-playground/log"
)

func main() {
	logger := log.New()
	app := issuer.NewApp(logger, issuer.DefaultConfig())

	err := app.Start()
	if err != nil {
		logger.Error("Error starting app", "err", err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	app.Shutdown()
}
