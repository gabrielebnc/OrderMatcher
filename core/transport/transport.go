package transport

import (
	"fmt"
	"io"
	"net"
)

// TCPTransport is responsible for:
//  - Listen for incoming connections
//  - Send request for creating connections
//  - Accept incoming connections
//  - Handle connections
//  - Send messages

type Message struct {
	from    string
	payload []byte
}

type TCPTransportConfigs struct {
	listenAddr string
	//port       string
}

type TCPTransport struct {
	configs  TCPTransportConfigs
	ln       net.Listener
	quitch   chan struct{}
	msgch    chan Message
	connsmap map[string]net.Conn
}

func NewTCPTransport(configs TCPTransportConfigs) *TCPTransport {
	return &TCPTransport{
		configs:  configs,
		quitch:   make(chan struct{}),
		msgch:    make(chan Message, 2048), // TODO should it be buffered or unbuffered? investigate
		connsmap: make(map[string]net.Conn),
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
	close(tcpt.msgch)

	return nil
}

func (tcpt *TCPTransport) acceptLoop() {
	for {
		conn, err := tcpt.ln.Accept()
		if err != nil {
			fmt.Println("ERROR (ACCEPT):", err)
			continue
		}

		tcpt.connsmap[conn.RemoteAddr().String()] = conn
		fmt.Println("Incoming Connection:", conn)

		go tcpt.readLoop(conn)
	}
}

func (tcpt *TCPTransport) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("ERROR (READLOOP):", err.Error())

			switch err {
			case io.EOF:
				delete(tcpt.connsmap, conn.RemoteAddr().String())
				fmt.Println("Conn closed:", conn.RemoteAddr().String())
				return
			default:
				continue

			}
		}
		msg := buf[:n]

		tcpt.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		fmt.Println("UPSTANDING CONNS:", tcpt.connsmap)
		fmt.Println("MSG FROM:", conn)
		fmt.Println("MSG:", string(msg))
	}
}
