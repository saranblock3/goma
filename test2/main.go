package main

import (
	"fmt"
	"log"

	"github.com/saranblock3/goma"
)

func main() {
	homaSocket, err := goma.NewHomaSocket(1)
	if err != nil {
		log.Fatal(err)
	}
	defer homaSocket.Close()
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
			"The quick brown fox jumps over the lazy dog\n",
	)
	homaSocket.WriteTo(content, "127.0.0.1", 3)

	content, address, id, err := homaSocket.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("content: %s\naddress: %s\nid: %d\n", content, address, id)
}
