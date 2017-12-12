package gamedata

import (
	"adventure/advserver/log"
	"adventure/common/clog"
	"adventure/common/csv"
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
		Coefficient         csv.Float64 `table:"hero" key:"" val:"成长系数"`
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
}

const (
	EmployReturnExp    = 50000 // 英雄.解雇返回经验需要的总经验值
	EmployReturnExpPer = 70    // 英雄.解雇返回经验比例
)

func Init() {
	logger = log.GetLogger()

	AllTemplates = Templates{}

	csv.LoadTemplates(&AllTemplates)
}
