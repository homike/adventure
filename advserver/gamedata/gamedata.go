package gamedata

import (
	"adventure/advserver/log"
	"adventure/common/clog"
	"adventure/common/csv"
	"adventure/common/structs"
	"errors"
	"fmt"
)

var logger *clog.Logger

var AllTemplates *Templates

type Templates struct {
	HeroTemplates           map[int32]structs.HeroTemplate           `table:"hero"`
	GameLevelTemplates      map[int32]structs.GameLevelTemplate      `table:"GameLevel"`
	HeroLevelTemplates      map[int32]structs.HeroLevelTemplate      `table:"HeroLevel"`
	UpgradeArtifactCosts    map[int32]structs.UpgradeArtifactCost    `table:"UpgradeArtifactCost"`
	UpgradeWeaponCosts      map[int32]structs.UpgradeWeaponCost      `table:"UpgradeWeaponCost"`
	AwakeCosts              map[int32]structs.AwakeCost              `table:"AwakeCost"`
	ItemTemplates           map[int32]structs.ItemTemplate           `table:"Item"`
	UnLockBagCosts          map[int32]structs.UnLockBagCost          `table:"UnLockBagCost"`
	Battlefields            map[int32]structs.Battlefield            `table:"Battlefield"`
	RewardTemplates         map[int32]structs.RewardTemplate         `table:"Reward"`
	GameLevelEventTemplates map[int32]structs.GameLevelEventTemplate `table:"GameLevelEvent"`
	SpellTemplates          map[int32]structs.SpellTemplate          `table:"spell"`
	ArtifactTemplates       map[int32]structs.ArtifactTemplate       `table:"Artifact"`
	ResourceTemplates       map[int32]structs.ResourceTemplate       `table:"Resouce"`
	CombinationSpells       map[int32]structs.CombinationSpell       `table:"CombinationSpell"`
}

const (
	EmployReturnExp         = 50000 // 英雄.解雇返回经验需要的总经验值
	EmployReturnExpPer      = 70    // 英雄.解雇返回经验比例
	HeroAwakeMinLevel       = 30    // 英雄.英雄觉醒需要的等级
	MaxStrength             = 72000 // 角色.饱足度上限
	FightFloatValueMin      = 20    // 战斗.普通攻击浮动下限
	FightFloatValueMax      = 25    // 战斗.普通攻击浮动上限
	FightSkillFloatValueMin = 95    // 战斗.技能攻击浮动下限
	FightSkillFloatValueMax = 105   // 战斗.技能攻击浮动上限
)

func Init() {
	logger = log.GetLogger()

	AllTemplates = &Templates{
		HeroTemplates:           make(map[int32]structs.HeroTemplate),
		GameLevelTemplates:      make(map[int32]structs.GameLevelTemplate),
		HeroLevelTemplates:      make(map[int32]structs.HeroLevelTemplate),
		UpgradeArtifactCosts:    make(map[int32]structs.UpgradeArtifactCost),
		UpgradeWeaponCosts:      make(map[int32]structs.UpgradeWeaponCost),
		AwakeCosts:              make(map[int32]structs.AwakeCost),
		ItemTemplates:           make(map[int32]structs.ItemTemplate),
		UnLockBagCosts:          make(map[int32]structs.UnLockBagCost),
		Battlefields:            make(map[int32]structs.Battlefield),
		RewardTemplates:         make(map[int32]structs.RewardTemplate),
		GameLevelEventTemplates: make(map[int32]structs.GameLevelEventTemplate),
		SpellTemplates:          make(map[int32]structs.SpellTemplate),
		ArtifactTemplates:       make(map[int32]structs.ArtifactTemplate),
		ResourceTemplates:       make(map[int32]structs.ResourceTemplate),
		CombinationSpells:       make(map[int32]structs.CombinationSpell),
	}
	csv.LoadTemplates2(AllTemplates)

	fmt.Println("Skill ", AllTemplates.HeroTemplates[16205].SkillID[1])
	fmt.Println("battle Name ", AllTemplates.Battlefields[1001].Name)
}

func GetHeroLevelExp(heroLv int32, awakeCnt int32) (int32, error) {
	exp := int(0)
	err := errors.New("GetHeroLevelExp failed")

	template, ok := AllTemplates.HeroLevelTemplates[heroLv]
	if !ok {
		return 0, err
	}
	switch awakeCnt {
	case 0:
		return template.EXP1, nil
	case 1:
		return template.EXP2, nil
	case 2:
		return template.EXP3, nil
	case 3:
		return template.EXP4, nil
	case 4:
		return template.EXP5, nil
	}
	return int32(exp), nil
}
