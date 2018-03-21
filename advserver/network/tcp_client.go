package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

type MsgHandler func(msgID uint16, msgBody []byte, tc *TCPClient)

type TCPClient struct {
	sync.Mutex
	AccountID uint
	Parser    *MsgParser
	conn      net.Conn
	WriteChan chan []byte
	//ReadChan  chan []byte
}

func NewTCPClient(conn net.Conn, parser *MsgParser) *TCPClient {
	client := &TCPClient{
		AccountID: 0,
		conn:      conn,
		WriteChan: make(chan []byte, 128),
		//ReadChan:  make(chan []byte, 128),
		Parser: parser,
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

func (tc *TCPClient) Write(msgID uint16, msgStruct interface{}) {
	tc.Lock()
	defer tc.Unlock()

	data := tc.Parser.Pack(msgID, msgStruct)
	tc.WriteChan <- data
}

func (tc *TCPClient) Run(handler MsgHandler) {
	bufReader := bufio.NewReader(tc.conn)
	for {
		msgID, msgBody, err := tc.Parser.UnPack(bufReader)
		if err != nil {
			log.Println("gate message read error")
			return
		}
		fmt.Println("msgID", msgID)

		//CZXDO: 是否有必要gorutine
		go handler(msgID, msgBody, tc)
	}
}

func (tc *TCPClient) UnMarshal(msgBody []byte, msgStruct interface{}) {
	tc.Parser.MsgProcessor.UnMarshal(msgBody, msgStruct)
}
