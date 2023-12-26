package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

type value struct {
	v          []byte
	expiration *int64
}

type server struct {
	l net.Listener

	db *struct {
		sync.RWMutex
		m map[string]*value
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
					conn.Write(encodeString(echo))

					input = input[2:]
				case "SET":
					key := input[1].(string)

					v := &value{
						v: []byte(input[2].(string)),
					}

					input = input[3:]

					if len(input) >= 1 {
						if val, ok := input[0].(string); ok {
							if strings.ToUpper(val) == "PX" {
								now := time.Now().UnixNano() / int64(time.Millisecond)
								expiration, err := strconv.Atoi(input[1].(string))
								exp := int64(expiration) + now

								if err != nil {
									fmt.Println("Error getting expiration")
									os.Exit(1)
								}
								v.expiration = &exp

								input = input[2:]
							}
						}
					}

					s.db.Lock()
					s.db.m[key] = v
					s.db.Unlock()

					conn.Write([]byte("+OK\r\n"))

				case "GET":
					key := input[1].(string)
					now := time.Now().UnixNano() / int64(time.Millisecond)

					s.db.Lock()
					val := s.db.m[key]

					if val != nil && val.expiration != nil && *val.expiration <= now {
						s.db.m[key] = nil
						val = nil
					}

					s.db.Unlock()

					if val == nil {
						conn.Write([]byte("$-1\r\n"))
					} else {
						conn.Write(encodeString(string(val.v)))
					}

					input = input[2:]
				}
			default:

			}
		}
	}
}

func encodeString(s string) []byte {
	buf := new(bytes.Buffer)
	body := []byte(s)

	buf.WriteByte('$')

	length := strconv.Itoa(len(body))

	buf.WriteString(length)
	buf.Write([]byte("\r\n"))

	buf.Write(body)
	buf.Write([]byte("\r\n"))

	return buf.Bytes()
}

func startServer(port string) (*server, error) {
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		return nil, err
	}

	var db = struct {
		sync.RWMutex
		m map[string]*value
	}{m: make(map[string]*value)}

	return &server{
		l:  l,
		db: &db,
	}, nil
}
