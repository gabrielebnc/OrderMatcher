package transport

import (
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
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
}

type TCPTransport struct {
	configs TCPTransportConfigs
	ln      net.Listener
	quitch  chan struct{}
	Msgch   chan Message

	connsmapMu sync.RWMutex
	connsmap   map[string]net.Conn
}

func NewTCPTransportConfigs(listenAddr string) TCPTransportConfigs {
	return TCPTransportConfigs{
		listenAddr: listenAddr,
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

		tcpt.addConnecton(conn.RemoteAddr().String(), conn)
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

			switch {
			case errors.Is(err, io.EOF):
				tcpt.removeConnecton(conn.RemoteAddr().String())
				fmt.Println("Conn closed (EOF):", conn.RemoteAddr().String())
				return
			case errors.Is(err, net.ErrClosed):
				tcpt.removeConnecton(conn.RemoteAddr().String())
				fmt.Println("Conn closed (ErrClosed):", conn.RemoteAddr().String())
				return
			case errors.As(err, new(*net.OpError)):
				tcpt.removeConnecton(conn.RemoteAddr().String())
				fmt.Println("Conn closed (OpError):", conn.RemoteAddr().String())
				return
			default:
				fmt.Println(reflect.TypeOf(err))
				continue
			}
		}

		tcpt.Msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}
	}
}

func (tcpt *TCPTransport) addConnecton(address string, conn net.Conn) {
	tcpt.connsmapMu.Lock()
	tcpt.connsmap[address] = conn
	tcpt.connsmapMu.Unlock()
}

func (tcpt *TCPTransport) removeConnecton(address string) {
	tcpt.connsmapMu.Lock()
	delete(tcpt.connsmap, address)
	tcpt.connsmapMu.Unlock()
}

func (tcpt *TCPTransport) CloseConnection(address string) {
	tcpt.connsmap[address].Close() //concurrency safety?
	tcpt.removeConnecton(address)
}

func (tcpt *TCPTransport) Dial(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("ERROR (DIAL):", err)
	}
	tcpt.addConnecton(address, conn) //what if there's already a connection ?
}

func (tcpt *TCPTransport) SendMessage(address string, message []byte) {
	fmt.Println(tcpt.connsmap)
	conn, ok := tcpt.connsmap[address]
	if ok {
		fmt.Println("Sending message to", address)
		conn.Write(message)
	} else {
		fmt.Printf("ERROR (SEND): %s is not connected\n", address)
	}
}
