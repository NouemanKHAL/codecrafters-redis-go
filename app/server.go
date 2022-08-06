package main

import (
	"fmt"
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

	for {
		conn, err := l.Accept()
		go func() {
			defer conn.Close()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}

			for {
				if _, err := conn.Read([]byte{}); err != nil {
					fmt.Println("Error reading from connection: ", err.Error())
					continue
				}
				conn.Write([]byte(toRestSimpleString("PONG")))
			}
		}()
	}
}
