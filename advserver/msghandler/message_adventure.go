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

	sess.RefreshPlayerInfo(nil)

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

	gameLevelT, ok := gamedata.AllTemplates.GameLevelTemplates[gameLevel.GameLevelID]
	if !ok {
		logger.Error("GameLevelTemplates(%v) failed", gameLevel.GameLevelID)
		sess.Send(structs.Protocol_OpenGameBox_Resp, resp)
		return
	}
	random := util.NewRandom(gameLevelT.GameBoxIDs, gameLevelT.GameBoxWeight)
	if random == nil {
		sess.Send(structs.Protocol_OpenGameBox_Resp, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	rewardIDs := []int32{}
	for i := 0; i < int(req.Count); i++ {
		rewardID := random.GetRandomNum()
		rewardIDs = append(rewardIDs, int32(rewardID))
	}

	clientRewards := []*structs.Reward{}
	for _, v := range rewardIDs {
		rewardT, ok := gamedata.AllTemplates.RewardTemplates[v]
		if ok {
			reward, _ := sess.DoReward(&rewardT, 1)
			for _, r := range clientRewards {
				if r.RewardType == reward.RewardType && r.Param1 == reward.Param1 {
					r.Param2 += reward.Param2
				} else {
					clientRewards = append(clientRewards, reward)
				}
			}
		}
	}

	sess.SyncCurrentGameLevelNtf()

	resp.Ret = structs.AdventureRet_Success
	resp.Count = req.Count
	resp.Rewards = clientRewards
	sess.Send(structs.Protocol_OpenGameBox_Resp, resp)

	//CZXDO: 开宝箱成就检测
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

func AdventureEventReq(sess *sessions.Session, msgBody []byte) {

	logger.Debug("AdventureEventReq")

	req := &structs.AdventureEventReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.AdventureEventResp{
		Ret:         structs.AdventureRet_Failed,
		GameLevelID: sess.PlayerData.PlayerGameLevel.CurrentGameLevelID,
		EventID:     req.EventID,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	sess.RefreshPlayerInfo(nil)

	gameLevelT, ok := gamedata.AllTemplates.GameLevelTemplates[sess.PlayerData.PlayerGameLevel.CurrentGameLevelID]
	if !ok {
		logger.Error("AdventureEventReq AdventureEventReq(%v)", sess.PlayerData.PlayerGameLevel.CurrentGameLevelID)
		sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
		return
	}

	if req.EventID < 0 || int(req.EventID) >= len(gameLevelT.EvnetIDs) {
		logger.Error("AdventureEventReq req.EventID(%v) is out of range", req.EventID)
		sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
		return
	}

	eventTemplateID := gameLevelT.EvnetIDs[req.EventID]
	eventT, ok := gamedata.AllTemplates.GameLevelEventTemplates[eventTemplateID]
	if !ok {
		logger.Error("AdventureEventReq GameLevelEventTemplates(%v)", eventTemplateID)
		sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
		return
	}

	if sess.PlayerData.Res.Strength < eventT.CostFood {
		logger.Error("AdventureEventReq (%v) food not enough", req.EventID)
		sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
		return
	}

	gameLevel, err := sess.PlayerData.PlayerGameLevel.GetCurGameLevelData()
	if err != nil {
		logger.Error("AdventureEventReq (%v) GetCurGameLevelData() faield %v", req.EventID, err)
		sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
		return
	}
	if int(req.EventID) > len(gameLevel.CompleteEvent) {
		logger.Error("AdventureEventReq (%v) cannot find event", req.EventID)
		sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
		return
	}
	eventStatus := gameLevel.CompleteEvent[req.EventID]
	if eventStatus == structs.AdventureEventStatus_UnActive || eventStatus == structs.AdventureEventStatus_Finish {
		logger.Error("AdventureEventReq (%v) event not active or finish", req.EventID)
		sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
		return
	}

	switch eventT.Type {
	case structs.GameLevelType_NpcTack:
	case structs.GameLevelType_Item:
		if !sess.PlayerData.Bag.HasEnoughItem(eventT.ItemId, eventT.Num) {
			logger.Error("AdventureEventReq (%v) item(%v, %v) not enough ", eventT.ItemId, eventT.Num)
			sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
			return
		}
	case structs.GameLevelType_Res:
		if !sess.PlayerData.Res.HasEnoughOres(eventT.ResId, eventT.Num) {
			logger.Error("AdventureEventReq (%v) res(%v, %v) not enough ", eventT.ItemId, eventT.Num)
			sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
			return
		}
	case structs.GameLevelType_Fight:
		leftTime := sess.PlayerData.ExtendData.GetGamelevelLeftTime(eventT.Num)
		if leftTime > 0 {
			logger.Error("AdventureEventReq (%v) gamelevel is CD(%v) ", leftTime)
			sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
			return
		}
	default:
		logger.Error("AdventureEventReq (%v) event type(%v) is error", eventT.Type)
		sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////

	err = sess.PlayerData.Res.StrengthChange(-eventT.CostFood)
	if err == nil {
		sess.SyncStrength()
	}

	switch eventT.Type {
	case structs.GameLevelType_NpcTack:
	case structs.GameLevelType_Item:
		sess.PlayerData.Bag.MinusItem(eventT.ItemId, eventT.Num)

	case structs.GameLevelType_Res:
		sess.PlayerData.Res.OresChange(eventT.ResId, eventT.Num)

	case structs.GameLevelType_Fight:
		ret, err := sess.DoFightTest(eventT.Num)
		if err != nil {
			logger.Error("DoFightTest (%v) dofight failed %v", eventT.Num, err)
			sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
			return
		}
		fightNtf := &structs.FightResultNtf{
			FType:  structs.FightType_EventFight,
			Result: ret,
		}
		sess.Send(structs.Protocol_FightResult_Ntf, fightNtf)
		if !ret.IsLeftWin {
			logger.Error("AdventureEventReq () fight lose")
			sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
			return
		}
	}
	gameLevel.CompleteEvent[req.EventID] = structs.AdventureEventStatus_Finish

	// 更新关卡数据
	sess.SyncCurrentGameLevelNtf()
	// 通知客户端，时间完成
	resp.Ret = structs.AdventureRet_Success
	sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
	// 发放奖励
	sess.DoSomeRewards(eventT.RewardIDs)

	//CZXDO: 通过公告

	//CZXDO: 成就检测

	sess.Send(structs.Protocol_AdventureEvent_Resp, resp)
}
