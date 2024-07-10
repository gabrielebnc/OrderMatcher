package transport

import (
	"fmt"
	"net"
)

// TCPTransport is responsible for:
//  - Listen for incoming connections
//  - Send request for creating connections
//  - Accept incoming connections
//  - Handle connections
//  - Send messages

type TCPTransportConfigs struct {
	listenAddr string
	//port       string
}

type TCPTransport struct {
	configs TCPTransportConfigs
	ln      net.Listener
	quitch  chan struct{}
}

func NewTCPTransport(configs TCPTransportConfigs) *TCPTransport {
	return &TCPTransport{
		configs: configs,
		quitch:  make(chan struct{}),
	}
}

func NewTCPTransportConfigs(listenAddr string) TCPTransportConfigs {
	return TCPTransportConfigs{
		listenAddr: listenAddr,
		//port:       port,
	}
}

func (tcpt *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", tcpt.configs.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	tcpt.ln = ln

	go tcpt.acceptLoop()

	<-tcpt.quitch

	return nil
}

func (tcpt *TCPTransport) acceptLoop() {
	for {
		conn, err := tcpt.ln.Accept()
		if err != nil {
			fmt.Println("ERROR (ACCEPT): ", err)
			continue
		}
		go tcpt.readLoop(conn)
	}
}

func (tcpt *TCPTransport) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("ERROR (READLOOP): ", err)
			continue
		}

		msg := buf[:n]
		fmt.Println("MSG: ", string(msg))
	}
}
