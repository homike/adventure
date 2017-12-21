package structs

type HeroTemplate struct { // 英雄模板
	HeroName            string  `val:"名字"`
	SkillID             []int32 `val:"技能ID列表"`
	CombinationSpllID   []int32 `val:"组合技能ID"`
	IconID              int     `val:"模型ID"`
	QualityType         int     `val:"品质"`
	Profession          string  `val:"职业"`
	Ethnicity           string  `val:"种族"`
	Sex                 string  `val:"性别"`
	Description         string  `val:"描述"`
	BaseHP              int     `val:"基础战力"`
	Coefficient         float32 `val:"成长系数"`
	HonorDebris         int     `val:"荣誉碎片"`
	AwakeCount          int     `val:"觉醒次数上限"`
	TempleAppearWeight  int     `val:"神殿出现概率权重"`
	EmployCostFragments int     `val:"神殿兑换该英雄的花费"`
	EmployWeight        string  `val:"招募权重"`
	// 招募权重，金币\混乱之门 301 \辉煌之门 302 \律动之门 303\万象之门（元宝）\传奇之门普通（元宝）\传奇之门特殊（元宝）
}

type GameLevelTemplate struct { // 游戏关卡
	Title              string  `val:"关卡名称"`
	EvnetIDs           []int32 `val:"事件ID队列"`
	GameBoxIDs         []int32 `val:"宝箱奖励ID队列"`
	GameBoxWeight      []int32 `val:"宝箱奖励权重队列"`
	SpeedGameBoxIDs    []int32 `val:"加速宝箱奖励ID队列"`
	SpeedGameBoxWeight []int32 `val:"加速宝箱奖励权重队列"`
	ActiveGameBoxSec   int32   `val:"宝箱需要时间"`
	MoneyPer           int32   `val:"每秒产出金币"`
	ExpPer             int32   `val:"每秒产出经验"`
	MinHP              int32   `val:"需要战力"`
	IconID             string  `val:"图标ID"`
}
type HeroLevelTemplate struct { // 英雄等级模板数据
	EXP1 int `val:"经验1"`
	EXP2 int `val:"经验2"`
	EXP3 int `val:"经验3"`
	EXP4 int `val:"经验4"`
	EXP5 int `val:"经验5"`
}

type UpgradeArtifactCost struct { // 神器升级消耗表
	//Level                 int    `val:"神器等级"`
	NeedResourceIdList    string `val:"需要消耗的资源id列表"`
	NeedResourceCountList string `val:"需要消耗的资源数量列表"`
	WeaponParam           int    `val:"神器系数"`
}

type UpgradeWeaponCost struct { // 武具级消耗表
	//Level                 int    `val:"武具等级"`
	NeedResourceIdList    string `val:"需要消耗的资源id列表"`
	NeedResourceCountList string `val:"需要消耗的资源数量列表"`
	WeaponParam           int    `val:"武具系数"`
	SuccessRate           int    `val:"成功率"`
	AddForgePoint         int    `val:"增加的锻造点"`
}

type AwakeCost struct { // 觉醒消耗表
	//Level                 int    `val:"等级"`
	Money  int     `val:"金币"`
	Statue int     `val:"雕像"`
	Param  float32 `val:"觉醒系数"`
}

type ItemTemplate struct { // 物品模板
	//ID                 int    `val:"ID"`
	Name   string `val:"名字"`
	Type   int    `val:"类型"`
	ExType int    `val:"扩展类型"`
	//SellMoney         int    `val:"出售价格"`
	IconID      string `val:"图标"`
	Description string `val:"描述"`
	//ShopID            int    `val:"商城id"`
	//ShopPrice         int    `val:"商城价格"`
	//Param1            int    `val:"参数1"`
	//Param2            int    `val:"参数2"`
	//Param3            int    `val:"参数3"`
	RewardIDs         []int32 `val:"奖励id"`
	OccurWeight       []int32 `val:"权重"`
	IsOnceEveryday    bool    `val:"是否每天只能使用一次"`
	RewardDescription string  `val:"奖励描述"`
}

type UnLockBagCost struct { // 背包格子解锁表
	//UnlockLevel string `val:"开启次数"`
	BagCount    int    `val:"开启格子数"`
	CostResIDs  string `val:"资源类型"`
	CostResNums string `val:"资源数量"`
}

// Battlefield
type Battlefield struct {
	Name             string  `val:"名字"`
	HP               int     `val:"战力"`
	NpcIDs           []int32 `val:"英雄ID列表"`
	XianGong         int     `val:"先攻"`
	FangYu           int     `val:"防御"`
	ShanBi           int     `val:"闪避"`
	WangZhe          int     `val:"王者"`
	BackgroundID     string  `val:"背景ID"`
	ForegroundID     string  `val:"前景ID"`
	RewardType       int     `val:"类型"` // 0: 固定, 1: 随机
	RewardIDs        []int32 `val:"奖励Id列表"`
	Weights          []int32 `val:"奖励权重列表"`
	LostWarDelayTime int     `val:"冷却时间"`
	CostFood         int     `val:"消耗的饱足度"`
}
