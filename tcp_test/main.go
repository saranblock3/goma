package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// Define the CIDR
	if len(os.Args) < 2 {
		log.Fatal("Usage: tcp_test <ip>")
	}
	addresses := os.Args[1:]

	// Parse the CIDR
	var counter float64 = 0
	var total int64 = 0

	// Iterate over all IPs in the subnet
	for i := 0; i < 100000; i++ {
		for _, address := range addresses {
			start := time.Now()
			conn, err := net.Dial("tcp", address+":8080")
			if err != nil {
				fmt.Println("Dial failed!")
				continue
			}
			_, err = conn.Write([]byte("The quick brown fox jumps over the lazy dog"))
			if err != nil {
				fmt.Println("Write failed!")
				continue
			}
			readChan := make(chan int)
			go func() {
				buf := make([]byte, 1024)
				n, _ := conn.Read(buf)
				fmt.Println(string(buf[:n]))
				readChan <- n
			}()
			select {
			case n := <-readChan:
				if n == 0 {
					continue
				}
			case <-time.After(time.Millisecond * 300):
				fmt.Println("Timeout!")
				continue
			}
			conn.Close()
			counter += 1
			total += time.Since(start).Microseconds()
		}
	}

	fmt.Println("No. connections: ", counter)
	fmt.Println("Average RTT (microseconds): ", float64(total)/counter)
}
