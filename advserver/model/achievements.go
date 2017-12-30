package model

import (
	"adventure/common/structs"
)

type PlayerAchievenment struct {
	Achievements          map[int32]*structs.Achievement
	NextRefreshTimeDaily  int64 // 下一次刷新每日成就的时间
	NextRefreshTimeWeekly int64 // 下一次刷新每周成就的时间
}

func NewPlayerAchievenment() *PlayerAchievenment {
	return &PlayerAchievenment{
		Achievements:          make(map[int32]*structs.Achievement),
		NextRefreshTimeDaily:  0,
		NextRefreshTimeWeekly: 0,
	}
}

func (pa *PlayerAchievenment) GetAchieveMentsArray() []*structs.Achievement {
	arrAchv := make([]*structs.Achievement, len(pa.Achievements))
	for _, v := range pa.Achievements {
		arrAchv = append(arrAchv, v)
	}
	return arrAchv
}
