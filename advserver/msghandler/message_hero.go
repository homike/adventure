package msghandler

import "adventure/common/structs"
import "adventure/advserver/sessions"
import "adventure/advserver/gamedata"

func InitMessageHero() {
	heroMessage := map[uint16]ProcessFunc{
		uint16(structs.Protocol_Employ_Req):            TestReq, // 雇佣英雄
		uint16(structs.Protocol_UnEmploy_Req):          TestReq, // 解雇英雄
		uint16(structs.Protocol_UnEmployManyHeros_Req): TestReq, // 解雇多名英雄
		uint16(structs.Protocol_ResetHeroIndex_Req):    TestReq, // 调整英雄出站顺序
		uint16(structs.Protocol_Work_Req):              TestReq, // 英雄出战
		uint16(structs.Protocol_SomeWork_Req):          TestReq, // 一些英雄出战
		uint16(structs.Protocol_Rest_Req):              TestReq, // 英雄休息
		uint16(structs.Protocol_SomeRest_Req):          TestReq, // 一些英雄休息
		uint16(structs.Protocol_Awake_Rep):             TestReq, // 英雄觉醒
		uint16(structs.Protocol_UpgradeWeapon_Rep):     TestReq, // 武具升级
		uint16(structs.Protocol_SyncEmploy_Req):        TestReq, // 同步招募信息
	}

	for k, v := range heroMessage {
		MapFunc[k] = v
	}
}

func EmployReq(sess sessions.Session, msgBody []byte) {
	logger.Debug("EmployReq")

	req := &structs.EmployReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.EmployResp{
		Ret: structs.EmployRet_Failed,
	}
	/////////////////////////////////////////////Data Check////////////////////////////////////////
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

func UnEmployReq(sess sessions.Session, msgBody []byte) {
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

	honorDebris, err := gamedata.AllTemplates.HeroTemplate.HonorDebris(req.HeroID)
	if err != nil {
		logger.Error("player(%v) hero (%v) get gamedata error: %v", sess.AccountID, req.HeroID, err)
		sess.Send(structs.Protocol_UnEmploy_Resp, resp)
	}

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
