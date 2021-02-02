package sessions

import (
	"adventure/common/structs"
)

func (sess *Session) SyncGameLevelNtf() {

	resp := &structs.SyncGameLevelNtf{
		GameLevels:      sess.PlayerData.PlayerGameLevel.GameLevels,
		CurLevelID:      sess.PlayerData.PlayerGameLevel.CurrentGameLevelID,
		LastRefreshTime: sess.PlayerData.PlayerGameLevel.LastRefreshTime,
		SpeedCount:      sess.PlayerData.PlayerGameLevel.TodaySpeedAdventure,
	}

	//resp.GameLevels[0].CompleteEvent = []uint8{0, 0}
	// fmt.Println("GameLevels len ", len(sess.PlayerData.PlayerGameLevel.GameLevels), " event len", len(sess.PlayerData.PlayerGameLevel.GameLevels[0].CompleteEvent))
	// fmt.Println(sess.PlayerData.PlayerGameLevel.GameLevels[0].GameLevelID, sess.PlayerData.PlayerGameLevel.GameLevels[0].IsUnlock, sess.PlayerData.PlayerGameLevel.GameLevels[0].CompleteEvent)
	// fmt.Println("gameLeveel time : ", sess.PlayerData.PlayerGameLevel.LastRefreshTime, " CurLevelID", resp.CurLevelID)

	sess.Send(structs.Protocol_SyncGameLevel_Ntf, resp)
}

func (sess *Session) SyncCurrentGameLevelNtf() {
	gameLevel, err := sess.PlayerData.PlayerGameLevel.GetCurGameLevelData()
	if err != nil {
		return
	}

	resp := &structs.SyncCurrentGameLevelNtf{
		GameLevel: gameLevel,
	}
	sess.Send(structs.Protocol_SyncCurrentGameLevel_Ntf, resp)
}

func (sess *Session) SyncCurrentGameLevelNtf2() {
	gameLevel, err := sess.PlayerData.PlayerGameLevel.GetCurGameLevelData()
	if err != nil {
		return
	}
	resp := &structs.SyncCurrentGameLevelNtf{
		GameLevel: gameLevel,
	}
	sess.Send(structs.Protocol_SyncCurrentGameLevel_Ntf, resp)
}

// func (sess *Session) gameLevelEventForNone(eventID int32, playerGameLevel *structs.GameLevel, eventTemplate *structs.GameLevelEventTemplate) {

// }
