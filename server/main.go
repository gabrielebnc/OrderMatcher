package main

import (
	"fmt"

	"github.com/gabrielebnc/OrderMatcher/core/transport"
)

func main() {

	tcp_configs := transport.NewTCPTransportConfigs(":3000")

	server := transport.NewTCPTransport(tcp_configs)

	server.Start()
	fmt.Println("ok server")
}
