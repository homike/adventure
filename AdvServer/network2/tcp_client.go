package network2

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type TCPClient struct {
	sync.Mutex
	conn      net.Conn
	WriteChan chan []byte
	ReadChan  chan []byte
}

func NewTCPClient(conn net.Conn) *TCPClient {
	client := &TCPClient{
		conn:      conn,
		WriteChan: make(chan []byte, 128),
		ReadChan:  make(chan []byte, 128),
	}

	go func() {
		for b := range client.WriteChan {
			if b == nil {
				break
			}

			_, err := conn.Write(b)
			if err != nil {
				break
			}
		}

		conn.Close()
	}()

	return client
}

func (tc *TCPClient) Write(b []byte) {
	tc.Lock()
	defer tc.Unlock()

	tc.WriteChan <- b
}

func (tc *TCPClient) Run() {
	bufReader := bufio.NewReader(tc.conn)
	for {
		msgID, msgBody, err := MsgParserSingleton.Read(bufReader)
		if err != nil {
			log.Println("gate message read error")
			return
		}
		fmt.Println("msgID", msgID)

		processFunc, ok := MapFunc[msgID]
		if ok {
			go processFunc(tc, msgBody)
		}
	}
}
