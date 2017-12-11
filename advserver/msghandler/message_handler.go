package msghandler

import (
	"adventure/advserver/log"
	"adventure/advserver/model"
	"adventure/advserver/network"
	"adventure/advserver/sessions"
	"adventure/common/clog"
	"adventure/common/structs"
	"fmt"
)

type ProcessFunc func(sess *sessions.Session, msgBody []byte)

var (
	MapFunc map[uint16]ProcessFunc
	logger  *clog.Logger
)

func Dispatch(msgID uint16, msgBody []byte, tc *network.TCPClient) {

	sess, err := sessions.SessionMgr.FindSession(tc.AccountID)
	if err != nil {
		sess = sessions.NewSession(tc)
	}

	processFunc, ok := MapFunc[msgID]
	if ok {
		processFunc(sess, msgBody)
	}
}

func init() {
	logger = log.GetLogger()

	MapFunc = map[uint16]ProcessFunc{
		uint16(structs.Protocol_Test_Req):                   TestReq,
		uint16(structs.Protocol_GetSystemTime_Req):          GetSystemTime,
		uint16(structs.Protocol_CreatePlayer_Req):           CreatePlayer,
		uint16(structs.Protocol_LoginServerPlatform_Req):    LoginServerPlatform,
		uint16(structs.Protocol_NameExists_Req):             NameExists,
		uint16(structs.Protocol_SetPlayerBarrageConfig_Req): SetPlayerBarrageConfig,
		uint16(structs.Protocol_UpdateUserGuidRecord_Req):   UpdateUserGuidRecord,
	}
}

// 1
func TestReq(sess *sessions.Session, msgBody []byte) {

	player, err := model.NewPlayer("czx", 1)
	if err != nil {
		fmt.Println("NewPlayer Error", err)
		return
	}
	sess.SetPlayer(player)
	sessions.SessionMgr.AddSession(sess)

	resp := &structs.SyncLoginDataFinishNtf{}
	sess.Send(structs.Protocol_Test_Resp, resp) //structs.Protocol_Test_Resp, resp)
}
