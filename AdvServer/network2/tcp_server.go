package network2

import (
	"fmt"
	"net"
)

type TCPServer struct {
	Addr       string
	Port       int
	MaxConnNum int

	Listen *net.TCPListener
}

// NewTCPServer :
func NewTCPServer() *TCPServer {
	server := &TCPServer{
		Addr: "127.0.0.1",
		Port: 9110,
	}
	// Socket Listen
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(server.Addr), server.Port, ""})
	if err != nil {
		fmt.Println("Listen failed", err.Error())
		return nil
	}
	server.Listen = listen

	return server
}

// Run :
func (t *TCPServer) Run() {
	fmt.Println("Server Listen at", t.Addr, t.Port)

	for {
		conn, err := t.Listen.AcceptTCP()
		if err != nil {
			fmt.Println("client connect error", err.Error())
			continue
		}

		agent := NewTCPClient(conn)
		go func() {
			agent.Run()

			// cleanup
			conn.Close()
		}()
	}
}
