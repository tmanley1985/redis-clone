package main

import (
	"fmt"
	"log/slog"
	"net"
)

const defaultListenAddress = ":5001"

type Config struct {
	ListenAddress string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgChan   chan []byte
	ds        *DataStore
}

func NewServer(cfg Config, ds *DataStore) *Server {

	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddress
	}

	return &Server{
		Config: cfg,
		peers:  make(map[*Peer]bool),
		ds:     ds,

		// Channels
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgChan:   make(chan []byte),
	}
}

func (s *Server) handleRawMsg(rawMsg []byte) error {
	fmt.Println(rawMsg)
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgChan:
			if err := s.handleRawMsg(rawMsg); err != nil {
				slog.Error("raw message error", "error", err)
			}
		case <-s.quitCh:
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		}
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)

	if err != nil {
		return err
	}

	s.ln = ln

	go s.loop()

	slog.Info("server running on", "listenAddr: ", defaultListenAddress)
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

	peer := NewPeer(conn, s.msgChan, s.ds)

	s.addPeerCh <- peer

	slog.Info("Peer connected", "remoteAddr", peer.conn.RemoteAddr())
	go peer.Handle()

}
