package sessions

import (
	"adventure/common/structs"
	"fmt"
)

func (sess *Session) SyncGameLevelNtf() {

	resp := &structs.SyncGameLevelNtf{
		GameLevels:      sess.PlayerData.PlayerGameLevel.GameLevels,
		CurLevelID:      sess.PlayerData.PlayerGameLevel.CurrentGameLevelID,
		LastRefreshTime: sess.PlayerData.PlayerGameLevel.LastRefreshTime,
		SpeedCount:      sess.PlayerData.PlayerGameLevel.TodaySpeedAdventure,
	}
	fmt.Println(" GameLevels len ", len(sess.PlayerData.PlayerGameLevel.GameLevels))
	fmt.Println(sess.PlayerData.PlayerGameLevel.GameLevels[0].IsUnlock, " event len", len(sess.PlayerData.PlayerGameLevel.GameLevels[0].CompleteEvent))
	fmt.Println("gameLeveel time : ", sess.PlayerData.PlayerGameLevel.LastRefreshTime, " CurLevelID", resp.CurLevelID)

	sess.Send(structs.Protocol_SyncGameLevel_Ntf, resp)
}

func (sess *Session) SyncCurrentGameLevelNtf() {
	gameLevel, err := sess.PlayerData.PlayerGameLevel.GetCurGameLevel()
	if err != nil {
		return
	}

	resp := &structs.SyncCurrentGameLevelNtf{
		GameLevel: gameLevel,
	}
	sess.Send(structs.Protocol_SyncCurrentGameLevel_Ntf, resp)
}
