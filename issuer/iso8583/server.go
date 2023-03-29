package iso8583

import (
	"fmt"

	iso8583Server "github.com/moov-io/iso8583-connection/server"
)

type server struct {
	Addr string

	server *iso8583Server.Server
}

func NewServer(addr string) *server {
	s := iso8583Server.New(spec, readMessageLength, writeMessageLength)

	return &server{
		Addr:   addr,
		server: s,
	}
}

func (s *server) Start() error {
	fmt.Printf("Starting ISO 8583 server on %s...\n", s.Addr)

	if err := s.server.Start(s.Addr); err != nil {
		return fmt.Errorf("starting iso8583 server: %w", err)
	}

	// if the server is started successfully, update the address as it might be
	// different from the one we passed to the Start() method (e.g. if we passed
	// ":0" to let the OS choose a free port)
	s.Addr = s.server.Addr

	fmt.Printf("ISO 8583 server started on %s\n", s.Addr)

	return nil
}

func (s *server) Shutdown() error {
	fmt.Println("Shutting down iso8583 server...")

	s.server.Close()

	return nil
}
