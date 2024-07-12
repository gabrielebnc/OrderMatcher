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

func (msg *Message) From() string {
	return msg.from
}

func (msg *Message) Payload() []byte {
	return msg.payload
}

type TCPTransportConfigs struct {
	listenAddr string
	//port       string
}

type TCPTransport struct {
	configs  TCPTransportConfigs
	ln       net.Listener
	quitch   chan struct{}
	Msgch    chan Message
	connsmap map[string]net.Conn // TODO concurrency safety?
}

func NewTCPTransportConfigs(listenAddr string) TCPTransportConfigs {
	return TCPTransportConfigs{
		listenAddr: listenAddr,
		//port:       port,
	}
}

func NewTCPTransport(configs TCPTransportConfigs) *TCPTransport {
	return &TCPTransport{
		configs:  configs,
		quitch:   make(chan struct{}),
		Msgch:    make(chan Message, 2048), // TODO should it be buffered or unbuffered? investigate
		connsmap: make(map[string]net.Conn),
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
	close(tcpt.Msgch)

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
		fmt.Println("Accepted Conn:", conn.RemoteAddr().String())

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

		tcpt.Msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}
	}
}

func (tcpt *TCPTransport) Dial(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("ERROR (DIAL):", err)
	}
	tcpt.connsmap[address] = conn //what if there's already a connection ?
}

func (tcpt *TCPTransport) CloseConnection(address string) {
	tcpt.connsmap[address].Close() //concurrency safety?
	delete(tcpt.connsmap, address)
}
