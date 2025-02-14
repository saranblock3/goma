package main

import (
	"fmt"
	"net"
	"os"
)

func tcpListener() {
	address := "localhost:8080"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer conn.Close()
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil || n == 0 {
				return
			}
			conn.Write([]byte("I think therefore I am"))
		}(conn)

	}

}

func main() {
	tcpListener()
}
