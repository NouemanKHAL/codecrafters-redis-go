package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var CRLF = []byte{'\r', '\n'}

const (
	KB          = (1 << 10)
	MB          = (KB << 10)
	BUFFER_SIZE = 512 * KB
)

// RESP Types
const (
	SIMPLE_STRING = '+'
	ERROR         = '-'
	INTEGER       = ':'
	BULK_STRING   = '$'
	ARRAY         = '*'
)

func isRESPType(r rune) bool {
	return r == SIMPLE_STRING ||
		r == ERROR ||
		r == INTEGER ||
		r == BULK_STRING ||
		r == ARRAY
}

func encodeRESPSimpleString(s string) string {
	return fmt.Sprintf("%v%s\r\n", SIMPLE_STRING, s)
}

func encodeRESPBulkString(ss []string) string {
	res := strings.Join(ss, "")
	return fmt.Sprintf("%v%d%v%s%v", BULK_STRING, len(res), CRLF, res, CRLF)
}

func decodeRESPArray(b []byte) ([]string, error) {
	if b[0] != ARRAY {
		return nil, fmt.Errorf("invalid argument, RESP Array must start with %v", ARRAY)
	}

	curr := 0
	nextDelim := bytes.Index(b[curr:], CRLF)

	if curr+1 > curr+nextDelim {
		return nil, fmt.Errorf("invalid argument, check RESP Array syntax %d %d", curr+1, curr+nextDelim)
	}
	size, err := strconv.Atoi(string(b[curr+1 : curr+nextDelim]))
	if err != nil {
		return nil, err
	}

	curr += nextDelim + 2

	res := make([]string, size)

	for i := 0; i < size; i++ {
		var data []byte
		nextDelim := bytes.Index(b[curr:], CRLF)
		switch b[curr] {
		case BULK_STRING:
			strSize, _ := strconv.Atoi(string(b[curr+1 : curr+nextDelim]))
			data = b[curr+nextDelim+2 : curr+nextDelim+2+strSize]
			curr += nextDelim + 2 + strSize + 2
		case INTEGER:
			data = b[curr+1 : curr+nextDelim]
			curr += nextDelim + 2
		}
		res[i] = string(data)
	}

	fmt.Println(res)
	return res, nil
}

func pingCommand() string {
	return encodeRESPSimpleString("PONG")
}

func echoCommand(in []string) string {
	return encodeRESPBulkString(in)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, BUFFER_SIZE)
	for {
		if _, err := conn.Read(buf); err != nil {
			fmt.Println("Error reading from connection: ", err.Error())
			continue
		}
		tokens, err := decodeRESPArray(buf)
		fmt.Println("received request: ", tokens)
		if err != nil {
			log.Printf("error decoding request: %v", err)
			continue
		}
		var response string
		switch tokens[0] {
		case "PING":
			response = pingCommand()
		case "ECHO":
			response = echoCommand(tokens[1:])
		}
		conn.Write([]byte(response))
	}
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}
