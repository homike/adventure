package model

import (
	"adventure/advserver/gamedata"
	"time"
)

type ExtendData struct {
	BattleRecords map[int32]int64
}

func NewExtendData() *ExtendData {
	extend := &ExtendData{
		BattleRecords: make(map[int32]int64),
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
