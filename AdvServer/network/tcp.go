package network

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type TCPClient struct {
	Listen    *net.TCPListener
	msgParser *MsgParser
}

func NewTCPListenter() *TCPClient {
	client := &TCPClient{}
	// Socket Listen
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP("127.0.0.1"), 9110, ""})
	if err != nil {
		fmt.Println("监听接口失败", err.Error())
		return nil
	}
	client.Listen = listen

	return client
}

func (t *TCPClient) Run() {
	fmt.Println("等待客户端连接")

	msgParse := NewMsgParser()
	for {
		conn, err := t.Listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常", err.Error())
			continue
		}

		go func() {
			defer conn.Close()
			bufReader := bufio.NewReader(conn)

			for {
				msgID, msgBody, err := msgParse.Read(bufReader)
				if err != nil {
					log.Println("gate message read error")
					return
				}
				fmt.Println("msgID", msgID)

				// Process
				processFunc, ok := MapFunc[msgID]
				if ok {
					go processFunc(conn, msgBody)
				}
			}
		}()
	}
}

func ConnectSend(conn net.Conn, msgID uint16, msgStruct interface{}) {
	message := MsgMarshal(msgStruct)

	writer := bufio.NewWriter(conn)
	binary.Write(writer, binary.LittleEndian, uint32(len(message)+6))
	binary.Write(writer, binary.LittleEndian, msgID)
	binary.Write(writer, binary.LittleEndian, message)
	writer.Flush()
}

func ConnectRead(bufReader *bufio.Reader) (uint16, []byte, error) {
	var headerSize uint32
	err := binary.Read(bufReader, binary.LittleEndian, &headerSize)
	if err != nil {
		log.Println("read headsize error")
		return 0, nil, err
	}

	var msgID uint16
	err = binary.Read(bufReader, binary.LittleEndian, &msgID)
	if err != nil {
		log.Println("read msgid error")
		return 0, nil, err
	}

	//fmt.Println("headerSize", headerSize, "msgID:", msgID)
	bodySize := headerSize - 6
	bodyData := make([]byte, bodySize)
	err = binary.Read(bufReader, binary.LittleEndian, &bodyData)
	if err != nil {
		log.Println("read body error")
		return 0, nil, err
	}

	return msgID, bodyData, nil
}
