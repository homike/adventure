package network

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
func (t *TCPServer) Run(handler MsgHandler) {
	fmt.Println("Server Listen at", t.Addr, t.Port)

	msgParser := NewMsgParser()
	for {
		conn, err := t.Listen.AcceptTCP()
		if err != nil {
			fmt.Println("client connect error", err.Error())
			continue
		}
		fmt.Println("new client : ", conn.RemoteAddr().String())
		agent := NewTCPClient(conn, msgParser)
		go func() {
			agent.Run(handler)

			// cleanup
			conn.Close()
		}()
	}
}
