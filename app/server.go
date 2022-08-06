package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func toRestSimpleString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	conn, err := l.Accept()
	defer conn.Close()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 256)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("Error reading from connection: ", err.Error())
		} else {
			conn.Write([]byte(toRestSimpleString("PONG")))
		}
	}
}
