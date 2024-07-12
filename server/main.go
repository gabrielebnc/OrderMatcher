package main

import (
	"fmt"

	"github.com/gabrielebnc/OrderMatcher/core/transport"
)

func main() {

	tcp_configs := transport.NewTCPTransportConfigs(":3000")

	server := transport.NewTCPTransport(tcp_configs)

	go func() {
		for msg := range server.Msgch {
			fmt.Printf("msg from %v: %v\n", msg.From(), string(msg.Payload()))
		}
	}()

	server.Start()

}
