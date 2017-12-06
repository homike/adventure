package network2

import (
	"Adventure/AdvServer/model"
	"fmt"
	"time"
)

type ProcessFunc func(client *TCPClient, msgBody []byte)

var MapFunc map[uint16]ProcessFunc

func init() {
	MapFunc = map[uint16]ProcessFunc{
		uint16(Protocol_Test_Req):                TestReq,
		uint16(Protocol_GetSystemTime_Req):       GetSystemTime,
		uint16(Protocol_CreatePlayer_Req):        CreatePlayer,
		uint16(Protocol_LoginServerPlatform_Req): LoginServerPlatform,
		uint16(Protocol_NameExists_Req):          NameExists,
	}
}

// 1
func TestReq(client *TCPClient, msgBody []byte) {

	player, err := model.NewPlayer("czx", 1)
	if err != nil {
		fmt.Println("NewPlayer Error", err)
		return
	}
	SessionMgr.CreateSession(player, client)

	resp := &SyncLoginDataFinishNtf{}
	MsgParserSingleton.Write(client, Protocol_Test_Resp, resp)
}

// 3
func GetSystemTime(client *TCPClient, msgBody []byte) {

	timeNow := time.Now().Unix()
	resp := &GetSystemTimeResp{
		Time: timeNow,
	}
	fmt.Println("czx@@@ GetSystemTime: ", timeNow)

	MsgParserSingleton.Write(client, Protocol_GetSystemTime_Resp, resp)
}

// 1002
func CreatePlayer(client *TCPClient, msgBody []byte) {
	fmt.Println("czx@@@ CreatePlayer:", string(msgBody))

	req := &CreatePlayerReq{}
	MsgParserSingleton.MsgProcessor.UnMarshal(msgBody, &req)

	resp := &CreatePlayerResp{
		Result: 0, // Success
	}

	player, err := model.NewPlayer(req.PlayerName, req.HeroTemplateId)
	if err != nil {
		fmt.Println("NewPlayer Error", err)
		resp.Result = 1
		MsgParserSingleton.Write(client, Protocol_CreatePlayer_Resp, resp)
		return
	}

	SessionMgr.CreateSession(player, client)

	MsgParserSingleton.Write(client, Protocol_CreatePlayer_Resp, resp)

	SyncPlayerBaseInfo(client)

	SyncUserGuidRecords(client)

	SyncLoginDataFinish(client)
}

// 1006
func SyncLoginDataFinish(client *TCPClient) {
	resp := &SyncLoginDataFinishNtf{}
	MsgParserSingleton.Write(client, Protocol_SyncLoginDataFinish_Ntf, resp)
}

// 1007
func LoginServerPlatform(client *TCPClient, msgBody []byte) {
	fmt.Println("czx@@@ LoginServerPlatform:", msgBody)

	req := LoginServerPlatformReq{}
	MsgParserSingleton.MsgProcessor.UnMarshal(msgBody, &req)
	fmt.Printf("takon: %v, version: %v, channnelid: %v", req.Takon, req.Version, req.ChannelID)

	isExistsPlayer := false
	resp := &LoginServerResultNtf{
		Result:         0,
		IsCreatePlayer: isExistsPlayer,
	}
	MsgParserSingleton.Write(client, Protocol_LoginServerResult_Ntf, resp)
	GetSystemTime(client, nil)

	if isExistsPlayer {
		SyncPlayerBaseInfo(client)

		SyncLoginDataFinish(client)
	}
	SyncUserGuidRecords(client)

}

// 1008
func SyncPlayerBaseInfo(client *TCPClient) {
	fmt.Println("czx@@@ SyncPlayerBaseInfo")

	resp := &SyncPlayerBaseInfoNtf{
		PlayerID:           1,
		GameZoonID:         1,
		IsSupperMan:        true,
		PlatformType:       1,
		Viplevel:           1,
		TotalRechargeIngot: 1,
	}
	MsgParserSingleton.Write(client, Protocol_SyncPlayerBaseInfo_Ntf, resp)
}

// 1009
func NameExists(client *TCPClient, msgBody []byte) {
	fmt.Println("czx@@@ NameExists:", string(msgBody))

	req := NameExistsReq{}
	MsgParserSingleton.MsgProcessor.UnMarshal(msgBody, &req)

	resp := &NameExistsResp{
		Name: req.Name,
	}
	MsgParserSingleton.Write(client, Protocol_NameExists_Resp, resp)
}

// 1413
func SyncUserGuidRecords(client *TCPClient) {
	fmt.Println("czx@@@ SyncUserGuidRecords:")

	records := []GuildRecord{}
	for i := 0; i < 2; i++ {
		records = append(records, GuildRecord{
			UserGuidTypes: uint8(i + 2),
			TriggerCount:  int32(i + 3),
		})
	}
	resp := &SyncUserGuidRecordsNtf{
		Records: records,
	}

	MsgParserSingleton.Write(client, Protocol_SyncUserGuidRecords_Ntf, resp)
}
