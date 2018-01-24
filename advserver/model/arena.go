package model

import (
	"adventure/common/structs"
)

type Arena struct {
	NextRefrshTime int64                    // 下次刷新时间
	ChallengeCount int32                    // 已挑战次数
	Targets        []*structs.FightTarget   // 战斗目标
	RewardRecord   []structs.ArenaRwdStatus // 奖励领取情况
}

func NewArena() *Arena {
	return &Arena{
		NextRefrshTime: 0,
		ChallengeCount: 0,
		RewardRecord:   []structs.ArenaRwdStatus{},
		Targets:        []*structs.FightTarget{},
	}
}

func (a *Arena) GetWinCount() int32 {
	isWinCnt := int32(0)
	for _, v := range a.Targets {
		if v.IsWin {
			isWinCnt++
		}
	}
	return isWinCnt
}
