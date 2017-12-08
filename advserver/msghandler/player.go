package msghandler

import (
	"Adventure/AdvServer/model"
	"Adventure/AdvServer/network"
	"Adventure/AdvServer/sessions"
	"Adventure/common/structs"
	"fmt"
)

// 1002
func CreatePlayer(client *network.TCPClient, msgBody []byte) {
	fmt.Println("czx@@@ CreatePlayer:", string(msgBody))

	req := &structs.CreatePlayerReq{}
	client.UnMarshal(msgBody, &req)

	resp := &structs.CreatePlayerResp{
		Result: 0, // Success
	}

	player, err := model.NewPlayer(req.PlayerName, req.HeroTemplateId)
	if err != nil {
		fmt.Println("NewPlayer Error", err)
		resp.Result = 1
		client.Write(structs.Protocol_CreatePlayer_Resp, resp)
		return
	}

	sessions.SessionMgr.CreateSession(player, client)

	client.Write(structs.Protocol_CreatePlayer_Resp, resp)

	SyncPlayerBaseInfo(client)

	SyncUserGuidRecords(client)

	SyncLoginDataFinish(client)
}

// 1008
func SyncPlayerBaseInfo(client *network.TCPClient) {
	fmt.Println("czx@@@ SyncPlayerBaseInfo")

	resp := &structs.SyncPlayerBaseInfoNtf{
		PlayerID:           1,
		GameZoonID:         1,
		IsSupperMan:        true,
		PlatformType:       1,
		Viplevel:           1,
		TotalRechargeIngot: 1,
	}
	client.Write(structs.Protocol_SyncPlayerBaseInfo_Ntf, resp)
}

// 1009
func NameExists(client *network.TCPClient, msgBody []byte) {
	fmt.Println("czx@@@ NameExists:", string(msgBody))

	req := structs.NameExistsReq{}
	client.UnMarshal(msgBody, &req)

	resp := &structs.NameExistsResp{
		Name: req.Name,
	}
	client.Write(structs.Protocol_NameExists_Resp, resp)
}
