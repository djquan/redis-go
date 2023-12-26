package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func Parse(reader *bufio.Reader) ([]interface{}, error) {
	return parseArray(reader)
}

func parseHelper(reader *bufio.Reader) (interface{}, error) {
	msg, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	reader.UnreadByte()

	switch msg {
	case '*':
		return parseArray(reader)
	case '$':
		return ParseBulkString(reader)
	default:
		return nil, fmt.Errorf("unhandled RESP")
	}
}

func parseArray(reader *bufio.Reader) ([]interface{}, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	length, err := strconv.Atoi(strings.TrimSuffix(line[1:], "\r\n"))
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, length)

	for i := 0; i < length; i++ {
		r, err := parseHelper(reader)
		if err != nil {
			return nil, err
		}

		result[i] = r
	}

	return result, nil
}

func ParseBulkString(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	length, err := strconv.Atoi(strings.TrimSuffix(line, "\r\n")[1:])

	if err != nil {
		return "", err
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	reader.ReadByte()
	reader.ReadByte()

	return string(buf), nil
}
