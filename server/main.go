package main

import (
	"fmt";
	"github.com/gabrielebnc/OrderMatcher/core/transport"
)

func main() {

	tcp_configs := transport.NewTCPTransportConfigs("3000")

	tcp_transport := transport.NewTCPTransport(tcp_configs)

	tcp_transport.Test()
	fmt.Println("ok server")
}
