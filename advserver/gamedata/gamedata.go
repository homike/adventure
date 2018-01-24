package gamedata

import (
	"adventure/advserver/log"
	"adventure/common/clog"
	"adventure/common/csv"
	"adventure/common/structs"
	"adventure/common/util"
	"errors"
	"fmt"
	"math"
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
	AchievementTemplates    map[int32]structs.AchievementTemplate    `table:"achievementAward"`
	GlobalData              structs.GlobalTemplate                   `table:"GlobalData"`
}

var ArenaRobtosGroup []*structs.PlayerGroup
var ArenaRobotsMap map[int32]*structs.PlayerBaseInfo

const (
	EmployReturnExp              = 50000        // 英雄.解雇返回经验需要的总经验值
	EmployReturnExpPer           = 70           // 英雄.解雇返回经验比例
	HeroAwakeMinLevel            = 30           // 英雄.英雄觉醒需要的等级
	MaxStrength                  = 72000        // 角色.饱足度上限
	FightFloatValueMin           = 20           // 战斗.普通攻击浮动下限
	FightFloatValueMax           = 25           // 战斗.普通攻击浮动上限
	FightSkillFloatValueMin      = 95           // 战斗.技能攻击浮动下限
	FightSkillFloatValueMax      = 105          // 战斗.技能攻击浮动上限
	MaxHeroWorkTopItem           = 20           // 英雄.英雄最大出战数
	FreeIngotEmployFirstTimeSpan = 180          // 英雄.首次免费元宝招募间隔
	FreeIngotEmployTimeSpanItem  = 60 * 60 * 48 // 英雄.免费元宝招募间隔
	SystemRefreshTime            = 10           // 英雄.系统每天的刷新时间
)

var (
	InitialEmployCost = []int32{3000, 1, 1, 1, 60, 580} // 英雄.招募英雄初始支出
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
		AchievementTemplates:    make(map[int32]structs.AchievementTemplate),
	}
	csv.LoadTemplates2(AllTemplates)

	LoadArenaRobots()

	fmt.Println("EmployReturnExp ", AllTemplates.GlobalData.EmployReturnExp)
	// fmt.Println("battle Name ", AllTemplates.Battlefields[1001].Name)
}

func GetHeroLevelExp(heroLv int32, awakeCnt int32) (int32, error) {
	exp := int(0)
	err := errors.New("GetHeroLevelExp failed")

	template, ok := AllTemplates.HeroLevelTemplates[heroLv]
	if !ok {
		return 0, err
	}
	switch awakeCnt {
	case 1:
		return template.EXP1, nil
	case 2:
		return template.EXP2, nil
	case 3:
		return template.EXP3, nil
	case 4:
		return template.EXP4, nil
	case 5:
		return template.EXP5, nil
	}
	return int32(exp), nil
}

func LoadArenaRobots() {
	ArenaRobtosGroup = make([]*structs.PlayerGroup, 0)
	ArenaRobotsMap = make(map[int32]*structs.PlayerBaseInfo)

	id, robotID, robotNum := int32(0), int32(0), 128
	minHP, maxHP := int32(0), int32(0)
	for i := 0; i < len(AllTemplates.GlobalData.HPRange); i++ {
		id++
		maxHP = AllTemplates.GlobalData.HPRange[i]

		group := structs.PlayerGroup{
			ID:      id,
			MinHP:   minHP,
			MaxHP:   maxHP,
			Players: make([]*structs.PlayerBaseInfo, robotNum),
		}
		for j := 0; j < robotNum; j++ {
			robotID++
			robot := CreateRobot(robotID)
			robot.HP = util.RandNum(minHP, maxHP-minHP)
			group.Players[j] = robot

			ArenaRobotsMap[robot.ID] = robot
		}

		ArenaRobtosGroup = append(ArenaRobtosGroup, &group)
		minHP = maxHP
	}

	id++
	ArenaRobtosGroup = append(ArenaRobtosGroup, &structs.PlayerGroup{
		ID:    id,
		MinHP: maxHP,
		MaxHP: math.MaxInt32,
	})
}

func CreateRobot(id int32) *structs.PlayerBaseInfo {
	robot := structs.PlayerBaseInfo{
		ID:   id,
		Name: AllTemplates.GlobalData.RobotName[util.RandNum(int32(0), int32(len(AllTemplates.GlobalData.RobotName)))],
	}

	heros := []*structs.HeroPostion{}
	for i := 0; i < 20; i++ {
		templateID := AllTemplates.GlobalData.RobotHeroIDs[util.RandNum(int32(0), int32(len(AllTemplates.GlobalData.RobotHeroIDs)))]
		template, ok := AllTemplates.HeroTemplates[templateID]
		if ok {
			heros = append(heros, &structs.HeroPostion{
				HeroTemplateID: templateID,
				HeroIconID:     int32(template.IconID),
				HeroPosition:   int32(i),
			})
		}
	}

	robot.HeroIDs = heros
	robot.IconID = heros[0].HeroIconID

	return &robot
}
