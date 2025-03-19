package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/saranblock3/goma"
)

func homaServer() {
	localIdInt, _ := strconv.Atoi(os.Args[1])
	localId := uint32(localIdInt)

	homaSocket, err := goma.NewHomaSocket(localId)
	if err != nil {
		log.Fatal("Socket error:", err)
	}
	defer homaSocket.Close()

	for {
		content, sourceAddress, destinationAddress, id, err := homaSocket.RecvFrom()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("received message", len(content), sourceAddress, id)

		response := append(content[:8], []byte("I think therefore I am")...)
		homaSocket.SendTo(response, destinationAddress, sourceAddress, id)
	}

}

func main() {
	homaServer()
}
