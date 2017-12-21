package msghandler

import (
	"adventure/advserver/gamedata"
	"adventure/advserver/sessions"

	"adventure/common/structs"
	"adventure/common/util"
)

func InitMessageAdventure() {
	message := map[uint16]ProcessFunc{
		uint16(structs.Protocol_SelectGameLevel_Req):     SelectGameLevelReq,     // 切换关卡
		uint16(structs.Protocol_AdventureEvent_Req):      AdventureEventReq,      // 冒险事件
		uint16(structs.Protocol_OpenGameBox_Req):         OpenGameBoxReq,         // 打开宝箱
		uint16(structs.Protocol_GetFightCoolingTime_Req): GetFightCoolingTimeReq, // 获取战败冷却时间
		uint16(structs.Protocol_SpeedAdventure_Req):      TestReq,                // 加速冒险
	}

	for k, v := range message {
		MapFunc[k] = v
	}
}

func SelectGameLevelReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("SelectGameLevelReq")

	req := &structs.SelectGameLevelReq{}
	sess.UnMarshal(msgBody, req)

	gameLevel, err := sess.PlayerData.PlayerGameLevel.GetGameLevelData(req.LevelID)
	if err != nil {
		logger.Error("GetGameLevelData(%v) error %v", req.LevelID, err)
		return
	}
	if !gameLevel.IsUnlock {
		logger.Error("%v is lock", req.LevelID)
		return
	}

	sess.PlayerData.PlayerGameLevel.SelectGameLevel(req.LevelID)

	//CZXDO: 同步关卡数据2
	sess.SyncCurrentGameLevelNtf()
}

func AdventureEventReq(sess *sessions.Session, msgBody []byte) {

}

func OpenGameBoxReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("GetFightCoolingTimeReq")

	req := &structs.OpenGameBoxReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.OpenGameBoxResp{
		Ret: structs.AdventureRet_Failed,
	}
	/////////////////////////////////////////////Data Check////////////////////////////////////////

	gameLevel, err := sess.PlayerData.PlayerGameLevel.GetCurGameLevelData()
	if err != nil {
		logger.Error("GetCurGameLevelData() Error(%v)", err)
		sess.Send(structs.Protocol_OpenGameBox_Resp, resp)
		return
	}

	if gameLevel.BoxCount <= 0 {
		logger.Error("BoxCount(%v) not enough", gameLevel.BoxCount)
		sess.Send(structs.Protocol_OpenGameBox_Resp, resp)
		return
	}

	boxIDs, boxWeights, err := gamedata.GetGameLevelGameBox(gameLevel.GameLevelID)
	if err != nil {
		logger.Error("GetGameLevelGameBox(%v) error(%v)", gameLevel.GameLevelID, boxIDs)
		sess.Send(structs.Protocol_OpenGameBox_Resp, resp)
		return
	}
	random := util.NewRandom(boxIDs, boxWeights)
	if random == nil {
		logger.Error("GetGameLevelGameBox(%v) error(%v)", gameLevel.GameLevelID, boxIDs)
		sess.Send(structs.Protocol_OpenGameBox_Resp, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	rewardIDs := []int{}
	for i := 0; i < int(req.Count); i++ {
		rewardID := random.GetRandomNum()
		rewardIDs = append(rewardIDs, rewardID)
	}

}

func GetFightCoolingTimeReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("GetFightCoolingTimeReq")

	req := &structs.GetFightCoolingTimeReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.GetFightCoolingTimeResp{
		LeftTime: sess.PlayerData.ExtendData.GetGamelevelLeftTime(req.BattleID),
	}

	sess.Send(structs.Protocol_GetFightCoolingTime_Resp, resp)
}
