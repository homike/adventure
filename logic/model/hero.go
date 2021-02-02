package model

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"errors"
	"sort"
)

type HeroTeams struct {
	Heros                []*structs.Hero              // 英雄列表
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
		MaxHeroId:    1,
		MaxWorker:    3,
		Heros:        make([]*structs.Hero, 0),
		EmployRecord: make(map[structs.EmployType]int32),
	}
	teams.EmployRecord[structs.EmployType_Money] = 0
	teams.EmployRecord[structs.EmployType_HunLuan] = 0
	teams.EmployRecord[structs.EmployType_HuiHuang] = 0
	teams.EmployRecord[structs.EmployType_LvDong] = 0
	teams.EmployRecord[structs.EmployType_Diamond] = 0
	teams.EmployRecord[structs.EmployType_ManyDiamond] = 0

	return teams
}

// 团队的hp(当前战力)
func (h *HeroTeams) MaxHP() int32 {
	totalHP := int32(0)
	for _, v := range h.Heros {
		if v.IsOutFight {
			totalHP += v.HP()
		}
	}
	return totalHP
}

func (h *HeroTeams) AddHero(name string, isPlayer bool, heroTemplateID int32) (*structs.Hero, error) {
	heroT, ok := gamedata.AllTemplates.HeroTemplates[heroTemplateID]
	if !ok {
		return nil, errors.New("HeroTemplates Error")
	}
	configName, configQualityType, configBaseHP := heroT.HeroName, heroT.QualityType, heroT.BaseHP

	isOutFight := false
	heroName := configName
	quality := structs.QualityType(configQualityType)
	if isPlayer {
		isOutFight = true
		heroName = name
		quality = structs.QualityType_Gold
	}

	hero := &structs.Hero{
		HeroID:         h.MaxHeroId,
		IsOutFight:     isOutFight,
		IsPlayer:       isPlayer,
		Level:          1,
		HeroTemplateID: heroTemplateID,
		Name:           heroName,
		Quality:        quality,
		AwakeCount:     1,
		LevelHP:        int32(configBaseHP),
		Index:          int32(len(h.Heros)),
	}
	h.Heros = append(h.Heros, hero)
	h.MaxHeroId++

	return hero, nil
}

func (h *HeroTeams) RemoveHero(hero *structs.Hero) error {
	for k, v := range h.Heros {
		if v.HeroID == hero.HeroID {
			h.Heros = append(h.Heros[0:k], h.Heros[k+1:]...)
			return nil
		}
	}

	return errors.New("RemoveHero error")
}

func (h *HeroTeams) GetHerosArray() []structs.Hero {
	heros := make([]structs.Hero, 0, len(h.Heros))
	for _, v := range h.Heros {
		heros = append(heros, *v)
	}
	return heros
}

func (h *HeroTeams) GetHero(heroID int32) (*structs.Hero, error) {
	for _, v := range h.Heros {
		if v.HeroID == heroID {
			return v, nil
		}
	}
	return nil, errors.New("has not this hero")
}

func (h *HeroTeams) GetMainHero() (*structs.Hero, error) {
	for _, v := range h.Heros {
		if v.IsPlayer {
			return v, nil
		}
	}
	return nil, errors.New("has not main hero")
}

func (h *HeroTeams) GetFightHeros() []*structs.Hero {
	heros := []*structs.Hero{}
	for _, v := range h.Heros {
		if v.IsOutFight {
			heros = append(heros, v)
		}
	}
	return heros
}

func (h *HeroTeams) GetPlayerHero() *structs.Hero {
	for _, v := range h.Heros {
		if v.IsPlayer {
			return v
		}
	}
	return nil
}

type HeroByIndex []*structs.Hero

func (a HeroByIndex) Len() int           { return len(a) }
func (a HeroByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a HeroByIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

func (h *HeroTeams) SortHeros() error {
	sort.Sort(HeroByIndex(h.Heros))
	return nil
}

func (h *HeroTeams) ReCalculateHeroLevelHp(hero *structs.Hero) error {
	heroT, ok := gamedata.AllTemplates.HeroTemplates[hero.HeroTemplateID]
	if !ok {
		return errors.New("HeroTemplates Error")
	}
	baseHP, coefficient := heroT.BaseHP, heroT.Coefficient

	artifactCostT, ok := gamedata.AllTemplates.UpgradeArtifactCosts[hero.WeaponLevel]
	if !ok {
		return errors.New("UpgradeArtifactCosts Error")
	}
	weaponCostT, ok := gamedata.AllTemplates.UpgradeWeaponCosts[hero.WeaponLevel]
	if !ok {
		return errors.New("UpgradeWeaponCosts Error")
	}

	weaponParam := float32(1)
	if hero.WeaponLevel > 0 {
		param := int(0)
		if hero.IsPlayer {
			param = artifactCostT.WeaponParam
		} else {
			param = weaponCostT.WeaponParam
		}
		weaponParam = float32(param / 100.0)
	}

	jieXianCount := float32(1) //暂时使用，代替突破界限系数
	awakeCostT, ok := gamedata.AllTemplates.AwakeCosts[hero.AwakeCount]
	if ok {
		jieXianCount = float32(awakeCostT.Param)
	}

	//计算英雄战力
	param1 := float32(baseHP) + float32(hero.Level)*coefficient
	param2 := param1*jieXianCount + float32(hero.ItemHP)
	hp := int32(param2 * weaponParam)
	hero.LevelHP = hp

	return nil
	//PlayerEvents.OnHPChange(player)
}
