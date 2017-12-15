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

func (pgl *PlayerGameLevel) GetCurGameLevel() (*structs.GameLevel, error) {
	for _, v := range pgl.GameLevels {
		if v.GameLevelID == pgl.CurrentGameLevelID {
			return v, nil
		}
	}
	return nil, errors.New("GetCurGameLevel() Error")
}
