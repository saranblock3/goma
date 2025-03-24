package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func abortAfter(n time.Duration) {
	time.Sleep(n * time.Second)
	os.Exit(0)
}

func readAll(conn net.Conn) ([]byte, error) {
	sizeBuf := make([]byte, 8)
	read := 0
	for read < 8 {
		n, err := conn.Read(sizeBuf[read:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, io.EOF
		}
		read += n
	}

	size := binary.LittleEndian.Uint64(sizeBuf)

	buf := make([]byte, size)
	read = 0
	for read < int(size) {
		n, err := conn.Read(buf[read:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return nil, io.EOF
		}
		read += n
	}

	return buf, nil
}

func writeAll(conn net.Conn, buf []byte) error {
	sizeBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuf, uint64(len(buf)))
	written := 0
	for written < 8 {
		n, err := conn.Write(sizeBuf[written:])
		if err != nil {
			return err
		}
		written += n
	}

	written = 0
	for written < len(buf) {
		n, err := conn.Write(buf[written:])
		if err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

func tcpListener() {
	address := os.Args[1]
	port := os.Args[2]

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", address, port))
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
		go func(conn net.Conn) {
			defer conn.Close()
			_, err := readAll(conn)
			if err != nil {
				return
			}
			err = writeAll(conn, []byte("I think therefore I am"))
			if err != nil {
				return
			}
		}(conn)

	}

}

func main() {
	go abortAfter(600)
	tcpListener()
}
