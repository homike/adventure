package msghandler

import (
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"time"
)

// 3
func GetSystemTime(sess *sessions.Session, msgBody []byte) {

	timeNow := time.Now().Unix()
	resp := &structs.GetSystemTimeResp{
		Time: timeNow,
	}
	//fmt.Println("czx@@@ GetSystemTime: ", timeNow)

	sess.Send(structs.Protocol_GetSystemTime_Resp, resp)
}

// 2801
func SetPlayerBarrageConfig(sess *sessions.Session, msgBody []byte) {
	//fmt.Println("czx@@@ SetPlayerBarrageConfig")
	//client.Write(structs.Protocol_SyncUserGuidRecords_Ntf, resp)
}
