package gamedata

import (
	"adventure/advserver/log"
	"adventure/common/clog"
	"adventure/common/csv"
	"adventure/common/structs"
	"fmt"
)

var logger *clog.Logger

var AllTemplates *Templates

type Templates struct {
	HeroTemplates        map[int32]structs.HeroTemplate        `table:"hero"`
	GameLevelTemplates   map[int32]structs.GameLevelTemplate   `table:"GameLevel"`
	HeroLevelTemplates   map[int32]structs.HeroLevelTemplate   `table:"HeroLevel"`
	UpgradeArtifactCosts map[int32]structs.UpgradeArtifactCost `table:"UpgradeArtifactCost"`
	UpgradeWeaponCosts   map[int32]structs.UpgradeWeaponCost   `table:"UpgradeWeaponCost"`
	AwakeCosts           map[int32]structs.AwakeCost           `table:"AwakeCost"`
	ItemTemplates        map[int32]structs.ItemTemplate        `table:"Item"`
	UnLockBagCosts       map[int32]structs.UnLockBagCost       `table:"UnLockBagCost"`
	Battlefields         map[int32]structs.Battlefield         `table:"Battlefield"`
}

const (
	EmployReturnExp    = 50000 // 英雄.解雇返回经验需要的总经验值
	EmployReturnExpPer = 70    // 英雄.解雇返回经验比例
	HeroAwakeMinLevel  = 30    // 英雄.英雄觉醒需要的等级
)

func Init() {
	logger = log.GetLogger()

	AllTemplates = &Templates{
		HeroTemplates:        make(map[int32]structs.HeroTemplate),
		GameLevelTemplates:   make(map[int32]structs.GameLevelTemplate),
		HeroLevelTemplates:   make(map[int32]structs.HeroLevelTemplate),
		UpgradeArtifactCosts: make(map[int32]structs.UpgradeArtifactCost),
		UpgradeWeaponCosts:   make(map[int32]structs.UpgradeWeaponCost),
		AwakeCosts:           make(map[int32]structs.AwakeCost),
		ItemTemplates:        make(map[int32]structs.ItemTemplate),
		UnLockBagCosts:       make(map[int32]structs.UnLockBagCost),
		Battlefields:         make(map[int32]structs.Battlefield),
	}
	csv.LoadTemplates2(AllTemplates)

	fmt.Println("Skill ", AllTemplates.HeroTemplates[16205].SkillID[1])
	fmt.Println("battle Name ", AllTemplates.Battlefields[1001].Name)
}

func GetHeroLevelExp(heroLv int32, awakeCnt int32) (int32, error) {
	//exp := int(0)
	//err := errors.New("GetHeroLevelExp failed")
	// switch awakeCnt {
	// case 0:
	// 	exp, err = AllTemplates.HeroLevelTemplate.EXP1(heroLv)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// case 1:
	// 	exp, err = AllTemplates.HeroLevelTemplate.EXP2(heroLv)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// case 2:
	// 	exp, err = AllTemplates.HeroLevelTemplate.EXP3(heroLv)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// case 3:
	// 	exp, err = AllTemplates.HeroLevelTemplate.EXP4(heroLv)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// case 4:
	// 	exp, err = AllTemplates.HeroLevelTemplate.EXP5(heroLv)
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// }
	//return int32(exp), nil

	return int32(0), nil
}

func GetUnlockBagCost(unlockLv int32) ([]int32, []int32, error) {
	// strResIDs, err := AllTemplates.UnLockBagCost.CostResIDs(unlockLv)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// resIDs := []int32{}
	// err = json.Unmarshal([]byte(strResIDs), resIDs)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// strResNums, err := AllTemplates.UnLockBagCost.CostResNums(unlockLv)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// resNums := []int32{}
	// err = json.Unmarshal([]byte(strResNums), resNums)
	// if err != nil {
	// 	return nil, nil, err
	// }

	// return resIDs, resNums, nil
	return nil, nil, nil
}

func GetGameLevelEvents(gameLv int32) ([]int32, error) {
	// events := []int32{}
	// strEvents, err := AllTemplates.GameLevelTemplate.EvnetIDs(gameLv)
	// if err != nil {
	// 	fmt.Println("GetGameLevelEvents() ", err)
	// 	return events, err
	// }

	// return nil, GetIntSliceByString(strEvents)
	return nil, nil
}

func GetGameLevelGameBox(gameLv int32) ([]int, []int, error) {
	// boxIDs := []int32{}
	// weights := []int32{}
	// strBoxIDs, err := AllTemplates.GameLevelTemplate.GameBoxIDs(gameLv)
	// if err != nil {
	// 	fmt.Println("GameBoxIDs() ", err)
	// 	return nil, boxIDs, err
	// }

	// strWeights, err := AllTemplates.GameLevelTemplate.GameBoxWeight(gameLv)
	// if err != nil {
	// 	fmt.Println("GameBoxIDs() ", err)
	// 	return nil, boxIDs, err
	// }

	// return GetIntSliceByString(strBoxIDs), GetIntSliceByString(strWeights), nil
	return nil, nil, nil
}

func GetIntSliceByString(strSlice string) []int32 {
	// sliceInt := []int32{}

	// sliceStr := strings.Split(strSlice, ";")
	// if len(sliceStr) <= 0 {
	// 	fmt.Println("GetIntSliceByString() error")
	// 	return sliceInt
	// }

	// sliceInt = make([]int32, 0, len(sliceStr))
	// for _, v := range sliceStr {
	// 	n, _ := strconv.Atoi(v)
	// 	sliceInt = append(sliceInt, int32(n))
	// }

	// return sliceInt
	return nil
}
