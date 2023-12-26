package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	server, err := startServer("6379")
	defer server.l.Close()
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	server.listenAndServe()
}

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
		input, err := Parse(reader)
		if err != nil {
			fmt.Printf("Got %v\n", err)
		}

		for len(input) != 0 {
			command := input[0]

			fmt.Println(command)
			switch v := command.(type) {
			case string:
				switch strings.ToUpper(v) {
				case "PING":
					conn.Write([]byte("+PONG\r\n"))
					input = input[1:]
				case "ECHO":
					echo := input[1].(string)

					conn.Write([]byte("+"))
					conn.Write([]byte(echo))
					conn.Write([]byte("\r\n"))

					input = input[2:]
				case "SET":
					key := input[1].(string)
					value := input[2].(string)

					s.db.Lock()
					s.db.m[key] = []byte(value)
					s.db.Unlock()

					conn.Write([]byte("+OK\r\n"))

					input = input[3:]

				case "GET":
					key := input[1].(string)

					s.db.RLock()
					val := s.db.m[key]
					s.db.RUnlock()

					conn.Write([]byte("+"))
					conn.Write(val)
					conn.Write([]byte("\r\n"))

					input = input[2:]
				}
			default:

			}
		}
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
