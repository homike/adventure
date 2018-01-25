package model

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"time"
)

type Temple struct {
	TempleHeros     []*structs.TempleHero // 神殿英雄
	NextRefreshTime int64                 // 最近一次的刷新时间
	ToDayTradeCount int32                 // 今天已经兑换的次数
	RefreshCount    int32                 // 今天已经刷新过的次数
}

func NewTemple() *Temple {
	return &Temple{
		TempleHeros:     make([]*structs.TempleHero, 0),
		NextRefreshTime: 0,
		ToDayTradeCount: 0,
		RefreshCount:    0,
	}
}

func (t *Temple) RefreshHeros(useIngot bool) {
	heros := []*structs.TempleHero{}

	for _, v := range gamedata.AllTemplates.GlobalData.TempleHeroQualityQueue {
		num := 1
		if v == int32(structs.QualityType_Purple) {
			num = 2
		} else if v == int32(structs.QualityType_Blue) {
			num = 4
		}

		for i := 0; i < num; i++ {
			heroT := gamedata.RandomTempleHero(structs.QualityType(v), useIngot)
			bRepeat := false
			for _, v := range heros {
				if v.HeroTemplateID == heroT.ID {
					bRepeat = true
				}
			}
			if bRepeat {
				i--
			}
			heros = append(heros, &structs.TempleHero{Quality: heroT.QualityType, HeroTemplateID: heroT.ID, Cost: heroT.EmployCostFragments})
		}
	}

	t.TempleHeros = heros
}

func (t *Temple) RefreshTemple() {
	if t.NextRefreshTime < time.Now().Unix() {
		t.RefreshHeros(false)

		t.NextRefreshTime = time.Now().AddDate(0, 0, 1).Unix()
		t.ToDayTradeCount = 0
		t.RefreshCount = 0
	}
}
