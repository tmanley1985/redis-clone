package main

const (
	CommandSet = "SET"
)

type Command interface {
}

type SetCommand struct {
	key, val string
}
