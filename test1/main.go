package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/saranblock3/goma"
)

const NUM_MESSAGES = 10

func main() {
	content := []byte(
		"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n" +
			"The quick brown fox jumps over the lazy dog\n",
	)

	// id := rand.Uint32()
	id := rand.Uint32()
	if len(os.Args) > 1 {
		j, _ := strconv.Atoi(os.Args[1])
		id = uint32(j)
	}

	homaSocket, err := goma.NewHomaSocket(id)
	if err != nil {
		log.Fatal(err)
	}

	// mu := sync.Mutex{}
	mu1 := sync.Mutex{}
	counter := 0

	time.Sleep(1 * time.Second)

	wg := sync.WaitGroup{}
	for i := 500; i < 500+NUM_MESSAGES; i++ {
		// go func() {
		time.Sleep(time.Duration(i * 100))
		start := time.Now()
		// start := time.Now()
		contentLength := len(content) - int(float32(len(content))*float32(i-500)/float32(NUM_MESSAGES))
		// contentLength := len(content)
		// mu.Lock()
		destination_id := uint32(rand.Intn(500) + 500)
		fmt.Println("====", destination_id)
		err = homaSocket.WriteTo(content[:contentLength], "127.0.0.1", destination_id)
		if err != nil {
			log.Fatal(err)
		}
		// mu.Unlock()

		wg.Add(1)
		go func() {
			mu1.Lock()
			content, address, id, messageId, err := homaSocket.Read()
			if err != nil {
				log.Fatal(err)
			}
			counter++
			mu1.Unlock()
			fmt.Printf("content: %s\naddress: %s\nsource id: %d\nmessage id: %d %d --- %d\n", content, address, id, messageId, len(content), i)
			end := time.Since(start)
			fmt.Printf("%dms\n", end.Milliseconds())
			wg.Done()
		}()
		// }()
	}
	wg.Wait()
	err = homaSocket.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("++++++", counter)
	// end := time.Since(start)
	// avg := end.Milliseconds() / 50
	// fmt.Printf("%dms\n", avg)
}
