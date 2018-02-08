package sessions

import "adventure/common/structs"

func (sess *Session) SyncMineMapData() {
	sess.Send(structs.Protocol_SyncMiningMap_Resp, &structs.SyncMiningMapResp{
		sess.PlayerData.MineMap.MapData.NodeList,
		sess.PlayerData.MineMap.UserData})
}

// 添加资源
func (sess *Session) AddBlock(nodeID int32, x, y int8, isVisible bool) {
	// userMine := sess.PlayerData.MineMap

	// _, err := userMine.AddBlock(nodeID, x, y, isVisible)
	// if err != nil {
	// }

	// // 如果是巨魔雕像，同步客户端
	// tileT, ok := gamedata.AllTemplates.TileTemplates[nodeID]
	// if ok && tileT.TileType == structs.TileType_Statue {
	// 	userMine.StatueCnt++
	// 	sess.Send(structs.Protocol_UpdateStatueLevelAndCount_Ntf, &structs.UpdateStatueLevelAndCountNtf{
	// 		Level: userMine.StatueLv,
	// 		Count: userMine.StatueLv,
	// 	})
	// }
}
