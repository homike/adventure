package msghandler

import (
	"adventure/advserver/network"
	"Adventure/common/structs"
	"fmt"
)

// 1006
func SyncLoginDataFinish(client *network.TCPClient) {
	resp := &structs.SyncLoginDataFinishNtf{}
	client.Write(structs.Protocol_SyncLoginDataFinish_Ntf, resp)
}

// 1007
func LoginServerPlatform(client *network.TCPClient, msgBody []byte) {
	fmt.Println("czx@@@ LoginServerPlatform:", msgBody)

	req := structs.LoginServerPlatformReq{}
	client.UnMarshal(msgBody, &req)
	fmt.Printf("takon: %v, version: %v, channnelid: %v", req.Takon, req.Version, req.ChannelID)

	isExistsPlayer := false
	resp := &structs.LoginServerResultNtf{
		Result:         0,
		IsCreatePlayer: isExistsPlayer,
	}
	client.Write(structs.Protocol_LoginServerResult_Ntf, resp)
	GetSystemTime(client, nil)

	if isExistsPlayer {
		SyncPlayerBaseInfo(client)

		SyncLoginDataFinish(client)
	}
	SyncUserGuidRecords(client)

}
