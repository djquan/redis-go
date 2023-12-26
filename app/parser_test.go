package main

import (
	"bufio"
	"strings"
	"testing"
)

func TestParseString(t *testing.T) {
	ping := bufio.NewReader(strings.NewReader("*1\r\n$4\r\nping\r\n"))

	result, err := Parse(ping)
	if err != nil {
		t.Errorf("got error %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected one command, got %v", len(result))
	}

	command := result[0]

	if s, ok := command.(string); ok {
		if s != "ping" {
			t.Errorf("Expected ping, got %v", s)
		}
	} else {
		t.Errorf("Did not receive a string back")
	}
}

func TestParseBulkString(t *testing.T) {
	type args struct {
		reader *bufio.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Simple",
			args: args{
				reader: bufio.NewReader(strings.NewReader("$3\r\nhey\r\n")),
			},
			want:    "hey",
			wantErr: false,
		},
		{
			name: "Greater than 10",
			args: args{
				reader: bufio.NewReader(strings.NewReader("$12\r\nheyheyheyhey\r\n")),
			},
			want:    "heyheyheyhey",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBulkString(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBulkString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseBulkString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
