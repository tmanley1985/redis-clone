package main

import (
	"net"
)

type Peer struct {
	conn net.Conn
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{conn}
}

func (p *Peer) readLoop() error {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)

		if err != nil {
			return err
		}

		msgBuff := make([]byte, n)
		copy(msgBuff, buf[:n])
	}
} 