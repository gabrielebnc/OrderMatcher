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
			//Messages may be aggregated, for this reason it will be necessary
			//to split them back properly

			fmt.Printf("msg from %v: %v\n", msg.From(), string(msg.Payload()))
			server.SendMessage(msg.From(), []byte("hello from server"))
		}
	}()

	server.Start()

}
