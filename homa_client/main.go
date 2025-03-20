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

func abortAfter(n time.Duration) {
	time.Sleep(n * time.Second)
	latencies, err := os.OpenFile(fmt.Sprintf("latencies_%d.txt", id), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer latencies.Close()
	for _, latency := range latenciesSlice {
		latencies.Write([]byte(strconv.FormatInt(latency, 10) + "\n"))
	}
	os.Exit(0)
}

func setup() (string, uint32, map[string]host, []byte, error) {
	hosts := make(map[string]host)

	data, err := os.ReadFile("config.json")
	if err != nil {
		return "", 0, nil, nil, err
	}

	err = json.Unmarshal(data, &hosts)
	if err != nil {
		return "", 0, nil, nil, err
	}

	content, err := os.ReadFile("content.txt")
	if err != nil {
		return "", 0, nil, nil, err
	}

	if len(os.Args) != 3 {
		return "", 0, nil, nil, fmt.Errorf("invalid arguments")
	}

	localAddress := hosts[os.Args[1]].Address
	localIdInt, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return "", 0, nil, nil, err
	}
	localId := uint32(localIdInt)

	return localAddress, localId, hosts, content, err
}

func main() {
	localAddress, localId, hosts, content, err := setup()
	if err != nil {
		log.Fatal(err)
	}
	id = localId

	homaSocket, err := goma.NewHomaSocket(localId)
	if err != nil {
		log.Fatal(err)
	}
	defer homaSocket.Close()

	go abortAfter(10)

	wg := sync.WaitGroup{}

	for _, remoteHost := range hosts {
		for _, remoteId := range remoteHost.Ids {
			for range NUM_MESSAGES {
				wg.Add(1)
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
					wg.Done()
				}(remoteHost, remoteId)
				wg.Add(1)
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
					wg.Done()
				}()
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
		latencies.Write([]byte(strconv.FormatInt(latency, 10) + "\n"))
	}
}
