package model

import "Adventure/common/structs"

type HeroTeams struct {
	Heros                map[int32]*structs.Hero      // 英雄列表
	MaxWorker            int32                        // 出战人数上限
	EmployRecord         map[structs.EmployType]int32 // 招募英雄记录表
	EmployTotalCount     int32                        // 招募英雄的总数量
	EmployCount          int32                        // 招募英雄的总数量(周期的活动值，可刷新)
	IngotEmployCount     int32                        // 元宝招募次数
	ManyIngotEmployCount int32                        // 十连抽的招募次数
	ManyIngotEmployRP    int32                        // 10连抽的人品值
	IngotEmployRP        int32                        // 元宝的人品值
	MaxHeroId            int32                        // 英雄的最大id
	LastEmployTime       uint64                       // 最近一次刷新的招募时间
}

func NewHeroTeams() *HeroTeams {
	teams := &HeroTeams{
		Heros:        make(map[int32]*structs.Hero),
		EmployRecord: make(map[structs.EmployType]int32),
	}
	teams.EmployRecord[structs.Money] = 0
	teams.EmployRecord[structs.HunLuan] = 0
	teams.EmployRecord[structs.HuiHuang] = 0
	teams.EmployRecord[structs.LvDong] = 0
	teams.EmployRecord[structs.Diamond] = 0
	teams.EmployRecord[structs.ManyDiamond] = 0

	return teams
}

// 团队的hp(当前战力)
func (h *HeroTeams) MaxHP() int32 {
	totalHP := int32(0)
	for _, v := range h.Heros {
		if v.IsOutFight {
			totalHP += v.HP
		}
	}
	return totalHP
}

func (h *HeroTeams) AddHero(hero *structs.Hero) error {
	return nil
}
