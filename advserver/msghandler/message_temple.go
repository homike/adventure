package msghandler

import (
	"adventure/advserver/gamedata"
	"adventure/advserver/sessions"
	"adventure/common/structs"
)

func InitMessageTemple() {
	message := map[uint16]ProcessFunc{
		uint16(structs.Protocol_SyncTemplateHeros_Req):    SyncTemplateHerosReq,    // 同步神殿英雄数据
		uint16(structs.Protocol_ExchangeTemplateHero_Req): ExchangeTemplateHeroReq, // 兑换英雄
		uint16(structs.Protocol_UnlockTemple_Req):         UnlockTempleReq,         // 解锁神殿
		uint16(structs.Protocol_RefreshTemple_Req):        RefreshTempleReq,        // 刷新神殿
	}

	for k, v := range message {
		MapFunc[k] = v
	}
}

func SyncTemplateHerosReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("SyncTemplateHerosReq")

	sess.SyncTemplateHeros()
}

func ExchangeTemplateHeroReq(sess *sessions.Session, msgBody []byte) {
	req := &structs.ExchangeTemplateHeroReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.ExchangeTemplateHeroResp{
		Ret:            structs.AdventureRet_Failed,
		HeroTemplateID: req.HeroTemplateID,
	}
	respID := uint16(structs.Protocol_ExchangeTemplateHero_Resp)

	/////////////////////////////////////////////Data Check////////////////////////////////////////

	vipT, ok := gamedata.AllTemplates.VipTemplates[sess.PlayerData.VipLevel]
	if !ok {
		logger.Error("cannot find vip (%v) template", sess.PlayerData.VipLevel)
		sess.Send(respID, resp)
		return
	}

	userTemple := sess.PlayerData.Temple
	if userTemple.ToDayTradeCount >= vipT.DayofTradeTempleHero {
		logger.Error(" exchange count %v is limit %v", userTemple.ToDayTradeCount, vipT.DayofTradeTempleHero)
		sess.Send(respID, resp)
		return
	}

	var exchangeHero *structs.TempleHero
	for _, v := range userTemple.TempleHeros {
		if v.HeroTemplateID == req.HeroTemplateID {
			exchangeHero = v
		}
	}
	if exchangeHero == nil {
		logger.Error("cannot find exchange hero(%v)", req.HeroTemplateID)
		sess.Send(respID, resp)
		return
	}
	if exchangeHero.IsTrade {
		logger.Error("hero(%v) has exchanged", req.HeroTemplateID)
		sess.Send(respID, resp)
		return
	}

	heroT, ok := gamedata.AllTemplates.HeroTemplates[exchangeHero.HeroTemplateID]
	if !ok {
		logger.Error("hero(%v) cannot find template", req.HeroTemplateID)
		sess.Send(respID, resp)
		return
	}

	if sess.PlayerData.Res.Fragments < heroT.EmployCostFragments {
		logger.Error("hero(%v) need fragments %v, user fragements %v", heroT.EmployCostFragments, sess.PlayerData.Res.Fragments)
		sess.Send(respID, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	// 碎片扣除
	sess.FregmentsChange(-heroT.EmployCostFragments, 0)

	// 兑换英雄
	sess.PlayerData.HeroTeam.AddHero("", false, heroT.ID)
	userTemple.ToDayTradeCount++
	exchangeHero.IsTrade = true

	// 返回消息
	resp.Ret = structs.AdventureRet_Success
	sess.Send(respID, resp)

	// 发送奖励
	reward := &structs.Reward{
		RewardType: structs.RewardType_Hero,
		Param1:     heroT.ID,
	}
	sess.RewardResults(false, []*structs.Reward{reward}, "")
}

func UnlockTempleReq(sess *sessions.Session, msgBody []byte) {

	resp := structs.UnlockTempleResp{
		Ret: structs.AdventureRet_Failed,
	}

	for _, v := range sess.PlayerData.MenuStates {
		if v.MenuID == int32(structs.MenuTypes_Temple) && v.MenuStatus != structs.MenuStatus_Close {
			return
		}
	}

	money := gamedata.AllTemplates.GlobalData.TempleUnlockCost
	if sess.PlayerData.Res.Money < money {
		logger.Error("money not enough, money : %v, need money: %v", sess.PlayerData.Res.Money, money)
		sess.Send(structs.Protocol_UnlockTemple_Req, resp)
		return
	}

	sess.MoneyChange(-money, 0)
	sess.UnLockMenu(int32(structs.MenuTypes_Temple))

	sess.Send(structs.Protocol_UnlockTemple_Resp, resp)
}

func RefreshTempleReq(sess *sessions.Session, msgBody []byte) {

	req := &structs.RefreshTempleReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.RefreshTempleResp{
		Ret:        structs.AdventureRet_Failed,
		SplashGold: false,
	}
	respID := uint16(structs.Protocol_RefreshTemple_Resp)

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	userTemple := sess.PlayerData.Temple
	if userTemple.RefreshCount != req.Count {
		logger.Error("count %v not equal refreshCount: %v", req.Count, userTemple.RefreshCount)
		sess.Send(respID, resp)
		return
	}

	arrlen := len(gamedata.AllTemplates.GlobalData.TempleRefreshIngot)
	ingot := gamedata.AllTemplates.GlobalData.TempleRefreshIngot[arrlen-1]
	if req.Count < int32(arrlen) {
		ingot = gamedata.AllTemplates.GlobalData.TempleRefreshIngot[req.Count]
	}

	if sess.PlayerData.Res.Ingot < ingot {
		logger.Error("req ingot: %v, use ingot: %v, not enough", ingot, sess.PlayerData.Res.Ingot)
		sess.Send(respID, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////

	sess.PlayerData.Temple.RefreshHeros(true)

	sess.PlayerData.Temple.RefreshCount++

	sess.IgnotChange(-ingot, 0)

	sess.SyncTemplateHeros()

	sess.CheckAchievements(structs.AchvCondType_CollectPoint, structs.PointType_RefreshTemple, 1)

	splashHero := &structs.TempleHero{}
	for _, v := range userTemple.TempleHeros {
		if v.Quality == structs.QualityType_SplashGold {
			splashHero = v
		}
	}

	resp.SplashGold = (splashHero.Quality == structs.QualityType_SplashGold)
	sess.Send(respID, resp)

	//CZXDO: 神殿英雄公告
}
