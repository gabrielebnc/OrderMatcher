package main

import (
	"fmt"

	"github.com/gabrielebnc/OrderMatcher/server/transport"
)

func main() {

	tcp_configs := *transport.NewTCPServerConfigs(":3000")

	server := *transport.NewTCPServer(tcp_configs)

	go func() {
		for msg := range server.Consume() {
			fmt.Printf("msg from %v: %v\n", msg.From(), string(msg.Payload()))
			server.SendMessage(msg.From(), []byte("ACK"))
		}
	}()

	server.Start()

}
