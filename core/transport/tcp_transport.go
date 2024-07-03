package transport

import (
	"fmt"
	"net"
)

// TCPTransport is a Transport
// It is responsible for:
//  - Listen for incoming connections
//  - Send request for creating connections
//  - Accept incoming connections
//  - Handle connections
//  - Send messages

type TCPTransportConfigs struct {
	port string
}

type TCPTransport struct {
	configs TCPTransportConfigs
}

func NewTCPTransport(configs TCPTransportConfigs) *TCPTransport {
	return &TCPTransport{
		configs: configs,
	}
}

func NewTCPTransportConfigs(port string) TCPTransportConfigs {
	return TCPTransportConfigs{
		port: port,
	}
}

func (tcp_t *TCPTransport) Test() {
	fmt.Println(net.InterfaceAddrs())
	fmt.Println(tcp_t)
}
