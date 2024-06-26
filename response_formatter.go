package main

// ResponseCommand represents a Redis server response
type ClientResponse interface {
	// Type returns the RESP type of the response (SimpleString, Error, Integer, BulkString, Array)
	// TODO: Figure out how to do enums or simple unions in go if possible.
	Type() string

	// Serialize converts the response to RESP format byte slice
	Serialize() []byte
}

// TODO: Put these into their own files.
// TODO: Add more responses.

type SimpleStringResponse struct {
	value string
}

func NewSimpleStringResponse(value string) *SimpleStringResponse {
	return &SimpleStringResponse{value: value}
}

func (r *SimpleStringResponse) Type() string {
	return "SimpleString"
}

func (r *SimpleStringResponse) Serialize() []byte {
	return []byte("+" + r.value + "\r\n")
}

type ErrorResponse struct {
	error string
}

func NewErrorResponse(error string) *ErrorResponse {
	return &ErrorResponse{error}
}

func (r *ErrorResponse) Type() string {
	return "Error"
}

func (r *ErrorResponse) Serialize() []byte {
	return []byte("-ERR" + r.error + "\r\n")
}
