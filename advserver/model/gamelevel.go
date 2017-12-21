package model

import (
	"adventure/common/structs"
	"errors"
	"time"
)

type PlayerGameLevel struct {
	GameLevels                    []*structs.GameLevel // 游戏关卡列表
	CurrentGameLevelID            int32                // 当前正在进行的关卡
	LastRefreshTime               int64                // 最后一次刷新的时间
	TodaySpeedAdventure           int32                // 今天的加速冒险次数
	NextRefreshSpeedAdventureTime int64                // 下一次刷新冒险次数的时间
}

func NewPlayerGameLevel() *PlayerGameLevel {
	gameLevel := &PlayerGameLevel{
		CurrentGameLevelID: 1,
		LastRefreshTime:    time.Now().Unix(),
		GameLevels:         make([]*structs.GameLevel, 0, 10),
	}

	return gameLevel
}

func (pgl *PlayerGameLevel) AddGameLevel(g *structs.GameLevel) error {
	pgl.GameLevels = append(pgl.GameLevels, g)
	return nil
}

func (pgl *PlayerGameLevel) GetCurGameLevelData() (*structs.GameLevel, error) {
	return pgl.GetGameLevelData(pgl.CurrentGameLevelID)
}

func (pgl *PlayerGameLevel) GetGameLevelData(levelID int32) (*structs.GameLevel, error) {
	for _, v := range pgl.GameLevels {
		if v.GameLevelID == levelID {
			return v, nil
		}
	}
	return nil, errors.New("GetGameLevelData() Error")
}

func (pgl *PlayerGameLevel) SelectGameLevel(levelID int32) error {
	gameLevel, err := pgl.GetGameLevelData(levelID)
	if err != nil {
		return err
	}
	gameLevel.IsNew = false
	pgl.CurrentGameLevelID = levelID
	pgl.LastRefreshTime = time.Now().Unix()
	return nil
}
