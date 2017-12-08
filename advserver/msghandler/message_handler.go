package msghandler

import (
	"Adventure/AdvServer/model"
	"Adventure/AdvServer/network"
	"Adventure/AdvServer/sessions"
	"Adventure/common/structs"
	"fmt"
)

type ProcessFunc func(client *network.TCPClient, msgBody []byte)

var MapFunc map[uint16]ProcessFunc

func Dispatch(msgID uint16, msgBody []byte, tc *network.TCPClient) {
	processFunc, ok := MapFunc[msgID]
	if ok {
		processFunc(tc, msgBody)
	}
}

func init() {
	MapFunc = map[uint16]ProcessFunc{
		uint16(structs.Protocol_Test_Req):                TestReq,
		uint16(structs.Protocol_GetSystemTime_Req):       GetSystemTime,
		uint16(structs.Protocol_CreatePlayer_Req):        CreatePlayer,
		uint16(structs.Protocol_LoginServerPlatform_Req): LoginServerPlatform,
		uint16(structs.Protocol_NameExists_Req):          NameExists,
	}
}

// 1
func TestReq(client *network.TCPClient, msgBody []byte) {

	player, err := model.NewPlayer("czx", 1)
	if err != nil {
		fmt.Println("NewPlayer Error", err)
		return
	}
	sessions.SessionMgr.CreateSession(player, client)

	resp := &structs.SyncLoginDataFinishNtf{}
	client.Write(structs.Protocol_Test_Resp, resp) //structs.Protocol_Test_Resp, resp)
}
