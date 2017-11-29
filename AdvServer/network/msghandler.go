package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"
)

type ProcessFunc func(conn net.Conn, msgBody []byte)

var MapFunc map[uint16]ProcessFunc

func init() {
	MapFunc = map[uint16]ProcessFunc{
		uint16(Protocol_GetSystemTime_Req):       GetSystemTime,
		uint16(Protocol_CreatePlayer_Req):        CreatePlayer,
		uint16(Protocol_LoginServerPlatform_Req): LoginServerPlatform,
		uint16(Protocol_NameExists_Req):          NameExists,
	}
}

func MsgUnMarshal(msgBody []byte, msgStruct interface{}) {
	readIndex := 0
	v := reflect.ValueOf(msgStruct).Elem()
	vType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)
		tf := vType.Field(i)

		fmt.Println(tf.Name, vf.Kind())
		switch vf.Kind() {
		case reflect.String:
			for i := readIndex; i < len(msgBody); i++ {
				if msgBody[i] == byte(0) {
					fmt.Println(readIndex, "string :", msgBody[readIndex:i])
					canSetValue := reflect.ValueOf(string(msgBody[readIndex:i]))
					vf.Set(canSetValue)
					readIndex = i + 1
					break
				}
			}

		case reflect.Int32:
			fmt.Println(readIndex, "int :", msgBody[readIndex:readIndex+4])
			intValue, err := strconv.Atoi(string(msgBody[readIndex : readIndex+4]))
			if err != nil {
			}
			canSetValue := reflect.ValueOf(int32(intValue))
			vf.Set(canSetValue)
			readIndex = readIndex + 4
		default:
		}
	}
}

func MsgMarshal(msgStruct interface{}) []byte {

	bytesBuffer := bytes.NewBuffer([]byte{})

	v := reflect.ValueOf(msgStruct).Elem()
	for i := 0; i < v.NumField(); i++ {
		vf := v.Field(i)

		switch vf.Kind() {
		case reflect.String:
			binary.Write(bytesBuffer, binary.LittleEndian, vf.String())
			binary.Write(bytesBuffer, binary.LittleEndian, byte(0))

		case reflect.Int32:
			binary.Write(bytesBuffer, binary.LittleEndian, int32(vf.Int()))

		case reflect.Int64:
			binary.Write(bytesBuffer, binary.LittleEndian, vf.Int())

		case reflect.Bool:
			b := 0
			if vf.Bool() {
				b = 1
			}
			binary.Write(bytesBuffer, binary.LittleEndian, uint8(b))

		case reflect.Slice:

		default:
			binary.Write(bytesBuffer, binary.LittleEndian, vf.Bytes())
		}
	}

	return bytesBuffer.Bytes()
}

// 3
func GetSystemTime(conn net.Conn, msgBody []byte) {

	timeNow := time.Now().Unix()
	resp := &GetSystemTimeResp{
		Time: timeNow,
	}
	fmt.Println("czx@@@ GetSystemTime: ", timeNow)

	ConnectSend(conn, Protocol_GetSystemTime_Resp, resp)
}

// 1002
func CreatePlayer(conn net.Conn, msgBody []byte) {
	fmt.Println("czx@@@ CreatePlayer:", string(msgBody))

	resp := &CreatePlayerResp{
		Result: 0, // Success
	}
	ConnectSend(conn, Protocol_CreatePlayer_Resp, resp)

	SyncPlayerBaseInfo(conn)

	SyncUserGuidRecords(conn)

	SyncLoginDataFinish(conn)
}

// 1006
func SyncLoginDataFinish(conn net.Conn) {
	resp := &SyncLoginDataFinishNtf{}
	ConnectSend(conn, Protocol_SyncLoginDataFinish_Ntf, resp)
}

// 1007
func LoginServerPlatform(conn net.Conn, msgBody []byte) {
	fmt.Println("czx@@@ LoginServerPlatform:", msgBody)

	loginReq := LoginServerPlatformReq{}
	MsgUnMarshal(msgBody, &loginReq)
	fmt.Printf("takon: %v, version: %v, channnelid: %v", loginReq.Takon, loginReq.Version, loginReq.ChannelID)

	isExistsPlayer := false
	resp := &LoginServerResultNtf{
		Result:         0,
		IsCreatePlayer: isExistsPlayer,
	}
	ConnectSend(conn, Protocol_LoginServerResult_Ntf, resp)
	GetSystemTime(conn, nil)

	if isExistsPlayer {
		SyncPlayerBaseInfo(conn)

		SyncLoginDataFinish(conn)

		SyncUserGuidRecords(conn)
	}
}

// 1008
func SyncPlayerBaseInfo(conn net.Conn) {
	fmt.Println("czx@@@ SyncPlayerBaseInfo")

	resp := &SyncPlayerBaseInfoNtf{
		PlayerID:           1,
		GameZoonID:         1,
		IsSupperMan:        true,
		PlatformType:       1,
		Viplevel:           1,
		TotalRechargeIngot: 1,
	}
	ConnectSend(conn, Protocol_SyncPlayerBaseInfo_Ntf, resp)
}

// 1009
func NameExists(conn net.Conn, msgBody []byte) {
	fmt.Println("czx@@@ NameExists:", string(msgBody))

	req := NameExistsReq{}
	MsgUnMarshal(msgBody, &req)

	resp := &NameExistsResp{
		Name: req.Name,
	}
	ConnectSend(conn, Protocol_NameExists_Resp, resp)
}

// 1413
func SyncUserGuidRecords(conn net.Conn) {
	fmt.Println("czx@@@ SyncUserGuidRecords:")
	resp := &SyncUserGuidRecordsNtf{}
	ConnectSend(conn, Protocol_SyncUserGuidRecords_Ntf, resp)
}
