package main

import (
	"net"
)

type Peer struct {
	conn 	net.Conn
	msgChan chan []byte
}

func NewPeer(conn net.Conn, msgChan chan []byte) *Peer {
	return &Peer{
		conn: conn,
		msgChan: msgChan,
	}
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

		p.msgChan <- msgBuff
	}
} 