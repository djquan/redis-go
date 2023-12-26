package main

import (
	"bufio"
	"net"
	"testing"
)

func Test_handleConnectionPing(t *testing.T) {
	conn1, conn2 := net.Pipe()
	defer conn1.Close()
	defer conn2.Close()

	s, err := startServer("0")
	if err != nil {
		t.Fatal(err)
	}

	go s.handleConn(conn1)

	_, err = conn2.Write([]byte("*1\r\n$4\r\nping\r\n"))

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

func Test_handleEcho(t *testing.T) {
	conn1, conn2 := net.Pipe()

	defer conn1.Close()
	defer conn2.Close()

	s, err := startServer("0")
	if err != nil {
		t.Fatal(err)
	}

	go s.handleConn(conn1)

	_, err = conn2.Write([]byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"))
	if err != nil {
		t.Fatal(err)
	}

	reply, err := bufio.NewReader(conn2).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}

	if reply != "+hey\r\n" {
		t.Errorf("Expected '+hey\r\n', but got '%s'", reply)
	}
}

func Test_handleSetGet(t *testing.T) {
	conn1, conn2 := net.Pipe()
	defer conn1.Close()
	defer conn2.Close()

	s, err := startServer("0")
	if err != nil {
		t.Fatal(err)
	}

	go s.handleConn(conn1)

	_, err = conn2.Write([]byte("*3\r\n$3\r\nSET\r\n$3\r\nhey\r\n$3\r\nbye\r\n"))
	if err != nil {
		t.Fatal(err)
	}

	reply, err := bufio.NewReader(conn2).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}

	if reply != "+OK\r\n" {
		t.Errorf("Expected '+OK\r\n', but got '%s'", reply)
	}

	_, err = conn2.Write([]byte("*2\r\n$3\r\nGET\r\n$3\r\nhey\r\n"))
	if err != nil {
		t.Fatal(err)
	}

	reply, err = bufio.NewReader(conn2).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}

	if reply != "+bye\r\n" {
		t.Errorf("Expected '+bye\r\n', but got '%s'", reply)
	}
}
