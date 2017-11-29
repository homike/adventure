package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func IntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	return bytesBuffer.Bytes()
}

func Int64ToBytes(n int64) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, n)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(bytesBuffer, binary.BigEndian, &tmp)
	return int(tmp)
}

func Int2Byte(data int) []byte {
	b := []byte{0x00, 0x00, 0x03, 0xe8}
	b_buf := bytes.NewBuffer(b)
	var x int32
	binary.Read(b_buf, binary.BigEndian, &x)
	return b_buf.Bytes()
}

type TCPListener struct {
	Listen *net.TCPListener
}

func NewTCPListenter() *TCPListener {
	tcpListener := &TCPListener{}
	// Socket Listen
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP("127.0.0.1"), 9110, ""})
	if err != nil {
		fmt.Println("监听接口失败", err.Error())
		return nil
	}
	tcpListener.Listen = listen

	return tcpListener
}

func (t *TCPListener) StartAccept() {
	fmt.Println("等待客户端连接")

	for {
		conn, err := t.Listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常", err.Error())
			continue
		}
		defer conn.Close()

		func() {
			defer conn.Close()
			bufReader := bufio.NewReader(conn)

			for {
				msgID, msgBody, err := ConnectRead(bufReader)
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

	fmt.Println("headerSize", headerSize, "msgID:", msgID)

	bodySize := headerSize - 6
	bodyData := make([]byte, bodySize)
	err = binary.Read(bufReader, binary.LittleEndian, &bodyData)
	if err != nil {
		log.Println("read body error")
		return 0, nil, err
	}

	return msgID, bodyData, nil
}
