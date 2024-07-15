package transport

import (
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"time"

	"github.com/gabrielebnc/OrderMatcher/commons/core"
)

type TCPClient struct {
	dl     net.Dialer
	quitch chan struct{}
	msgch  chan core.Message

	connsMapMu sync.RWMutex
	connsMap   map[string]net.Conn
}

func NewTCPClient() *TCPClient {
	return &TCPClient{
		dl: net.Dialer{
			Timeout: 10 * time.Second,
		},
		quitch:   make(chan struct{}),
		msgch:    make(chan core.Message),
		connsMap: make(map[string]net.Conn),
	}
}

func (c *TCPClient) Stop() {
	close(c.msgch)

}

func (c *TCPClient) Dial(addr string) {
	conn, err := c.dl.Dial("tcp", addr)
	if err != nil {
		fmt.Println("ERROR (DIAL):", err)
		fmt.Println("ERR TYPE:", reflect.TypeOf(err))
		return
	}
	c.addConnecton(addr, conn)

	go c.readLoop(conn)
}

func (c *TCPClient) SendMessage(addr string, payload []byte) {
	conn, ok := c.connsMap[addr]
	if !ok {
		fmt.Printf("ERROR (SendMessage): %v is not connected\n", addr)
		return
	}
	conn.Write(payload)
}

func (c *TCPClient) readLoop(conn net.Conn) {
	fmt.Println("readloop of conn: ", conn.RemoteAddr().String())
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			switch {
			case errors.Is(err, io.EOF):
				c.removeConnection(conn.RemoteAddr().String())
				fmt.Println("Conn closed (EOF):", conn.RemoteAddr().String())
				return
			case errors.Is(err, net.ErrClosed):
				c.removeConnection(conn.RemoteAddr().String())
				fmt.Println("Conn closed (ErrClosed):", conn.RemoteAddr().String())
				return
			case errors.As(err, new(*net.OpError)):
				c.removeConnection(conn.RemoteAddr().String())
				fmt.Println("Conn closed (OpError):", conn.RemoteAddr().String())
				return
			default:
				fmt.Println(reflect.TypeOf(err))
				continue
			}
		}
		c.msgch <- *core.NewMessage(conn.RemoteAddr().String(), buf[:n])
		fmt.Printf("msg from %v: %v\n", conn.RemoteAddr().String(), buf[:n])
	}
}

func (c *TCPClient) Consume() <-chan core.Message {
	return c.msgch
}

func (c *TCPClient) CloseConnection(addr string) {
	if conn, ok := c.connsMap[addr]; ok {
		conn.Close()
		c.removeConnection(addr)
	}
}

func (c *TCPClient) addConnecton(addr string, conn net.Conn) {
	c.connsMapMu.Lock()
	c.connsMap[addr] = conn
	c.connsMapMu.Unlock()
}

func (c *TCPClient) removeConnection(addr string) {
	c.connsMapMu.Lock()
	delete(c.connsMap, addr)
	c.connsMapMu.Unlock()
}
