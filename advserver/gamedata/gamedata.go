package gamedata

import (
	"adventure/advserver/log"
	"adventure/common/clog"
	"adventure/common/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var logger *clog.Logger

var AllTemplates Templates

type Templates struct {
	HeroTemplate struct { // 英雄模板
		HeroName            csv.String  `table:"hero" key:"" val:"名字"`
		SkillID             csv.Int     `table:"hero" key:"" val:"技能ID列表"`
		CombinationSpllID   csv.Int     `table:"hero" key:"" val:"组合技能ID"`
		IconID              csv.Int     `table:"hero" key:"" val:"模型ID"`
		QualityType         csv.Int     `table:"hero" key:"" val:"品质"`
		Profession          csv.String  `table:"hero" key:"" val:"职业"`
		Ethnicity           csv.String  `table:"hero" key:"" val:"种族"`
		Sex                 csv.String  `table:"hero" key:"" val:"性别"`
		Description         csv.String  `table:"hero" key:"" val:"描述"`
		BaseHP              csv.Int     `table:"hero" key:"" val:"基础战力"`
		Coefficient         csv.Float32 `table:"hero" key:"" val:"成长系数"`
		HonorDebris         csv.Int     `table:"hero" key:"" val:"荣誉碎片"`
		AwakeCount          csv.Int     `table:"hero" key:"" val:"觉醒次数上限"`
		TempleAppearWeight  csv.Int     `table:"hero" key:"" val:"神殿出现概率权重"`
		EmployCostFragments csv.Int     `table:"hero" key:"" val:"神殿兑换该英雄的花费"`
		EmployWeight        csv.String  `table:"hero" key:"" val:"招募权重"`
		// 招募权重，金币\混乱之门 301 \辉煌之门 302 \律动之门 303\万象之门（元宝）\传奇之门普通（元宝）\传奇之门特殊（元宝）
	}
	GameLevelTemplate struct { // 游戏关卡
		Title              csv.Int     `table:"GameLevel" key:"" val:"关卡"`
		EvnetIDs           csv.String  `table:"GameLevel" key:"" val:"事件ID队列"`
		GameBoxIDs         csv.String  `table:"GameLevel" key:"" val:"宝箱奖励ID队列"`
		GameBoxWeight      csv.String  `table:"GameLevel" key:"" val:"宝箱奖励权重队列"`
		SpeedGameBoxIDs    csv.String  `table:"GameLevel" key:"" val:"加速宝箱奖励ID队列"`
		SpeedGameBoxWeight csv.String  `table:"GameLevel" key:"" val:"加速宝箱奖励权重队列"`
		ActiveGameBoxSec   csv.Int     `table:"GameLevel" key:"" val:"宝箱需要时间"`
		MoneyPer           csv.Int     `table:"GameLevel" key:"" val:"每秒产出金币"`
		ExpPer             csv.Int     `table:"GameLevel" key:"" val:"每秒产出经验"`
		MinHP              csv.Int     `table:"GameLevel" key:"" val:"需要战力"`
		IconID             csv.Float64 `table:"GameLevel" key:"" val:"图标ID"`
	}
	HeroLevelTemplate struct { // 英雄等级模板数据
		EXP1 csv.Int `table:"HeroLevel" key:"" val:"经验1"`
		EXP2 csv.Int `table:"HeroLevel" key:"" val:"经验2"`
		EXP3 csv.Int `table:"HeroLevel" key:"" val:"经验3"`
		EXP4 csv.Int `table:"HeroLevel" key:"" val:"经验4"`
		EXP5 csv.Int `table:"HeroLevel" key:"" val:"经验5"`
	}
	UpgradeArtifactCost struct { // 神器升级消耗表
		//Level                 csv.Int    `table:"HeroLevel" key:"" val:"神器等级"`
		NeedResourceIdList    csv.String `table:"UpgradeArtifactCost" key:"" val:"需要消耗的资源id列表"`
		NeedResourceCountList csv.String `table:"UpgradeArtifactCost" key:"" val:"需要消耗的资源数量列表"`
		WeaponParam           csv.Int    `table:"UpgradeArtifactCost" key:"" val:"神器系数"`
	}
	UpgradeWeaponCost struct { // 武具级消耗表
		//Level                 csv.Int    `table:"HeroLevel" key:"" val:"武具等级"`
		NeedResourceIdList    csv.String `table:"HeroLevel" key:"" val:"需要消耗的资源id列表"`
		NeedResourceCountList csv.String `table:"HeroLevel" key:"" val:"需要消耗的资源数量列表"`
		WeaponParam           csv.Int    `table:"HeroLevel" key:"" val:"武具系数"`
		SuccessRate           csv.Int    `table:"HeroLevel" key:"" val:"成功率"`
		AddForgePoint         csv.Int    `table:"HeroLevel" key:"" val:"增加的锻造点"`
	}
	AwakeCost struct { // 觉醒消耗表
		//Level                 csv.Int    `table:"HeroLevel" key:"" val:"等级"`
		Money  csv.Int     `table:"AwakeCost" key:"" val:"金币"`
		Statue csv.Int     `table:"AwakeCost" key:"" val:"雕像"`
		Param  csv.Float32 `table:"AwakeCost" key:"" val:"觉醒系数"`
	}
	ItemTemplate struct { // 物品模板
		//ID                 csv.Int    `table:"HeroLevel" key:"" val:"ID"`
		Name              csv.String `table:"Item" key:"" val:"名字"`
		Type              csv.Int    `table:"Item" key:"" val:"类型"`
		ExType            csv.Int    `table:"Item" key:"" val:"扩展类型"`
		SellMoney         csv.Int    `table:"Item" key:"" val:"出售价格"`
		IconID            csv.String `table:"Item" key:"" val:"图标"`
		Description       csv.String `table:"Item" key:"" val:"描述"`
		ShopID            csv.Int    `table:"Item" key:"" val:"商城id"`
		ShopPrice         csv.Int    `table:"Item" key:"" val:"商城价格"`
		Param1            csv.Int    `table:"Item" key:"" val:"参数1"`
		Param2            csv.Int    `table:"Item" key:"" val:"参数2"`
		Param3            csv.Int    `table:"Item" key:"" val:"参数3"`
		RewardIDs         csv.String `table:"Item" key:"" val:"奖励id"`
		OccurWeight       csv.String `table:"Item" key:"" val:"权重"`
		IsOnceEveryday    csv.Bool   `table:"Item" key:"" val:"是否每天只能使用一次"`
		RewardDescription csv.String `table:"Item" key:"" val:"奖励描述"`
	}
	UnLockBagCost struct { // 背包格子解锁表
		//UnlockLevel csv.String `table:"UnLockBagCost" key:"" val:"开启次数"`
		BagCount    csv.Int    `table:"UnLockBagCost" key:"" val:"开启格子数"`
		CostResIDs  csv.String `table:"UnLockBagCost" key:"" val:"资源类型"`
		CostResNums csv.String `table:"UnLockBagCost" key:"" val:"资源数量"`
		UnLockCount csv.Counts `table:"UnLockBagCost" key:"" val:""`
	}
}

const (
	EmployReturnExp    = 50000 // 英雄.解雇返回经验需要的总经验值
	EmployReturnExpPer = 70    // 英雄.解雇返回经验比例
	HeroAwakeMinLevel  = 30    // 英雄.英雄觉醒需要的等级
)

func Init() {
	logger = log.GetLogger()

	AllTemplates = Templates{}

	csv.LoadTemplates(&AllTemplates)
}

func GetHeroLevelExp(heroLv int32, awakeCnt int32) (int32, error) {
	exp := int(0)
	err := errors.New("GetHeroLevelExp failed")
	switch awakeCnt {
	case 0:
		exp, err = AllTemplates.HeroLevelTemplate.EXP1(heroLv)
		if err != nil {
			return 0, err
		}
	case 1:
		exp, err = AllTemplates.HeroLevelTemplate.EXP2(heroLv)
		if err != nil {
			return 0, err
		}
	case 2:
		exp, err = AllTemplates.HeroLevelTemplate.EXP3(heroLv)
		if err != nil {
			return 0, err
		}
	case 3:
		exp, err = AllTemplates.HeroLevelTemplate.EXP4(heroLv)
		if err != nil {
			return 0, err
		}
	case 4:
		exp, err = AllTemplates.HeroLevelTemplate.EXP5(heroLv)
		if err != nil {
			return 0, err
		}
	}
	return int32(exp), nil
}

func GetUnlockBagCost(unlockLv int32) ([]int32, []int32, error) {
	strResIDs, err := AllTemplates.UnLockBagCost.CostResIDs(unlockLv)
	if err != nil {
		return nil, nil, err
	}
	resIDs := []int32{}
	err = json.Unmarshal([]byte(strResIDs), resIDs)
	if err != nil {
		return nil, nil, err
	}

	strResNums, err := AllTemplates.UnLockBagCost.CostResNums(unlockLv)
	if err != nil {
		return nil, nil, err
	}

	resNums := []int32{}
	err = json.Unmarshal([]byte(strResNums), resNums)
	if err != nil {
		return nil, nil, err
	}

	return resIDs, resNums, nil
}

func GetGameLevelEvents(gameLv int32) ([]int32, error) {
	events := []int32{}
	strEvents, err := AllTemplates.GameLevelTemplate.EvnetIDs(gameLv)
	if err != nil {
		fmt.Println("GetGameLevelEvents() ", err)
		return events, err
	}

	arrEvents := strings.Split(strEvents, ";")
	if len(arrEvents) <= 0 {
		fmt.Println("GetGameLevelEvents() ", err, " strEvents : ", strEvents)
		return events, err
	}

	events = make([]int32, 0, len(arrEvents))
	for _, v := range arrEvents {
		n, _ := strconv.Atoi(v)
		events = append(events, int32(n))
	}

	return events, nil
}
