package msghandler

import (
	"adventure/advserver/gamedata"
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"math"
)

func InitMessageArtifact() {
	message := map[uint16]ProcessFunc{
		uint16(structs.Protocol_EquipArtifact_Req):   EquipArtifactReq,   // 装备神器
		uint16(structs.Protocol_UpgradeArtifact_Req): UpgradeArtifactReq, // 升级神器
	}

	for k, v := range message {
		MapFunc[k] = v
	}
}

func EquipArtifactReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("EquipArtifactReq")

	req := &structs.EquipArtifactReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.EquipArtifactResp{
		Ret: structs.AdventureRet_Failed,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////

	aStatus := sess.PlayerData.Artifact.GetArtifactStatus(req.ArtifactID)
	if aStatus == nil {
		logger.Error("GetArtifactStatus(%v) Error", req.ArtifactID)
		sess.Send(structs.Protocol_EquipArtifact_Resp, resp)
		return
	}
	if aStatus.Status == structs.UnLock {
		logger.Error("%v is unlock", req.ArtifactID)
		sess.Send(structs.Protocol_EquipArtifact_Resp, resp)
		return
	}
	if aStatus.Status == structs.Use {
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////

	eStatusUse := sess.PlayerData.Artifact.GetArtifactStatusUse()
	if eStatusUse != nil {
		eStatusUse.Status = structs.UnLock
	}

	aStatus.Status = structs.Use

	resp.Ret = structs.AdventureRet_Success
	sess.Send(structs.Protocol_EquipArtifact_Resp, resp)
}

func UpgradeArtifactReq(sess *sessions.Session, msgBody []byte) {
	logger.Debug("EquipArtifactReq")

	req := &structs.UpgradeArtifactReq{}
	sess.UnMarshal(msgBody, req)

	resp := &structs.UpgradeArtifactResp{
		Ret:   structs.AdventureRet_Failed,
		AddHP: 0,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////

	unlockCnt := sess.PlayerData.Artifact.GetArtifactStatusUnLockCount()
	palyerHero := sess.PlayerData.HeroTeam.GetPlayerHero()
	if palyerHero == nil {
		logger.Error("GetPlayerHero() error")
		sess.Send(structs.Protocol_UpgradeArtifact_Req, resp)
		return
	}

	if palyerHero.WeaponLevel >= unlockCnt {
		logger.Error("cannot unlock weaplevel: %v, unlockCnt: %v", palyerHero.WeaponLevel, unlockCnt)
		sess.Send(structs.Protocol_UpgradeArtifact_Req, resp)
		return
	}

	if sess.PlayerData.Res.Ingot < req.Ingot {
		logger.Error("Ingot not enough, %v, %v ", sess.PlayerData.Res.Ingot, req.Ingot)
		sess.Send(structs.Protocol_UpgradeArtifact_Req, resp)
		return
	}

	costTemplate, ok := gamedata.AllTemplates.UpgradeArtifactCosts[palyerHero.WeaponLevel+1]
	if !ok {
		logger.Error("UpgradeArtifactCosts() failed", palyerHero.WeaponLevel+1)
		sess.Send(structs.Protocol_UpgradeArtifact_Req, resp)
		return
	}

	useIngot := float32(0)
	costRes := make(map[int32]int32)
	for i := 0; i < len(costTemplate.NeedResourceIdList); i++ {
		resID := costTemplate.NeedResourceIdList[i]
		resCnt := costTemplate.NeedResourceCountList[i]

		userResCnt := sess.PlayerData.Res.GetOresCount(resID)
		if userResCnt < resCnt {
			costRes[resID] = userResCnt
			resT, ok := gamedata.AllTemplates.ResourceTemplates[resID]
			if ok {
				useIngot += resT.UpgradeWeaponCost * float32(resCnt-userResCnt)
			}
		} else {
			costRes[resID] = resCnt
		}
	}

	if useIngot > 0 && req.Ingot == 0 {
		logger.Error("ignot not enough, cannot update")
		sess.Send(structs.Protocol_UpgradeArtifact_Req, resp)
		return
	}

	useIngot2 := math.Ceil(float64(useIngot))

	if int32(useIngot2) != req.Ingot {
		logger.Error("use ignot %v not equal req ignot %v ", useIngot2, req.Ingot)
		sess.Send(structs.Protocol_UpgradeArtifact_Req, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	// 扣除资源
	for k, v := range costRes {
		sess.PlayerData.Res.OresChange(k, -v)
	}

	// 扣除元宝
	if req.Ingot > 0 {
		sess.PlayerData.Res.Ingot -= req.Ingot
	}

	oldHP := palyerHero.HP()

	palyerHero.WeaponLevel++

	sess.PlayerData.HeroTeam.ReCalculateHeroLevelHp(palyerHero)
	// 同步英雄数据
	sess.SyncHeroNtf(structs.SyncHeroType_Update, []*structs.Hero{palyerHero})
	// 升级神器结果
	resp.Ret = structs.AdventureRet_Success
	resp.AddHP = palyerHero.HP() - oldHP
	sess.Send(structs.Protocol_UpgradeArtifact_Req, resp)

	//CZXDO: 成就
}
