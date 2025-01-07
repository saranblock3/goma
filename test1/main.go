package main

import (
	"fmt"
	"log"

	"github.com/saranblock3/goma"
)

func main() {
	homaSocket, err := goma.NewHomaSocket(3)
	if err != nil {
		log.Fatal(err)
	}
	defer homaSocket.Close()
	content := []byte("The quick brown fox jumps over the lazy dog")
	homaSocket.WriteTo(content, "130.127.133.61", 4)

	content, address, id, err := homaSocket.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("content: %s\naddress: %s\nid: %d\n", content, address, id)
}
