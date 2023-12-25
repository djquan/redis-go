package main

import (
	"bufio"
	"net"
	"testing"
)

func Test_handleConnectionPing(t *testing.T) {
	conn1, conn2 := net.Pipe()
	defer conn1.Close()

	go handleConnection(conn1)

	_, err := conn2.Write([]byte("*1\r\n$4\r\nping\r\n"))
	if err != nil {
		t.Fatal(err)
	}

	reply, err := bufio.NewReader(conn2).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}

	if reply != "+PONG\r\n" {
		t.Errorf("Expected '+PONG\r\n', but got '%s'", reply)
	}
}
