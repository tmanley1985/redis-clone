package main

import (
	"fmt"
	"io"
	"log/slog"
	"net"
)

type Peer struct {
	conn      net.Conn
	dataStore *DataStore
}

func NewPeer(conn net.Conn, ds *DataStore) *Peer {
	return &Peer{
		conn:      conn,
		dataStore: ds,
	}
}

func (p *Peer) Handle() error {
	slog.Info("Handling new peer connection", "remoteAddr", p.conn.RemoteAddr())

	responseChan := make(chan ClientResponse)
	defer close(responseChan)

	// Start a goroutine to process responses
	go func() {
		err := p.processResponses(responseChan)
		if err != nil {
			slog.Error("Error processing responses", "remoteAddr", p.conn.RemoteAddr(), "error", err)
		}
	}()

	for {
		cmd, err := ParseCommand(p.conn)
		if err != nil {
			if err == io.EOF {
				slog.Info("Client disconnected", "remoteAddr", p.conn.RemoteAddr())
				return nil // Client closed connection normally
			}
			slog.Error("Error parsing command", "remoteAddr", p.conn.RemoteAddr(), "error", err)
			p.sendErrorResponse(responseChan, fmt.Errorf("parse error: %w", err))
			continue // Continue to next command instead of returning
		}

		slog.Info("Handling command", "command", cmd, "remoteAddr", p.conn.RemoteAddr())

		// TODO: Possibly we need a different channel for errors and responses
		// I'm sure it isn't the responsibility of the datastore to handle errors unrelated to storage like peer connection issues.
		cmd.Execute(p.dataStore, responseChan)
		// if err != nil {
		// 	slog.Error("Error executing command", "remoteAddr", p.conn.RemoteAddr(), "error", err)
		// 	p.sendErrorResponse(responseChan, fmt.Errorf("execution error: %w", err))
		// 	// Don't return here, allow the client to send more commands
		// }
	}
}

func (p *Peer) sendErrorResponse(responseChan chan ClientResponse, err error) {
	responseChan <- NewErrorResponse("parse error")
}

func (p *Peer) processResponses(responseChan chan ClientResponse) error {
	for response := range responseChan {

		// Format and send successful response
		_, err := p.conn.Write(response.Serialize())
		if err != nil {
			return fmt.Errorf("error sending response: %w", err)
		}
	}
	return nil
}
