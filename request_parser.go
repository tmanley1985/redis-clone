package main

import (
	"bufio"
	"errors"
	"io"
	"log/slog"
	"strconv"
	"strings"
)

func ParseCommand(r io.Reader) (RequestCommand, error) {

	// Note to self: A buffered reader WRAPS another reader in order to prevent so many
	// network calls. It's a bit like preventing someone from going to the grocery store for
	// each item when instead they can go once at put them in the fridge (the buffer).
	reader := bufio.NewReader(r)

	// Read the first byte
	firstByte, err := reader.ReadByte()
	slog.Info("First byte" + string(firstByte))
	if err != nil {
		return nil, err
	}

	// Example: *3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nHello\r\n
	if firstByte != '*' {
		slog.Info("Got here")
		return nil, errors.New("invalid RESP: expected array")
	}

	slog.Info("Got here now boah")
	// All redis requests (this is a useful lie) are sent as arrays where each item is a bulk string.
	// So we need to see what the length of the array is here.
	numElements, err := readInteger(reader)
	if err != nil {
		return nil, err
	}

	slog.Info("About to read this bulk ass string")
	// Read each element
	elements := make([]string, numElements)
	for i := 0; i < numElements; i++ {
		element, err := readBulkString(reader)
		if err != nil {
			return nil, err
		}
		elements[i] = element
	}

	slog.Info("About to parse")
	// Parse the command based on the elements
	switch strings.ToUpper(elements[0]) {
	case "SET":
		if len(elements) != 3 {
			return nil, errors.New("SET command requires 2 arguments")
		}
		slog.Info("We have a set command here")
		return NewSetCommand(elements[1], elements[2]), nil
	case "GET":
		if len(elements) != 2 {
			return nil, errors.New("GET command requires 1 argument")
		}
		return NewGetCommand(elements[1]), nil
	default:
		return nil, errors.New("unknown command")
	}
}

func readInteger(reader *bufio.Reader) (int, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(line))
}

func readBulkString(reader *bufio.Reader) (string, error) {
	// Read the $ byte
	firstByte, err := reader.ReadByte()
	if err != nil {
		return "", err
	}
	if firstByte != '$' {
		return "", errors.New("invalid RESP: expected bulk string")
	}

	// Read the length
	length, err := readInteger(reader)
	if err != nil {
		return "", err
	}

	// Read the string
	str := make([]byte, length)
	_, err = io.ReadFull(reader, str)
	if err != nil {
		return "", err
	}

	// Read the trailing \r\n so that this isn't hanging around for the next time we read.
	_, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return string(str), nil
}

// Command represents a Redis command
type RequestCommand interface {
	// Name returns the name of the command (e.g., "GET", "SET")
	Type() string

	// Args returns the arguments of the command
	Args() []string

	Execute(ds *DataStore, responseChannel chan ClientResponse)
}

type SetCommand struct {
	key, val string
}

func (c SetCommand) Type() string {
	return "SET"
}

func (c SetCommand) Args() []string {
	return []string{c.key, c.val}
}

func (c *SetCommand) Execute(ds *DataStore, responseChannel chan ClientResponse) {
	// The data store will be responsible for handling all of it's operations and it will send a message to the response channel
	// so it's doing it's own thing over there. This is a good way to think about goroutines, they're just little dudes doing their own
	// thing.
	ds.Set(c.key, c.val, responseChannel)
}

func NewSetCommand(key, val string) *SetCommand {
	return &SetCommand{
		key: key,
		val: val,
	}
}

type GetCommand struct {
	key string
}

func (c GetCommand) Type() string {
	return "GET"
}

func (c GetCommand) Args() []string {
	return []string{c.key}
}

func (c *GetCommand) Execute(ds *DataStore, responseChannel chan ClientResponse) {
	ds.Get(c.key, responseChannel)
}

func NewGetCommand(key string) *GetCommand {
	return &GetCommand{
		key: key,
	}
}
