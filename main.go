package main

import (
	"log/slog"
	"net"
)

const defaultListenAddress = ":5001"

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	ln net.Listener
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)

	if err != nil {
		return err
	}

	s.ln = ln

	return s.acceptLoop()

}

func (s *Server) acceptLoop() error {

	for {

		conn, err := s.ln.Accept()
	
		if err != nil {
			slog.Error("accept error: ", err)
			continue
		}

		s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	for {

	}
}

func NewServer(cfg Config) *Server {

	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress
	}

	return &Server{Config: cfg}
}

func main() {

}