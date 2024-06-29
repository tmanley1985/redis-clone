package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    RequestCommand
		wantErr bool
	}{
		{
			name:    "valid SET command",
			input:   "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$5\r\nHello\r\n",
			want:    &SetCommand{key: "mykey", val: "Hello"},
			wantErr: false,
		},
		{
			name:    "valid GET command",
			input:   "*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n",
			want:    &GetCommand{key: "mykey"},
			wantErr: false,
		},
		{
			name:    "invalid command",
			input:   "*1\r\n$7\r\nINVALID\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "malformed input",
			input:   "NOT A VALID RESP INPUT",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := ParseCommand(reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
