package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/alovak/cardflow-playground/acquirer"
	"github.com/alovak/cardflow-playground/log"
)

func main() {
	iso8583Address := flag.String("iso8583Addr", "127.0.0.1:8583", "Address of the ISO8583 server")

	flag.Parse()

	logger := log.New()
	app := acquirer.NewApp(logger, *iso8583Address)

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
