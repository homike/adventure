package msghandler

import (
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"fmt"
)

func InitMessageLogin() {
	handlers := map[uint16]ProcessFunc{
		uint16(structs.Protocol_LoginServerPlatform_Req): LoginServerPlatform,
	}

	for k, v := range handlers {
		MapFunc[k] = v
	}
}

// 1007
func LoginServerPlatform(sess *sessions.Session, msgBody []byte) {
	fmt.Println("czx@@@ LoginServerPlatform:", msgBody)

	req := structs.LoginServerPlatformReq{}
	sess.UnMarshal(msgBody, &req)
	fmt.Printf("takon: %v, version: %v, channnelid: %v", req.Takon, req.Version, req.ChannelID)

	isExistsPlayer := false
	resp := &structs.LoginServerResultNtf{
		Result:         0,
		IsCreatePlayer: isExistsPlayer,
	}
	sess.Send(structs.Protocol_LoginServerResult_Ntf, resp)
	GetSystemTime(sess, nil)

	if isExistsPlayer {
		sess.SyncPlayerBaseInfo()

		SyncLoginDataFinish(sess)
	}
	sess.SyncUserGuidRecords()
}

// 1006
func SyncLoginDataFinish(sess *sessions.Session) {
	resp := &structs.SyncLoginDataFinishNtf{}
	sess.Send(structs.Protocol_SyncLoginDataFinish_Ntf, resp)
}
