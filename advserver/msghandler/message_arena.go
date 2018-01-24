package msghandler

import (
	"adventure/advserver/gamedata"
	"adventure/advserver/model"
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"adventure/common/util"
	"math"
	"time"
)

func InitMessageArena() {
	message := map[uint16]ProcessFunc{
		uint16(structs.Protocol_OpenArena_Req):          OpenArenaReq,          // 打开竞技场
		uint16(structs.Protocol_ArenaChallenge_Req):     ChallengeReq,          // 挑战
		uint16(structs.Protocol_RefreshArena_Req):       RefreshArenaReq,       // 刷新
		uint16(structs.Protocol_RecieveArenaReward_Req): RecieveArenaRewardReq, // 领取奖励
	}

	for k, v := range message {
		MapFunc[k] = v
	}
}

func OpenArenaReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("OpenArenaReq")

	sess.RefreshPlayerInfo(nil)

	sess.RefreshArenaData()

	sess.SyncPlayerArena()
}

func CheckCourageAchievement(winNum int32) int {
	index := -1
	for i := 0; i < len(gamedata.AllTemplates.GlobalData.TotalWinCount); i++ {
		if gamedata.AllTemplates.GlobalData.TotalWinCount[i] == winNum {
			index = i
			break
		}
	}
	return index
}

func CheckChampionAchievement(winNum int32) int {
	index := -1
	for i := 0; i < len(gamedata.AllTemplates.GlobalData.TotalWinCount); i++ {
		if gamedata.AllTemplates.GlobalData.TotalWinCount[i] == winNum {
			index = i
			break
		}
	}
	return index
}

