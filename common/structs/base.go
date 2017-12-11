package structs

type IDNUM struct {
	ID  int32 `json:"id"`
	Num int32 `json:"num"`
}

type GameItem struct {
	ID         int32 `json:"id"`
	TemplateID int32 `json:"tid"` // 物品模板ID
	Num        int32 `json:"num"` // 物品数量
}

type UsedGameItem struct {
	TemplateID  int32 `json:"tid"`      // 物品模板ID
	LastUseTime int64 `json:"lastdate"` // 最近一次的使用时间
}

const (
	GuidTypes_None = iota // 无
	/// <summary>
	/// 剧情
	/// </summary>
	GuidTypes_Plot = 1
	/// <summary>
	/// 开始冒险
	/// </summary>
	GuidTypes_Adventure = 2
	/// <summary>
	/// 冒险事件1
	/// </summary>
	GuidTypes_AdventureEvent = 3
	/// <summary>
	/// 冒险领取宝箱
	/// </summary>
	GuidTypes_AdventureBox = 4
	/// <summary>
	/// 冒险吃食物
	/// </summary>
	GuidTypes_AdventureFood = 5
	/// <summary>
	/// 英雄出战
	/// </summary>
	GuidTypes_HeroOutFight = 6
	/// <summary>
	/// 武器强化
	/// </summary>
	GuidTypes_HeroWeaponForging = 7
	/// <summary>
	/// 英雄招募
	/// </summary>
	GuidTypes_HeroRecruit = 8
	/// <summary>
	/// 开始挖矿
	/// </summary>
	GuidTypes_MiningDig = 9
	/// <summary>
	/// 挖矿拾取矿石(冒险界面触发->需要定位环节)
	/// </summary>
	GuidTypes_MiningPickOre1 = 10
	/// <summary>
	/// 挖矿拾取矿石(冒险界面触发->无定位环节)
	/// </summary>
	GuidTypes_MiningPickOre2 = 11
	/// <summary>
	/// 挖矿拾取矿石(挖矿界面触发->需要定位环节)
	/// </summary>
	GuidTypes_MiningPickOre3 = 12
	/// <summary>
	/// 挖矿拾取矿石(挖矿界面触发->无定位环节)
	/// </summary>
	GuidTypes_MiningPickOre4 = 13
	/// <summary>
	/// 英雄出战2
	/// </summary>
	GuidTypes_HeroOutFight2 = 14
	/// <summary>
	/// 英雄出战3
	/// </summary>
	GuidTypes_HeroOutFight3 = 15
	/// <summary>
	/// 冒险事件2
	/// </summary>
	GuidTypes_AdventureEvent2 = 16
	/// <summary>
	/// 冒险事件3
	/// </summary>
	GuidTypes_AdventureEvent3 = 17
	/// <summary>
	/// 选择关卡2
	/// </summary>
	GuidTypes_SelectGameLevel2 = 18
	/// <summary>
	/// 选择关卡3
	/// </summary>
	GuidTypes_SelectGameLevel3 = 19
	/// <summary>
	/// 神器锻造
	/// </summary>
	GuidTypes_HeroArtifactForging = 20
	/// <summary>
	/// 完成事件1
	/// </summary>
	GuidTypes_FinishGameLevel = 21
	/// <summary>
	/// 完成事件2
	/// </summary>
	GuidTypes_FinishGameLevel1 = 22
	/// <summary>
	/// 完成事件3
	/// </summary>
	GuidTypes_FinishGameLeve2 = 23
	/// <summary>
	/// 其他
	/// </summary>
	GuidTypes_Other = 99
)
