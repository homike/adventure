package msghandler

import (
	"adventure/advserver/gamedata"
	"adventure/advserver/sessions"
	"adventure/common/structs"
	"adventure/common/util"
)

func InitMessageMine() {
	handlers := map[uint16]ProcessFunc{
		uint16(structs.Protocol_SyncMiningMap_Req):  SyncMiningMapReq,
		uint16(structs.Protocol_ResetMiningMap_Req): ResetMiningMapReq,
	}

	for k, v := range handlers {
		MapFunc[k] = v
	}
}

func SyncMiningMapReq(sess *sessions.Session, msgBody []byte) {
	sess.SyncMineMapData()
}

func ResetMiningMapReq(sess *sessions.Session, msgBody []byte) {
	req := structs.ResetMiningMapReq{}
	sess.UnMarshal(msgBody, &req)

	resp := &structs.ResetMiningMapResp{
		Ret: structs.AdventureRet_Failed,
	}

	/////////////////////////////////////////////Data Check////////////////////////////////////////
	shopPrice := int32(0)
	seconds := util.TimeSubNow(sess.PlayerData.MineMap.UserData.LastResetDate) - 1
	if seconds < gamedata.AllTemplates.GlobalData.MineMapResetDelayTime {
		if !req.UserIgnot {
			sess.Send(structs.Protocol_ResetMiningMap_Resp, resp)
			return
		}
		shopT, ok := gamedata.AllTemplates.ShopTemplates[int32(structs.ShopEnum_MiningMapReset)]
		if !ok {
			sess.Send(structs.Protocol_ResetMiningMap_Resp, resp)
			return
		}
		if sess.PlayerData.Res.Ingot < shopT.ShopPrice {
			sess.Send(structs.Protocol_ResetMiningMap_Resp, resp)
			return
		}
	}

	///////////////////////////////////////////Logic Process///////////////////////////////////////
	if shopPrice > 0 {
		sess.IgnotChange(shopPrice, 0)
	}

	userMine := sess.PlayerData.MineMap
	userMine.UserData.DigNode.NodeID = -1
	userMine.UserData.DigNode.X = int8(gamedata.AllTemplates.GlobalData.MineMapWidth / 2)
	userMine.UserData.DigNode.Y = 0
	userMine.UserData.Boss.Status = structs.BossStatus_NoAppear
	userMine.UserData.StatueLv = 1
	userMine.UserData.StatueCnt = 0

	userMine.MapData.BossIDs = []int32{}
	userMine.ResetMineMap(false)

	//  同步客户端巨魔数量
	sess.Send(structs.Protocol_UpdateStatueLevelAndCount_Ntf, &structs.UpdateStatueLevelAndCountNtf{userMine.UserData.StatueLv, userMine.UserData.StatueCnt})

	// 重置完成
	resp.Ret = structs.AdventureRet_Success
	sess.Send(structs.Protocol_ResetMiningMap_Resp, resp)

	// 同步地图数据
	sess.SyncMineMapData()
}
