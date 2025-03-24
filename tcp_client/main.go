package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const NUM_MESSAGES = 10

var mu sync.Mutex
var latenciesSlice []int64 = make([]int64, 0)

type host struct {
	Address string   `json:"address"`
	Ports   []uint32 `json:"ports"`
}

func setup() (map[string]host, []byte, error) {
	hosts := make(map[string]host)

	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(data, &hosts)
	if err != nil {
		return nil, nil, err
	}

	content, err := os.ReadFile("content.txt")
	if err != nil {
		return nil, nil, err
	}

	if len(os.Args) != 2 {
		return nil, nil, fmt.Errorf("invalid arguments")
	}

	num, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return nil, nil, err
	}

	exp := int(math.Pow(1.1, float64(num)) * float64(len(content)))

	for i := 0; i < exp; i++ {
		content = append(content, 'x')
	}

	return hosts, content, err
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

func main() {
	hosts, content, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(os.Stderr, len(content))

	wg := sync.WaitGroup{}

	for _, remoteHost := range hosts {
		for _, remotePort := range remoteHost.Ports {
			for range NUM_MESSAGES {
				wg.Add(1)
				go func(remoteHost host, remotePort uint32) {
					conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", remoteHost.Address, remotePort))
					if err != nil {
						return
					}
					defer conn.Close()

					err = writeAll(conn, content)
					if err != nil {
						return
					}

					start := time.Now()

					_, err = readAll(conn)
					if err != nil {
						return
					}

					latency := time.Since(start).Nanoseconds()
					mu.Lock()
					latenciesSlice = append(latenciesSlice, latency)
					mu.Unlock()
					wg.Done()
				}(remoteHost, remotePort)
			}
		}
	}
	ch := make(chan bool)
	go func() {
		wg.Wait()
		ch <- true
	}()
	select {
	case <-ch:
	case <-time.After(20 * time.Second):
	}
	mu.Lock()
	for _, latency := range latenciesSlice {
		fmt.Printf("%d\n", latency)
	}
	mu.Unlock()
}
