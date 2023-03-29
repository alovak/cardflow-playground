package iso8583

import (
	"fmt"

	iso8583Server "github.com/moov-io/iso8583-connection/server"
	"golang.org/x/exp/slog"
)

type server struct {
	Addr string

	server *iso8583Server.Server
	logger *slog.Logger
}

func NewServer(logger *slog.Logger, addr string) *server {
	logger = logger.With(slog.String("type", "iso8583-server"), slog.String("addr", addr))

	s := iso8583Server.New(spec, readMessageLength, writeMessageLength)

	return &server{
		Addr:   addr,
		server: s,
		logger: logger,
	}
}

func (s *server) Start() error {
	s.logger.Info("starting ISO 8583 server...")

	if err := s.server.Start(s.Addr); err != nil {
		return fmt.Errorf("starting ISO 8583 server: %w", err)
	}

	// if the server is started successfully, update the address as it might be
	// different from the one we passed to the Start() method (e.g. if we passed
	// ":0" to let the OS choose a free port)
	s.Addr = s.server.Addr

	s.logger.Info("ISO 8583 server started", slog.String("addr", s.Addr))

	return nil
}

func (s *server) Close() error {
	s.logger.Info("shutting down ISO 8583 server...")

	s.server.Close()

	s.logger.Info("ISO 8583 server shut down")

	return nil
}
