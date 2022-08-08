package main

import (
	"bufio"
	"build-your-own-redis/app/resp"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		req, err := resp.Decode(bufio.NewReader(conn))
		if err != nil {
			if err == io.EOF {
				continue
			}
			log.Println("ERROR: ", err.Error())
			continue
		}

		cmd := req.Array()[0].String()
		args := req.Array()[1:]

		var response []byte

		switch strings.ToUpper(cmd) {
		case "PING":
			response = resp.Ping()
		case "ECHO":
			response = resp.Echo(args[0])
		}

		conn.Write(response)
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
