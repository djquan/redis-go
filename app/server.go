package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type server struct {
	l net.Listener

	db *struct {
		sync.RWMutex
		m map[string][]byte
	}
}

func (s *server) listenAndServe() {
	for {
		conn, err := s.l.Accept()
		fmt.Println("accepted new connection")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		go s.handleConn(conn)
	}
}

func (s *server) handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		s.handleConnectionHelper(conn, reader)
	}
}

func (s *server) handleConnectionHelper(conn net.Conn, reader *bufio.Reader) {
	msg, err := reader.ReadByte()
	if err != nil {
		if err != io.EOF {
			log.Printf("Failed to read from socket: %v", err)
		}

		return
	}

	switch msg {
	case '*':
		s.parseArray(conn, reader)
	case '$':
		reader.UnreadByte()
		s.handleBulkString(conn, reader)
	}
}

func startServer(port string) (*server, error) {
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return nil, err
	}

	var db = struct {
		sync.RWMutex
		m map[string][]byte
	}{m: make(map[string][]byte)}

	return &server{
		l:  l,
		db: &db,
	}, nil
}

func main() {
	server, err := startServer("6379")
	defer server.l.Close()
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	server.listenAndServe()
}

func (s *server) handleBulkString(conn net.Conn, reader *bufio.Reader) {
	command, err := ParseBulkString(reader)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	s.handleCommand(conn, command, reader)
}

func (s *server) handleCommand(conn net.Conn, command string, reader *bufio.Reader) {
	switch strings.ToUpper(command) {
	case "PING":
		conn.Write([]byte("+PONG\r\n"))
	case "ECHO":
		echo, err := ParseBulkString(reader)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		conn.Write([]byte("+"))
		conn.Write([]byte(echo))
		conn.Write([]byte("\r\n"))
	case "SET":
		key, err := ParseBulkString(reader)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		value, err := ParseBulkString(reader)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		s.db.Lock()
		s.db.m[key] = []byte(value)
		s.db.Unlock()

		conn.Write([]byte("+OK\r\n"))
	case "GET":
		key, err := ParseBulkString(reader)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		s.db.RLock()
		val := s.db.m[key]
		s.db.RUnlock()

		conn.Write([]byte("+"))
		conn.Write(val)
		conn.Write([]byte("\r\n"))
	}
}

func (s *server) parseArray(conn net.Conn, reader *bufio.Reader) {
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	length, err := strconv.Atoi(strings.TrimSuffix(line, "\r\n"))

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for i := 0; i < length; i++ {
		s.handleConnectionHelper(conn, reader)
	}
}
