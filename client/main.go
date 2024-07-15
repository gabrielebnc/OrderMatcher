package main

import (
	"fmt"
	"time"

	"github.com/gabrielebnc/OrderMatcher/client/transport"
)

func main() {
	serverAddr := ":3000"

	client := transport.NewTCPClient()

	client.Dial(serverAddr)

	go func() {
		for msg := range client.Consume() {
			fmt.Printf("msg from %v: %v\n", msg.From(), string(msg.Payload()))
		}
	}()

	client.SendMessage(serverAddr, []byte("eskere"))
	time.Sleep(500 * time.Millisecond)
	client.SendMessage(serverAddr, []byte("ciaone"))
	time.Sleep(10 * time.Second)
	client.CloseConnection(serverAddr)
}
