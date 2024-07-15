package main

import (
	"fmt"
	"time"

	"github.com/gabrielebnc/OrderMatcher/core/transport"
)

func main() {

	tcp_configs := transport.NewTCPTransportConfigs(":42068")

	client := transport.NewTCPTransport(tcp_configs)

	client.Dial(":3000")

	go func() {
		for msg := range client.Msgch {
			fmt.Printf("msg from %v: %v\n", msg.From(), string(msg.Payload()))
		}
	}()

	client.SendMessage(":3000", []byte("ciao"))

	time.Sleep(time.Second * 10)

	client.CloseConnection(":3000")
}
