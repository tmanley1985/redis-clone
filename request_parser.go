package main

import (
	"bufio"
	"errors"
	"io"
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
	if err != nil {
		return nil, err
	}

	// Example: *3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nHello\r\n
	if firstByte != '*' {
		return nil, errors.New("invalid RESP: expected array")
	}

	// All redis requests (this is a useful lie) are sent as arrays where each item is a bulk string.
	// So we need to see what the length of the array is here.
	numElements, err := readInteger(reader)
	if err != nil {
		return nil, err
	}

	// Read each element
	elements := make([]string, numElements)
	for i := 0; i < numElements; i++ {
		element, err := readBulkString(reader)
		if err != nil {
			return nil, err
		}
		elements[i] = element
	}

	// Parse the command based on the elements
	switch strings.ToUpper(elements[0]) {
	case "SET":
		if len(elements) != 3 {
			return nil, errors.New("SET command requires 2 arguments")
		}
		return NewSetCommand(elements[1], elements[2]), nil
	// case "GET":
	// TODO:
	// 	if len(elements) != 2 {
	// 		return nil, errors.New("GET command requires 1 argument")
	// 	}
	// 	return NewGetCommand(elements[1]), nil
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

	// // String returns a string representation of the command
	// String() string
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

func (c *SetCommand) Execute() ResponseCommand {
	// I still need to set the command in storage and possibly append-only log it to disk
	// and all that jazz but I ain't that good just yet. Baby steps.
	//
	// Also, I think I may have to put this command on a channel instead of returning it.
	// Not sure what I want to do here, I'd like to keep this function pure.
	return NewSimpleStringResponse("OK")
}

func NewSetCommand(key, val string) *SetCommand {
	return &SetCommand{
		key: key,
		val: val,
	}
}
