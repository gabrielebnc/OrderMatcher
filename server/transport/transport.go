package transport

import (
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"

	"github.com/gabrielebnc/OrderMatcher/commons/core"
)

// TCPServer is responsible for:
//  - Listen for incoming connections
//  - Accept incoming connections
//  - Handle its connections
//  - Send messages

type TCPServerConfig struct {
	listenAddr string
}

type TCPServer struct {
	configs TCPServerConfig
	ln      net.Listener
	quitch  chan struct{}
	msgch   chan core.Message

	connsMapMu sync.RWMutex
	connsMap   map[string]net.Conn
}

func NewTCPServerConfigs(listenAddr string) *TCPServerConfig {
	return &TCPServerConfig{
		listenAddr: listenAddr,
	}
}

func NewTCPServer(configs TCPServerConfig) *TCPServer {
	return &TCPServer{
		configs:  configs,
		quitch:   make(chan struct{}),
		msgch:    make(chan core.Message, 512), // TODO should it be buffered or unbuffered?
		connsMap: make(map[string]net.Conn),
	}
}

func (s *TCPServer) Start() error {
	ln, err := net.Listen("tcp", s.configs.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *TCPServer) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("ERROR (ACCEPT):", err)
			continue
		}

		s.addConnecton(conn.RemoteAddr().String(), conn)
		fmt.Println("Accepted Conn:", conn.RemoteAddr().String())

		go s.readLoop(conn)
	}
}

func (s *TCPServer) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {

			switch {
			case errors.Is(err, io.EOF):
				s.removeConnection(conn.RemoteAddr().String())
				fmt.Println("Conn closed (EOF):", conn.RemoteAddr().String())
				return
			case errors.Is(err, net.ErrClosed):
				s.removeConnection(conn.RemoteAddr().String())
				fmt.Println("Conn closed (ErrClosed):", conn.RemoteAddr().String())
				return
			case errors.As(err, new(*net.OpError)):
				s.removeConnection(conn.RemoteAddr().String())
				fmt.Println("Conn closed (OpError):", conn.RemoteAddr().String())
				return
			default:
				fmt.Println(reflect.TypeOf(err))
				continue
			}
		}

		s.msgch <- *commons.NewMessage(conn.RemoteAddr().String(), buf[:n])
	}
}

func (s *TCPServer) Consume() <-chan commons.Message {
	return s.msgch
}

func (s *TCPServer) addConnecton(addr string, conn net.Conn) {
	s.connsMapMu.Lock()
	s.connsMap[addr] = conn
	s.connsMapMu.Unlock()
}

func (s *TCPServer) removeConnection(addr string) {
	s.connsMapMu.Lock()
	delete(s.connsMap, addr)
	s.connsMapMu.Unlock()
}

func (s *TCPServer) CloseConnection(addr string) {
	s.connsMap[addr].Close() //concurrency safety?
	s.removeConnection(addr)
}

func (s *TCPServer) SendMessage(addr string, payload []byte) {
	fmt.Println(s.connsMap)
	conn, ok := s.connsMap[addr]

	if ok {
		fmt.Printf("Sending message to %v: %v\n", addr, string(payload))
		conn.Write(payload)
	} else {
		fmt.Printf("ERROR (SEND): %s is not connected\n", addr)
	}
}
