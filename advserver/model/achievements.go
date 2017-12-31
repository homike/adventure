package model

import (
	"adventure/advserver/gamedata"
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

func (pa *PlayerAchievenment) GetAchievements(condType structs.AchvCondType, condID int32) ([]*structs.Achievement,
	[]*structs.AchievementTemplate) {
	arrAchv := []*structs.Achievement{}
	arrAchvT := []*structs.AchievementTemplate{}
	for _, v := range pa.Achievements {
		achvT, ok := gamedata.AllTemplates.AchievementTemplates[v.TemplateID]
		if !ok {
			logger.Error("AchievementTemplates(%v) failed", v.TemplateID)
			continue
		}
		if achvT.ConditionType == condType && achvT.ConditionID == condID && v.Status == structs.AchvStatus_Active {
			arrAchv = append(arrAchv, v)
			arrAchvT = append(arrAchvT, &achvT)

		}
	}
	return arrAchv, arrAchvT
}
