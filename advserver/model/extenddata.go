package model

import (
	"adventure/advserver/gamedata"
	"time"
)

type ExtendData struct {
	BattleRecords    map[int32]int64 // 战斗记录
	EatedFoodRecords map[int32]int64 // 吃过的食物列表
}

func NewExtendData() *ExtendData {
	extend := &ExtendData{
		BattleRecords:    make(map[int32]int64),
		EatedFoodRecords: make(map[int32]int64),
	}
	return extend
}

func (e *ExtendData) GetGamelevelLeftTime(battleID int32) int32 {
	leftTime := int32(0)
	for k, v := range e.BattleRecords {
		if k == battleID {
			lTime := time.Unix(v, 0)
			lostTime := int(time.Now().Sub(lTime).Seconds())

			battleFiledsT, ok := gamedata.AllTemplates.Battlefields[battleID]
			if !ok {
				logger.Debug("Battlefields(%v)", battleID)
				return leftTime
			}
			delayTime := battleFiledsT.LostWarDelayTime

			leftTime := delayTime - int32(lostTime)
			if leftTime < 0 {
				leftTime = 0
			}
			break
		}
	}
	return leftTime
}

func (e *ExtendData) GetEatedFoodRecord(foodID int32) int64 {
	v, ok := e.EatedFoodRecords[foodID]
	if !ok {
		return 0
	}
	return v
}

func (e *ExtendData) AddEatedFoodRecord(foodID int32) {
	e.EatedFoodRecords[foodID] = time.Now().Unix()
}
