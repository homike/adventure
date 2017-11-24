package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type ProcessFunc func(conn net.Conn, msgBody []byte)

var MapFunc map[uint16]ProcessFunc

func init() {
	MapFunc = map[uint16]ProcessFunc{
		uint16(1007): LoginServerPlatform,
		uint16(1002): CreatePlayer,
		uint16(1009): NameExists,
		uint16(3):    GetSystemTime,
	}
}

type SLoginServerPlatform struct {
	Takon     string
	Version   int
	ChannelID string
}

// func Float32ToByte(number int64) []byte {
// 	bits := math.int64(number)
// 	bytes := make([]byte, 4)
// 	binary.LittleEndian.PutUint32(bytes, bits)

// 	return bytes
// }

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

// 3
func GetSystemTime(conn net.Conn, msgBody []byte) {
	timeNow := time.Now().Unix()
	fmt.Println("czx@@@ GetSystemTime: ", timeNow)

	retByte := []byte{}
	retByte = append(retByte, Int64ToBytes(timeNow)...)
	ConnectSend(conn, 4, retByte)
}

// 1002
func CreatePlayer(conn net.Conn, msgBody []byte) {
	fmt.Println("czx@@@ CreatePlayer:", string(msgBody))

	retByte := []byte{}
	retByte = append(retByte, byte(0))
	ConnectSend(conn, 1003, retByte)

	SyncPlayerBaseInfo(conn)

	SyncUserGuidRecords(conn)

	SyncLoginDataFinish(conn)
}

// 1006
func SyncLoginDataFinish(conn net.Conn) {
	retByte := []byte{}
	ConnectSend(conn, 1006, retByte)
}

// 1007
func LoginServerPlatform(conn net.Conn, msgBody []byte) {
	fmt.Println("czx@@@ LoginServerPlatform:", string(msgBody))

	loginData := SLoginServerPlatform{}
	for k, v := range msgBody {
		if v == byte(0) {
			err := errors.New("body error")
			loginData.Takon = string(msgBody[0:k])
			loginData.Version, err = strconv.Atoi(string(msgBody[k+1 : k+5]))
			if err != nil {
			}
			loginData.ChannelID = string(msgBody[k+5 : len(msgBody)])

			fmt.Printf("takon: %v, version: %v, channnelid: %v", loginData.Takon, loginData.Version, loginData.ChannelID)
			break
		}
	}

	retByte := []byte{}
	retByte = append(retByte, byte(0))
	retByte = append(retByte, byte(0)) // 0:not hasCreatePlayer, 1:hasCreatePlayer
	ConnectSend(conn, 1001, retByte)

	GetSystemTime(conn, nil)

	SyncPlayerBaseInfo(conn)

	SyncLoginDataFinish(conn)

	SyncUserGuidRecords(conn)
}

// 1008
type SPlayerBaseInfo struct {
	PlayerID           int
	GameZoonID         int
	IsSupperMan        bool
	PlatformType       int
	VipLevel           int
	TotalRechargeIngot int
}

func SyncPlayerBaseInfo(conn net.Conn) {
	fmt.Println("czx@@@ SyncPlayerBaseInfo")
	playerInfo := SPlayerBaseInfo{
		PlayerID:           1,
		GameZoonID:         1,
		IsSupperMan:        true,
		PlatformType:       1,
		VipLevel:           1,
		TotalRechargeIngot: 1,
	}

	retByte := []byte{}
	retByte = append(retByte, IntToBytes(playerInfo.PlayerID)...)
	retByte = append(retByte, IntToBytes(playerInfo.GameZoonID)...)
	retByte = append(retByte, byte(1))
	retByte = append(retByte, IntToBytes(playerInfo.PlatformType)...)
	retByte = append(retByte, IntToBytes(playerInfo.VipLevel)...)
	retByte = append(retByte, IntToBytes(playerInfo.TotalRechargeIngot)...)

	ConnectSend(conn, 1008, retByte)
}

// 1009
func NameExists(conn net.Conn, msgBody []byte) {
	fmt.Println("czx@@@ NameExists:", string(msgBody))

	newName := ""
	for k, v := range msgBody {
		if v == byte(0) {
			newName = string(msgBody[0:k])
			fmt.Println("newName", newName)
			break
		}
	}

	nameBytes := []byte(newName)
	retByte := []byte{}
	retByte = append(retByte, nameBytes...)
	retByte = append(retByte, byte(0))

	ConnectSend(conn, 1010, retByte)
}

// 1413
func SyncUserGuidRecords(conn net.Conn) {
	fmt.Println("czx@@@ SyncUserGuidRecords:")

	retByte := []byte{}
	ConnectSend(conn, 1413, retByte)
}

func main() {
	// Socket Listen
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP("127.0.0.1"), 9110, ""})
	if err != nil {
		fmt.Println("监听接口失败", err.Error())
		return
	}
	fmt.Println("等待客户端连接")
	Server(listen)
}

func Server(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
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

func ConnectSend(conn net.Conn, msgID uint16, message []byte) {
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