func ChallengeReq(sess *sessions.Session, msgBody []byte) {
	req := &structs.ArenaChallengeReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.ArenaChallengeResp{
		Ret:      structs.AdventureRet_Failed,
		PlayerID: req.PlayerID,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	if time.Now().Unix() > sess.PlayerData.Arena.NextRefrshTime {
		logger.Error("%v ", "need refresh info ")
		sess.Send(structs.Protocol_ArenaChallenge_Resp, resp)
		return
	}

	if sess.PlayerData.Arena.ChallengeCount >= gamedata.AllTemplates.GlobalData.MaxChallengeCount {
		logger.Error("challenge count is max")
		sess.Send(structs.Protocol_ArenaChallenge_Resp, resp)
		return
	}

	var targetPlayer *structs.FightTarget
	for _, v := range sess.PlayerData.Arena.Targets {
		if v.PlayerID == req.PlayerID {
			targetPlayer = v
		}
	}
	if targetPlayer == nil {
		logger.Error("cannot find player")
		sess.Send(structs.Protocol_ArenaChallenge_Resp, resp)
		return
	}
	if targetPlayer.IsWin {
		logger.Error("has challenged")
		sess.Send(structs.Protocol_ArenaChallenge_Resp, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	sess.PlayerData.Arena.ChallengeCount++

	left := sess.PlayerData.GetFightTeamByHeroTeam()
	var right *structs.FightTeam
	{
		// if targetPlayer.PlayerID > 100000000 {

		// } else {
		robot, ok := gamedata.ArenaRobotsMap[targetPlayer.PlayerID]
		if !ok {
			logger.Error("cannot find robot : %v", targetPlayer.PlayerID)
			sess.Send(structs.Protocol_ArenaChallenge_Resp, resp)
			return
		}
		right = model.GetFightTeamByHeroTemplateIDs(robot.GetHeroIDs())
		right.DefaultHP = targetPlayer.HP
		//}
	}

	randIndex := util.RandNum(int32(0), int32(len(gamedata.AllTemplates.GlobalData.FightWithFriendBattleIDs)))
	battleID := gamedata.AllTemplates.GlobalData.FightWithFriendBattleIDs[randIndex]
	battlefieldTemplate, ok := gamedata.AllTemplates.Battlefields[battleID]
	if !ok {
		logger.Error("cannot find battleField : %v", battleID)
		sess.Send(structs.Protocol_ArenaChallenge_Resp, resp)
		return
	}

	sim := model.NewFightSim(left, right)
	fightRet := sim.Fight()
	fightRet.BackgroundID = battlefieldTemplate.BackgroundID
	fightRet.ForegroundID = battlefieldTemplate.ForegroundID

	if fightRet.IsLeftWin {
		targetPlayer.IsWin = true
	}

	sess.CheckAchievements(structs.AchvCondType_ChallengePlayer, 0, 1)

	fightNtf := &structs.FightResultNtf{
		FType:  structs.FightType_ArenaPK,
		Result: fightRet,
	}
	sess.Send(structs.Protocol_FightResult_Ntf, fightNtf)

	resp.Ret = structs.AdventureRet_Success
	resp.IsWin = targetPlayer.IsWin
	resp.PlayerID = targetPlayer.PlayerID
	sess.Send(structs.Protocol_ArenaChallenge_Resp, resp)

	if fightRet.IsLeftWin {
		winCnt := sess.PlayerData.Arena.GetWinCount()

		index1 := CheckCourageAchievement(winCnt)
		index2 := CheckChampionAchievement(winCnt)
		index := int(math.Max(float64(index1), float64(index2)))

		if index >= 0 && sess.PlayerData.Arena.RewardRecord[index] != structs.ArenaRwdStatus_Recieve {
			// 奖励领取
			sess.PlayerData.Arena.RewardRecord[index] = structs.ArenaRwdStatus_Recieve

			// 勇气点
			addNum := gamedata.AllTemplates.GlobalData.CourageAward[index]
			if addNum > 0 {
				sess.CheckAchievements(structs.AchvCondType_CollectPoint, structs.PointType_Courage, addNum)
			}

			// 冠军点
			addNum = gamedata.AllTemplates.GlobalData.ChampionAward[index]
			if addNum > 0 {
				sess.CheckAchievements(structs.AchvCondType_CollectPoint, structs.PointType_Champion, addNum)
			}
		}

		condID := int32(0)
		if winCnt == 1 {
			condID = structs.PointType_WinArena1Player
		}
		if winCnt == 4 {
			condID = structs.PointType_WinArena4Player
		}
		if winCnt == 9 {
			condID = structs.PointType_WinArena9Player
		}
		sess.CheckAchievements(structs.AchvCondType_CollectPoint, condID, 1)
		sess.CheckAchievements(structs.AchvCondType_WinArenaPlayer, 0, 1)

		//CZXDO: 9 连胜全服公告
	}
}

func RefreshArenaReq(sess *sessions.Session, msgBody []byte) {
	req := &structs.ArenaRefreshReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.ArenaRefreshResp{
		Ret: structs.AdventureRet_Failed,
	}
	/////////////////////////////////////////////Data Check////////////////////////////////////////
	userArena := sess.PlayerData.Arena
	if userArena.NextRefrshTime > time.Now().Unix() {
		if !req.UserIngot {
			return
		}
		if sess.PlayerData.Res.Ingot < gamedata.AllTemplates.GlobalData.RefreshIngot {
			logger.Error("ingot not enough ignot: %v, need ignot : %v", sess.PlayerData.Res.Ingot, gamedata.AllTemplates.GlobalData.RefreshIngot)
			sess.Send(structs.Protocol_ArenaChallenge_Resp, resp)
			return
		}
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	sess.IgnotChange(gamedata.AllTemplates.GlobalData.RefreshIngot, 0)

	userArena.NextRefrshTime = time.Now().AddDate(0, 0, -1).Unix()

	sess.RefreshArenaData()

	userArena.RewardRecord = []structs.ArenaRwdStatus{}

	sess.SyncPlayerArena()
}

func RecieveArenaRewardReq(sess *sessions.Session, msgBody []byte) {
	req := &structs.ArenaRecieveRewardReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.ArenaRecieveRewardResp{}

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	if req.StarIndex >= len(gamedata.AllTemplates.GlobalData.TotalWinCount) {
		logger.Error("start index is error %v", req.StarIndex)
		sess.Send(structs.Protocol_RecieveArenaReward_Resp, resp)
		return
	}

	userArena := sess.PlayerData.Arena
	///////////////////////////////////////////Logic Process///////////////////////////////////////
	userArena.RewardRecord[req.StarIndex] = structs.ArenaRwdStatus_Recieve

	resp.RewardRecords = userArena.RewardRecord
	sess.Send(structs.Protocol_RecieveArenaReward_Resp, resp)

	// 勇气点
	addNum := gamedata.AllTemplates.GlobalData.CourageAward[req.StarIndex]
	sess.CheckAchievements(structs.AchvCondType_CollectPoint, structs.PointType_Courage, addNum)
	sess.CheckAchievements(structs.AchvCondType_CollectPoint, structs.PointType_Champion, addNum)
}
