package msghandler

import (
	"adventure/advserver/model"
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"fmt"
)

// 1002
func CreatePlayer(sess *sessions.Session, msgBody []byte) {
	fmt.Println("CreatePlayer data:", msgBody)

	req := structs.CreatePlayerReq{}
	sess.UnMarshal(msgBody, &req)

	fmt.Println("CreatePlayer name: ", req.PlayerName, "heroTemplateID: ", req.HeroTemplateId)
	resp := &structs.CreatePlayerResp{
		Result: 0, // Success
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	///////////////////////////////////////////Logic Process///////////////////////////////////////
	player, err := model.NewPlayer(req.PlayerName, req.HeroTemplateId)
	if err != nil {
		fmt.Println("NewPlayer Error", err)
		resp.Result = 1
		sess.Send(structs.Protocol_CreatePlayer_Resp, resp)
		return
	}
	sess.SetPlayer(player)
	sessions.SessionMgr.AddSession(sess)

	sess.Send(structs.Protocol_CreatePlayer_Resp, resp)

	SyncPlayerBaseInfo(sess)

	SyncUserGuidRecords(sess)

	SyncLoginDataFinish(sess)

	OnEnterGame(sess)
}

// 1008
func SyncPlayerBaseInfo(sess *sessions.Session) {
	fmt.Println("czx@@@ SyncPlayerBaseInfo")

	resp := &structs.SyncPlayerBaseInfoNtf{
		PlayerID:           1,
		GameZoonID:         1,
		IsSupperMan:        true,
		PlatformType:       1,
		Viplevel:           1,
		TotalRechargeIngot: 1,
	}
	sess.Send(structs.Protocol_SyncPlayerBaseInfo_Ntf, resp)
}

// 1009
func NameExists(sess *sessions.Session, msgBody []byte) {
	fmt.Println("czx@@@ NameExists1:", string(msgBody))

	req := structs.NameExistsReq{}
	sess.UnMarshal(msgBody, &req)

	resp := &structs.NameExistsResp{
		Name: req.Name,
	}
	sess.Send(structs.Protocol_NameExists_Resp, resp)
}

func UpdateUserGuidRecord(sess *sessions.Session, msgBody []byte) {
	fmt.Println("czx@@@ UpdateUserGuidRecord:", string(msgBody))

	req := structs.UpdateUserGuidRecordReq{}
	sess.UnMarshal(msgBody, &req)

	sess.PlayerData.UpdateGuidRecords(req.UserGuidTypes)
}

// 1413
func SyncUserGuidRecords(sess *sessions.Session) {
	//fmt.Println("czx@@@ SyncUserGuidRecords:")

	records := []structs.GuildRecord{}
	for i := 0; i < 24; i++ {
		records = append(records, structs.GuildRecord{
			UserGuidTypes: uint8(i),
			TriggerCount:  int32(5),
		})
	}
	resp := &structs.SyncUserGuidRecordsNtf{
		Records: records,
	}

	sess.Send(structs.Protocol_SyncUserGuidRecords_Ntf, resp)
}

func OnEnterGame(sess *sessions.Session) {
	SyncHeroNtf(sess) // 同步英雄信息

	SyncStrength(sess)        // 同步饱食度
	SyncHeroWorkTop(sess)     // 同步最大出战英雄数
	SyncUserGuidRecords(sess) // 同步新手引导进度
	//同步客户端已食用过的食物列表
	SyncUnlockMenus(sess)      //同步客户端已解锁菜单列表
	SyncGameBoxTopNumNtf(sess) //同步客户端附加的宝箱上限数量
	//同步购买商城物品记录
}

func SyncHeroNtf(sess *sessions.Session) {
	fmt.Println("SyncHeroNtf heros num : ", len(sess.PlayerData.HeroTeam.Heros))
	resp := &structs.SyncHeroNtf{
		SyncHeroType: structs.SyncHeroType_First,
		Heros:        sess.PlayerData.HeroTeam.GetHerosArray(),
	}
	sess.Send(structs.Protocol_SyncHero_Ntf, resp)
}

func SyncStrength(sess *sessions.Session) {
	resp := &structs.SyncStrengthNtf{
		Strength: sess.PlayerData.Res.Strength,
	}
	sess.Send(structs.Protocol_SyncStrength_Ntf, resp)
}

func SyncHeroWorkTop(sess *sessions.Session) {
	resp := &structs.SyncHeroWorkTopNtf{
		MaxWorker: sess.PlayerData.HeroTeam.MaxWorker,
	}
	sess.Send(structs.Protocol_SyncWorkHeroTop_Ntf, resp)
}

func SyncUnlockMenus(sess *sessions.Session) {
	resp := &structs.SyncUnlockMenusNtf{
		MenuStates: sess.PlayerData.MenuStates,
	}
	sess.Send(structs.Protocol_SyncUnlockMenus_Ntf, resp)
}

func SyncGameBoxTopNumNtf(sess *sessions.Session) {
	resp := &structs.SyncGameBoxTopNumNtf{
		AddNum: sess.PlayerData.AddGameBoxCount,
	}
	sess.Send(structs.Protocol_SyncGameBoxTopNum_Ntf, resp)
}
