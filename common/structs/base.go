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

type MenuStatusItem struct {
	MenuID     int32 `json:"mid"`     // 菜单ID
	MenuStatus uint8 `json:"mstatus"` // 菜单状态
}

type Reward struct {
	RewardType uint8
	Param1     int32
	Param2     int32
}

const (
	RewardType_Property            = 0  // 属性奖励
	RewardType_RandProperty        = 1  // 随机属性奖励
	RewardType_HP                  = 2  // 战斗力奖励
	RewardType_Item                = 3  // 奖励物品
	RewardType_Exp                 = 4  // 经验奖励
	RewardType_UnlockGameLevel     = 5  // 解锁游戏关卡
	RewardType_Hero                = 6  // 英雄奖励
	RewardType_UnlockMenu          = 7  // 解锁菜单
	RewardType_AddHeroWorkTop      = 8  // 增加英雄出战数上限
	RewardType_AddMiningPickNumTop = 9  // 增加挖掘次数上限
	RewardType_AddMiningPickLevel  = 10 // 增加矿镐等级
	RewardType_AddMiningPickNum    = 11 // 增加挖矿次数，无视上限
	RewardType_UnlockArtifact      = 12 // 解锁神器
	RewardType_AddGetGiftDayNum    = 13 // 增加好友中每日领取礼物次数
	RewardType_AddSendGiftDayNum   = 14 // 增加好友中每日送礼次数
	RewardType_TradeTaskReset      = 15 // 商会任务重置
	RewardType_AddGameBoxNumTop    = 16 // 增加宝箱上限
)

const (
	MenuStatus_Close = iota // 关闭
	MenuStatus_New          // 刚解锁
	MenuStatus_Open         // 已开放
)

const (
	GuidTypes_None                = iota // 无
	GuidTypes_Plot                = 1    // 剧情
	GuidTypes_Adventure           = 2    // 开始冒险
	GuidTypes_AdventureEvent      = 3    // 冒险事件1
	GuidTypes_AdventureBox        = 4    // 冒险领取宝箱
	GuidTypes_AdventureFood       = 5    // 冒险吃食物
	GuidTypes_HeroOutFight        = 6    // 英雄出战
	GuidTypes_HeroWeaponForging   = 7    // 武器强化
	GuidTypes_HeroRecruit         = 8    // 英雄招募
	GuidTypes_MiningDig           = 9    // 开始挖矿
	GuidTypes_MiningPickOre1      = 10   // 挖矿拾取矿石(冒险界面触发->需要定位环节)
	GuidTypes_MiningPickOre2      = 11   // 挖矿拾取矿石(冒险界面触发->无定位环节)
	GuidTypes_MiningPickOre3      = 12   // 挖矿拾取矿石(挖矿界面触发->需要定位环节)
	GuidTypes_MiningPickOre4      = 13   // 挖矿拾取矿石(挖矿界面触发->无定位环节)
	GuidTypes_HeroOutFight2       = 14   // 英雄出战2
	GuidTypes_HeroOutFight3       = 15   // 英雄出战3
	GuidTypes_AdventureEvent2     = 16   // 冒险事件2
	GuidTypes_AdventureEvent3     = 17   // 冒险事件3
	GuidTypes_SelectGameLevel2    = 18   // 选择关卡2
	GuidTypes_SelectGameLevel3    = 19   // 选择关卡3
	GuidTypes_HeroArtifactForging = 20   // 神器锻造
	GuidTypes_FinishGameLevel     = 21   // 完成事件1
	GuidTypes_FinishGameLevel1    = 22   // 完成事件2
	GuidTypes_FinishGameLeve2     = 23   // 完成事件3
	GuidTypes_Other               = 99   // 其他
)

type EmployType int32 // 雇佣类型
const (
	EmployType_Money        EmployType = 0 // 金币抽奖
	EmployType_HunLuan      EmployType = 1 // 混乱之门
	EmployType_HuiHuang     EmployType = 2 // 辉煌之门
	EmployType_LvDong       EmployType = 3 // 律动之门
	EmployType_Diamond      EmployType = 4 // 万象之门
	EmployType_ManyDiamond  EmployType = 5 // 传奇之门(10连抽）
	EmployType_ManyDiamond2 EmployType = 6 // 传奇之门(10连抽特殊，保证必须出一个紫色英雄）
	EmployType_Exchange     EmployType = 7 // 碎片兑换
	EmployType_Reward       EmployType = 8 // 系统奖励
)

