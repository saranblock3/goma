package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/saranblock3/goma"
)

const NUM_MESSAGES = 10

var mu sync.Mutex
var latenciesMap map[uint64]time.Time = make(map[uint64]time.Time)

var id uint32

var latenciesSlice []int64 = make([]int64, 0)

type host struct {
	Address string   `json:"address"`
	Ids     []uint32 `json:"ids"`
}

func setup() (string, map[string]host, []byte, error) {
	hosts := make(map[string]host)

	data, err := os.ReadFile("config.json")
	if err != nil {
		return "", nil, nil, err
	}

	err = json.Unmarshal(data, &hosts)
	if err != nil {
		return "", nil, nil, err
	}

	content, err := os.ReadFile("content.txt")
	if err != nil {
		return "", nil, nil, err
	}

	if len(os.Args) != 2 {
		return "", nil, nil, fmt.Errorf("invalid arguments")
	}

	localAddress := hosts[os.Args[1]].Address

	return localAddress, hosts, content, err
}

func homaClient(localAddress string, localId uint32, hosts map[string]host, content []byte) {
	homaSocket, err := goma.NewHomaSocket(localId)
	if err != nil {
		log.Fatal(err)
	}
	defer homaSocket.Close()

	wg1 := &sync.WaitGroup{}

	for _, remoteHost := range hosts {
		for _, remoteId := range remoteHost.Ids {
			for range NUM_MESSAGES {
				wg1.Add(1)
				go func(remoteHost host, remoteId uint32) {
					id := rand.Uint64()
					buf := make([]byte, 8)
					binary.LittleEndian.PutUint64(buf, id)
					content = append(buf, content...)
					err := homaSocket.SendTo(content, localAddress, remoteHost.Address, remoteId)
					start := time.Now()
					mu.Lock()
					latenciesMap[id] = start
					mu.Unlock()
					if err != nil {
						log.Fatal(err)
					}
					wg1.Done()
				}(remoteHost, remoteId)
				wg1.Add(1)
				go func() {
					content, _, _, _, err := homaSocket.RecvFrom()
					if err != nil {
						log.Fatal(err)
					}
					id := binary.LittleEndian.Uint64(content[:8])
					end := time.Now()
					mu.Lock()
					latency := (end.Sub(latenciesMap[id])).Nanoseconds()
					latenciesSlice = append(latenciesSlice, latency)
					mu.Unlock()
					wg1.Done()
				}()
			}
		}
	}
	ch := make(chan bool)
	go func() {
		wg1.Wait()
		ch <- true
	}()
	select {
	case <-ch:
	case <-time.After(60 * time.Second):
	}
	wg1.Wait()
}

func main() {
	localAddress, hosts, content, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	var i uint32 = 1000
	homaClient(localAddress, i, hosts, content)
	latencies, err := os.OpenFile("latencies.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer latencies.Close()
	mu.Lock()
	for _, latency := range latenciesSlice {
		latencies.Write([]byte(strconv.FormatInt(latency, 10) + "\n"))
	}
	mu.Unlock()
}
