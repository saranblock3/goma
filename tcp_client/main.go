package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

func abortAfter(n time.Duration) {
	time.Sleep(n * time.Second)
	os.Exit(0)
}

func setup() (uint32, map[string]host, []byte, error) {
	hosts := make(map[string]host)

	data, err := os.ReadFile("config.json")
	if err != nil {
		return 0, nil, nil, err
	}

	err = json.Unmarshal(data, &hosts)
	if err != nil {
		return 0, nil, nil, err
	}

	content, err := os.ReadFile("content.txt")
	if err != nil {
		return 0, nil, nil, err
	}

	if len(os.Args) != 2 {
		return 0, nil, nil, fmt.Errorf("invalid arguments")
	}

	localIdInt, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return 0, nil, nil, err
	}
	localId := uint32(localIdInt)

	return localId, hosts, content, err
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
	localId, hosts, content, err := setup()
	if err != nil {
		log.Fatal(err)
	}

	go abortAfter(10)

	wg := sync.WaitGroup{}

	for _, remoteHost := range hosts {
		for _, remotePort := range remoteHost.Ports {
			for range NUM_MESSAGES {
				wg.Add(1)
				go func(remoteHost host, remotePort uint32) {
					conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", remoteHost.Address, remotePort))
					if err != nil {
						log.Fatal(err)
					}
					defer conn.Close()

					err = writeAll(conn, content)
					if err != nil {
						log.Fatal(err)
					}
					start := time.Now()

					_, err = readAll(conn)
					if err != nil {
						log.Fatal(err)
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
	wg.Wait()
	latencies, err := os.OpenFile(fmt.Sprintf("latencies_%d.txt", localId), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer latencies.Close()
	for _, latency := range latenciesSlice {
		_, err := fmt.Fprintf(latencies, "%d\n", latency)
		if err != nil {
			log.Fatal(err)
		}
	}
}
