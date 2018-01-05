package msghandler

import (
	"adventure/advserver/gamedata"
	"adventure/advserver/model"
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"adventure/common/util"
	"fmt"
	"time"
)

// 1002
func CreatePlayer(sess *sessions.Session, msgBody []byte) {
	//fmt.Println("CreatePlayer data:", msgBody)

	req := structs.CreatePlayerReq{}
	sess.UnMarshal(msgBody, &req)

	//fmt.Println("CreatePlayer name: ", req.PlayerName, "heroTemplateID: ", req.HeroTemplateId)
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

	sess.OnEnterGame()

	SyncLoginDataFinish(sess)
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

func EatFood(sess *sessions.Session, msgBody []byte) {

	req := structs.EatFoodReq{}
	sess.UnMarshal(msgBody, &req)

	//fmt.Println("CreatePlayer name: ", req.PlayerName, "heroTemplateID: ", req.HeroTemplateId)
	resp := &structs.EatFoodResp{
		Ret:      structs.AdventureRet_Failed,
		Strength: sess.PlayerData.Res.Strength,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	sess.RefreshPlayerInfo(nil)

	if !sess.PlayerData.Res.Foods.Exist(req.FoodID, 1) {
		logger.Error("has not enough food(%v, %v)", req.FoodID, 1)
		sess.Send(structs.Protocol_EatFood_Resp, resp)
		return
	}

	foodT, ok := gamedata.AllTemplates.ResourceTemplates[req.FoodID]
	if !ok {
		logger.Error("cannot find food(%v) template", req.FoodID)
		sess.Send(structs.Protocol_EatFood_Resp, resp)
		return
	}

	if sess.PlayerData.Res.Strength > foodT.MinEatLimit && foodT.MinEatLimit != 0 {
		logger.Error("Strength: %v, MinEatLimit: %v", sess.PlayerData.Res.Strength, foodT.MinEatLimit)
		sess.Send(structs.Protocol_EatFood_Resp, resp)
		return
	}

	eatDate := sess.PlayerData.ExtendData.GetEatedFoodRecord(req.FoodID)
	if foodT.IsOnceEveryday && util.TimeSub(time.Now(), time.Unix(eatDate, 0)) == 0 {
		logger.Error("already eat today (%v)", req.FoodID)
		sess.Send(structs.Protocol_EatFood_Resp, resp)
		return
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////

	if foodT.IsOnceEveryday {
		sess.PlayerData.ExtendData.AddEatedFoodRecord(req.FoodID)
		sess.SyncEatFoodList()
	}
	// minus food
	sess.PlayerData.Res.Foods.Add(req.FoodID, -1)
	// strength change
	sess.PlayerData.Res.StrengthChange(foodT.AddActive)
	// add rewards
	if len(foodT.RewardIDs) > 0 {
		sess.DoSomeRewards(foodT.RewardIDs)
	}

	resp.Ret = structs.AdventureRet_Success
	resp.Strength = sess.PlayerData.Res.Strength
	sess.Send(structs.Protocol_EatFood_Resp, resp)
}

func GetEatedFoodList(sess *sessions.Session, msgBody []byte) {
	sess.SyncEatFoodList()
}
