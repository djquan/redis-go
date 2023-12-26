package main

import (
	"bufio"
	"strings"
	"testing"
)

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
