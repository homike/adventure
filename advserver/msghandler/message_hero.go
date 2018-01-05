package msghandler

import "adventure/common/structs"
import "adventure/advserver/sessions"
import "adventure/advserver/gamedata"

func InitMessageHero() {
	heroMessage := map[uint16]ProcessFunc{
		uint16(structs.Protocol_Employ_Req):            EmployReq,         // 雇佣英雄
		uint16(structs.Protocol_UnEmploy_Req):          UnEmployReq,       // 解雇英雄
		uint16(structs.Protocol_UnEmployManyHeros_Req): TestReq,           // 解雇多名英雄
		uint16(structs.Protocol_ResetHeroIndex_Req):    ResetHeroIndexReq, // 调整英雄出站顺序
		uint16(structs.Protocol_Work_Req):              WorkReq,           // 英雄出战
		uint16(structs.Protocol_SomeWork_Req):          SomeWorkReq,       // 一些英雄出战
		uint16(structs.Protocol_Rest_Req):              ResetReq,          // 英雄休息
		uint16(structs.Protocol_SomeRest_Req):          SomeResetReq,      // 一些英雄休息
		uint16(structs.Protocol_Awake_Req):             AwakeReq,          // 英雄觉醒
		uint16(structs.Protocol_UpgradeWeapon_Req):     UpgradeWeaponReq,  // 武具升级
		uint16(structs.Protocol_SyncEmploy_Req):        SyncEmployReq,     // 同步招募信息
	}

	for k, v := range heroMessage {
		MapFunc[k] = v
	}
}

// 雇佣
func EmployReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("EmployReq")

	req := &structs.EmployReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.EmployResp{
		Ret: structs.EmployRet_Failed,
	}
	/////////////////////////////////////////////Data Check////////////////////////////////////////
	sess.RefreshPlayerInfo(nil)

	switch structs.EmployType(req.EmployType) {
	case structs.EmployType_Money:
	case structs.EmployType_HunLuan:
	case structs.EmployType_HuiHuang:
	case structs.EmployType_LvDong:
	case structs.EmployType_Diamond:
	case structs.EmployType_ManyDiamond:
	case structs.EmployType_ManyDiamond2:
	case structs.EmployType_Exchange:
	case structs.EmployType_Reward:
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	resp.Ret = structs.EmployRet_Success
	resp.HeroIDs = []int32{10110}

	sess.Send(structs.Protocol_Employ_Resp, resp)
}

// 解雇
func UnEmployReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("UnEmployReq")

	req := &structs.UnEmployReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.UnEmployResp{
		Ret:    structs.AdventureRet_Failed,
		HeroID: req.HeroID,
	}
	/////////////////////////////////////////////Data Check////////////////////////////////////////

	hero, err := sess.PlayerData.HeroTeam.GetHero(req.HeroID)
	if err != nil {
		logger.Error("player(%v) has not hero (%v)", sess.AccountID, req.HeroID)
		sess.Send(structs.Protocol_UnEmploy_Resp, resp)
	}

	if hero.IsOutFight {
		logger.Error("player(%v) hero (%v) is out fight, cannot unemploy", sess.AccountID, req.HeroID)
		sess.Send(structs.Protocol_UnEmploy_Resp, resp)
	}

	honorDebris := gamedata.AllTemplates.HeroTemplates[req.HeroID].HonorDebris

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	rewards := []structs.Reward{}
	// 奖励碎片
	rewards = append(rewards, structs.Reward{
		RewardType: structs.RewardType_Property,
		Param1:     3, // 碎片
		Param2:     int32(honorDebris),
	})
	if hero.Exp >= gamedata.EmployReturnExp {
		//deltaExp := int(hero.Exp * gamedata.EmployReturnExpPer / 100)
		// 背包中添加该物品
		rewards = append(rewards, structs.Reward{
			RewardType: structs.RewardType_Item,
			Param1:     6, // 英雄志道具ID
			Param2:     1,
		})
	}

	sess.PlayerData.HeroTeam.RemoveHero(hero)

	sess.Send(structs.Protocol_UnEmploy_Resp, resp)
}

// 调整英雄顺序
func ResetHeroIndexReq(sess *sessions.Session, msgBody []byte) {

	logger.Debug("ResetHeroIndexReq")

	req := &structs.ResetHeroIndexReq{}
	sess.UnMarshal(msgBody, req)

	newIndex := int32(0)
	for _, v := range req.HeroIDs {
		hero, err := sess.PlayerData.HeroTeam.GetHero(v)
		if err == nil {
			hero.Index = newIndex
			newIndex++
		}
	}

	sess.PlayerData.HeroTeam.SortHeros()
}

// 上阵
func workResp(sess *sessions.Session, heroID int32) {
	resp := &structs.WorkResp{
		Ret:    structs.AdventureRet_Failed,
		HeroID: heroID,
	}
	hero, err := sess.PlayerData.HeroTeam.GetHero(heroID)
	if err != nil {
		logger.Error("cannot find hero %v, err: %v", heroID, err)
		sess.Send(structs.Protocol_Work_Resp, resp)
		return
	}

	hero.IsOutFight = true

	resp.Ret = structs.AdventureRet_Success
	sess.Send(structs.Protocol_Work_Resp, resp)
}

func WorkReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("WorkReq")

	req := &structs.WorkReq{}
	sess.UnMarshal(msgBody, req)

	workResp(sess, req.HeroID)
}

func SomeWorkReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("SomeWorkReq")

	req := &structs.SomeWorkReq{}
	sess.UnMarshal(msgBody, req)

	for _, v := range req.HeroIDs {
		workResp(sess, v)
	}
}

// 下阵
func resetResp(sess *sessions.Session, heroID int32) {
	resp := &structs.ResetResp{
		Ret:    structs.AdventureRet_Failed,
		HeroID: heroID,
	}
	hero, err := sess.PlayerData.HeroTeam.GetHero(heroID)
	if err != nil {
		sess.Send(structs.Protocol_Rest_Resp, resp)
	}
	hero.IsOutFight = false
	sess.Send(structs.Protocol_Rest_Resp, resp)
}

func ResetReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("ResetReq")

	req := &structs.ResetReq{}
	sess.UnMarshal(msgBody, req)

	workResp(sess, req.HeroID)
}

func SomeResetReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("SomeResetReq")

	req := &structs.SomeResetReq{}
	sess.UnMarshal(msgBody, req)

	for _, v := range req.HeroIDs {
		workResp(sess, v)
	}
}

// 觉醒
func AwakeReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("AwakeReq")

	req := &structs.AwakeReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.AwakeResp{
		Ret:    structs.AdventureRet_Failed,
		HeroID: req.HeroID,
		AddHP:  0,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	hero, err := sess.PlayerData.HeroTeam.GetHero(req.HeroID)
	if err != nil {
		logger.Error("cannot find the hero(%v)", req.HeroID)
		sess.Send(structs.Protocol_Awake_Resp, resp)
		return
	}

	if hero.Level < gamedata.HeroAwakeMinLevel {
		logger.Error("hero(%v) level not enough", req.HeroID)
		sess.Send(structs.Protocol_Awake_Resp, resp)
		return
	}

	MaxAwakeCnt := gamedata.AllTemplates.HeroTemplates[req.HeroID].AwakeCount

	if hero.AwakeCount >= int32(MaxAwakeCnt) {
		logger.Error("hero(%v) awake max", req.HeroID)
		sess.Send(structs.Protocol_Awake_Resp, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////

	//CZXDO: 扣除金币及巨魔雕像
	hero.AwakeCount++
	hero.Level = 1
	nextLevelExp, err := gamedata.GetHeroLevelExp(hero.Level-1, hero.AwakeCount-1)
	if err != nil {
		logger.Error("hero(%v) GetHeroLevelExp error(%v)", req.HeroID, err)
		sess.Send(structs.Protocol_Awake_Resp, resp)
		return
	}

	// 英雄升级
	for {
		if hero.Exp < nextLevelExp {
			break
		}
		hero.Level++
		nextLevelExp, err = gamedata.GetHeroLevelExp(hero.Level-1, hero.AwakeCount-1)
		if err != nil {
			logger.Error("hero(%v) GetHeroLevelExp error(%v)", req.HeroID, err)
			sess.Send(structs.Protocol_Awake_Resp, resp)
			break
		}
	}

	oldHp := hero.HP()
	// 重新计算该英雄等级战力
	sess.PlayerData.HeroTeam.ReCalculateHeroLevelHp(hero)
	// 同步该英雄信息
	sess.SyncHeroNtf(structs.SyncHeroType_Update, []*structs.Hero{hero})

	resp.Ret = structs.AdventureRet_Success
	resp.HeroID = req.HeroID
	resp.AddHP = hero.HP() - oldHp
	sess.Send(structs.Protocol_Awake_Resp, resp)
}

// 武具升级
func UpgradeWeaponReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("UpgradeWeaponReq")

	req := &structs.UpgradeWeaponReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.UpgradeWeaponResp{
		Ret:    structs.AdventureRet_Failed,
		HeroID: req.HeroID,
		AddHP:  0,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	hero, err := sess.PlayerData.HeroTeam.GetHero(req.HeroID)
	if err != nil {
		logger.Error("cannot find the hero(%v)", req.HeroID)
		sess.Send(structs.Protocol_UpgradeWeapon_Resp, resp)
		return
	}

	if sess.PlayerData.Res.Ingot < req.Ingot {
		logger.Error("Ingot not enough")
		sess.Send(structs.Protocol_UpgradeWeapon_Resp, resp)
		return
	}

	if hero.IsPlayer {
		logger.Error("hero is player cannot update")
		sess.Send(structs.Protocol_UpgradeWeapon_Resp, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////

	// CZXDO: 消耗资源
	hero.WeaponLevel++
	oldHP := hero.HP()

	// 重新计算该英雄等级战力
	sess.PlayerData.HeroTeam.ReCalculateHeroLevelHp(hero)
	// 同步该英雄信息
	sess.SyncHeroNtf(structs.SyncHeroType_Update, []*structs.Hero{hero})

	resp.Ret = structs.AdventureRet_Success
	resp.HeroID = req.HeroID
	resp.AddHP = hero.HP() - oldHP
	sess.Send(structs.Protocol_UpgradeWeapon_Resp, resp)

	// CZXDO: 全服通告
}

func SyncEmployReq(sess *sessions.Session, msgBody []byte) {
	sess.SyncEmployInfo()
}
