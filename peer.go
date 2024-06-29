package main

import (
	"net"
)

type Peer struct {
	conn      net.Conn
	msgChan   chan []byte
	dataStore *DataStore
}

func NewPeer(conn net.Conn, msgChan chan []byte, ds *DataStore) *Peer {
	return &Peer{
		conn:      conn,
		msgChan:   msgChan,
		dataStore: ds,
	}
}

func (p *Peer) Handle() error {
	// Make a channel for the response. We need this here because there could be many client
	// connections and we wouldn't want to send someone another client's response.
	responseChan := make(chan ClientResponse)
	defer close(responseChan)

	// In a goroutine, have a little function that just sits there and tries to process responses.
	go p.processResponses(responseChan)

	// Parse the command given the input for as long as we need to.
	for {
		cmd, err := ParseCommand(p.conn)
		if err != nil {
			return nil
		}
		cmd.Execute(p.dataStore, responseChan)
	}
}

func (p *Peer) processResponses(chan ClientResponse) {

	// This will just keep looking for responses for this particular peer. We don't want to mix up
	// the responses between peers. That would be... bad.

}
