package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/saranblock3/goma"
)

func homaListener() {

	id, _ := strconv.Atoi(os.Args[1])

	homaSocket, err := goma.NewHomaSocket(uint32(id))
	if err != nil {
		log.Fatal("Socket error:", err)
	}
	defer homaSocket.Close()

	//content, address, id, err := homaSocket.Read()
	//if err != nil {
	//	log.Fatal("Read error:", err)
	//}
	//fmt.Printf("content: %s\naddress: %s\nid: %d\n", content, address, id)
	i := 0

	for {
		i++
		content, address, id, messageId, err := homaSocket.Read()
		if err != nil {
			log.Fatal("Read error:", err)
		}
		fmt.Printf("content: %s\naddress: %s\nsource id: %d\nmessage id: %d\n", content, address, id, messageId)
		fmt.Printf("-+- %d\n", messageId)

		x := fmt.Sprintf("%d", messageId)

		response := []byte("I think therefore I am " + strconv.Itoa(len(content)) + " " + x)
		homaSocket.WriteTo(response, address, id)
		fmt.Println(i)
	}

}

func main() {
	homaListener()
}
