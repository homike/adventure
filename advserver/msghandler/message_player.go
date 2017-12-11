package msghandler

import (
	"adventure/advserver/model"
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"fmt"
)

// 1002
func CreatePlayer(sess *sessions.Session, msgBody []byte) {
	fmt.Println("czx@@@ CreatePlayer:", string(msgBody))

	req := structs.CreatePlayerReq{}
	sess.UnMarshal(msgBody, &req)

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
	fmt.Println("czx@@@ NameExists:", string(msgBody))

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
	fmt.Println("czx@@@ SyncUserGuidRecords:")

	records := []structs.GuildRecord{}
	for i := 0; i < 2; i++ {
		records = append(records, structs.GuildRecord{
			UserGuidTypes: uint8(i + 2),
			TriggerCount:  int32(i + 3),
		})
	}
	resp := &structs.SyncUserGuidRecordsNtf{
		Records: records,
	}

	sess.Send(structs.Protocol_SyncUserGuidRecords_Ntf, resp)
}
