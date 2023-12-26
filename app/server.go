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
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		handleConnectionHelper(conn, reader)
	}
}

func handleConnectionHelper(conn net.Conn, reader *bufio.Reader) {
	msg, err := reader.ReadByte()
	if err != nil {
		if err != io.EOF {
			log.Printf("Failed to read from socket: %v", err)
		}

		return
	}

	switch msg {
	case '*':
		parseArray(conn, reader)
	case '$':
		parseBulkString(conn, reader)
	}
}

func parseBulkString(conn net.Conn, reader *bufio.Reader) {
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

	buf := make([]byte, length)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	reader.ReadByte()
	reader.ReadByte()

	handleCommand(conn, string(buf), reader)
}

func handleCommand(conn net.Conn, s string, reader *bufio.Reader) {
	println(s)
	switch strings.ToUpper(s) {
	case "PING":
		conn.Write([]byte("+PONG\r\n"))
	case "ECHO":

		echo, err := ParseBulkString(reader)
		println(echo)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		conn.Write([]byte("+"))
		conn.Write([]byte(echo))
		conn.Write([]byte("\r\n"))
	}
}

func parseArray(conn net.Conn, reader *bufio.Reader) {
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
		handleConnectionHelper(conn, reader)
	}
}
