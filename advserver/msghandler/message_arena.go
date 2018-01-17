package msghandler

import (
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"time"
)

func InitMessageArena() {
	message := map[uint16]ProcessFunc{
		uint16(structs.Protocol_OpenArena_Req): OpenArenaReq, // 打开竞技场
	}

	for k, v := range message {
		MapFunc[k] = v
	}
}

func OpenArenaReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("OpenArenaReq")

	sess.RefreshPlayerInfo(nil)

	sess.RefreshArenaData()

	lessTime := time.Unix(sess.PlayerData.Arena.NextRefrshTime, 0).Sub(time.Now()).Seconds()

	respSyncArena := &structs.SyncArenaNtf{
		Targets:         sess.PlayerData.Arena.Targets,
		LessRefreshTime: int32(lessTime),
		ChallengeCount:  sess.PlayerData.Arena.ChallengeCount,
	}

	respRwd := &structs.ArenaRecieveRewardResp{
		RewardRecords: sess.PlayerData.Arena.RewardRecord,
	}

	sess.Send(structs.Protocol_SyncArena_Ntf, respSyncArena)

	sess.Send(structs.Protocol_RecieveArenaReward_Resp, respRwd)
}