const (
	SyncHeroType_Add    = iota // 添加
	SyncHeroType_Update        // 更新
	SyncHeroType_First         // 首次更新
)

const (
	AdventureRet_Success = iota
	AdventureRet_Failed
)

const (
	ItemType_Fix    = iota // 属性固定
	ItemType_Random        // 属性随机
)

const (
	ResourceType_Money             = 1   // 金币
	ResourceType_Ingot             = 2   // 钻石
	ResourceType_Fragment          = 3   // 碎片
	ResourceType_Strength          = 4   // 饱食度(体力)
	ResourceType_Statue            = 5   // 巨魔雕像ID
	ResourceType_Detonator         = 6   // 雷管
	ResourceType_MiningToolkit     = 7   // 挖矿工具包
	ResourceType_MiningPickTop     = 8   // 矿镐耐久度上限(属性类的)
	ResourceType_Coupon            = 9   // 礼券
	ResourceType_Exp               = 15  // 经验值
	ResourceType_TradeTask         = 16  // 商会任务
	ResourceType_GameBoxTop        = 17  // 宝箱数量上限
	ResourceType_MiningPickup      = 18  // 矿镐耐久度
	ResourceType_MiningPickLevel   = 19  // 矿镐等级
	ResourceType_FriendGetGiftNum  = 20  // 好友收礼上限
	ResourceType_FriendSendGiftNum = 21  // 好友送礼上限
	ResourceType_BagLimit          = 22  // 背包上限
	ResourceType_OriMin            = 100 // 矿产资源ID最小值
	ResourceType_OriMax            = 150 // 矿产资源ID最大值
	ResourceType_FoodMin           = 200 // 食物资源id最小值
	ResourceType_FoodMax           = 250 // 食物资源id最大值
	ResourceType_BadgesMin         = 300 // 徽章资源最小值
	ResourceType_BadgesMax         = 350 // 徽章资源最大值
	ResourceType_HunLuanCrest      = 301 // 混乱之门徽章id
	ResourceType_HuiHuangCrest     = 302 // 辉煌之门徽章id
	ResourceType_LvDongCrest       = 303 // 律动之门徽章id
)

const (
	ResouceChangeType_Employ_Money = iota // 招募_金币招募
)

type FightResult struct {
	//ID        int32         // 战报日志
	//Time      int64         // 日期
	LeftTeam     *FightTeam   // 左边队伍
	RightTeam    *FightTeam   // 左边队伍
	Rounds       []FightRound // 回合数据
	IsLeftWin    bool         // 是否是左边胜利
	BackgroundID string       //
	ForegroundID string       // 前景ID
}

type FightTeam struct {
	XianGong  int32   // 先攻
	FangYu    int32   // 防御
	ShanBi    int32   // 闪避
	WangZhe   int32   // 王者
	DefaultHP int32   // 获取初始HP
	Models    []int32 // 英雄模板id
	SpellIDs  []int32 // 技能ID列表
	Name      string  // 名字
}

type FightRound struct {
	SkillID int32 // 技能ID
	LeftHP  int32
	RightHP int32
}

type FightType int8

const (
	FightType_GMTest         FightType = -1
	FightType_EventFight               = 0 // 关卡战斗事件
	FightType_ArenaPK                  = 1 // 竞技场
	FightType_MiningMonster            = 2 // 挖矿地图怪物战
	FightType_MiningBoss               = 3 // 挖矿地图boss战
	FightType_TestFight                = 4 // 模拟战斗
	FightType_OrderTestFight           = 5 // 排行榜中的挑战
	FightType_FSPK                     = 6 // 封神之阶
)

type SyncType uint8

const (
	First SyncType = iota
	Add
	Remove
	Update
)

type OfflineReward struct {
	OfflineTimeSec int32   // 离线时间 单位：秒
	HasStrength    bool    // 是否还有饱足度
	Money          int32   // 获得的游戏币
	Exp            int32   // 获得的经验
	OfflineHP      int32   // 离线前的HP
	OnlineHP       int32   // 上线时的HP
	UpLevelHero    []int32 // 产生升级的英雄
}

type ActiveType uint8

const (
	ActiveType_Full ActiveType = iota // 正常收益
	ActiveType_Half                   // 减半收益
	ActiveType_None                   // 无收益
)
