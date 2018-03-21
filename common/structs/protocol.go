package structs

const (
	Protocol_Test_Req                = 1
	Protocol_Test_Resp               = 2
	Protocol_GetSystemTime_Req       = 3
	Protocol_GetSystemTime_Resp      = 4
	Protocol_LoginServerResult_Ntf   = 1001
	Protocol_CreatePlayer_Req        = 1002 // 创建角色
	Protocol_CreatePlayer_Resp       = 1003
	Protocol_SyncLoginDataFinish_Ntf = 1006
	Protocol_LoginServerPlatform_Req = 1007
	Protocol_SyncPlayerBaseInfo_Ntf  = 1008
	Protocol_NameExists_Req          = 1009
	Protocol_NameExists_Resp         = 1010

	// 英雄相关的消息
	Protocol_Employ_Req         = 1100 // 雇佣英雄
	Protocol_Employ_Resp        = 1101
	Protocol_UnEmploy_Req       = 1102 // 解雇英雄
	Protocol_UnEmploy_Resp      = 1103
	Protocol_ResetHeroIndex_Req = 1104 // 调整英雄出站顺序
	Protocol_SyncHero_Ntf       = 1105 // 同步英雄信息
	Protocol_Work_Req           = 1106 // 英雄出战
	Protocol_SomeWork_Req       = 1107 // 一些英雄出战
	Protocol_Rest_Req           = 1108 // 英雄休息
	Protocol_SomeRest_Req       = 1109 // 一些英雄出战
	Protocol_Work_Resp          = 1110
	//Protocol_SomeWork_Resp          = 1111
	Protocol_Rest_Resp              = 1112
	Protocol_SomeRest_Resp          = 1113
	Protocol_Awake_Req              = 1114 // 英雄觉醒
	Protocol_Awake_Resp             = 1115
	Protocol_UpgradeWeapon_Req      = 1116 // 武具升级
	Protocol_UpgradeWeapon_Resp     = 1117
	Protocol_SyncEmploy_Req         = 1118 // 同步招募信息
	Protocol_SyncEmploy_Resq        = 1119
	Protocol_HeroHpAdd_Ntf          = 1120 // 英雄HP的变化
	Protocol_UnEmployManyHeros_Req  = 1121 // 解雇多名英雄
	Protocol_UnEmployManyHeros_Resp = 1122

	// 冒险相关
	Protocol_SelectGameLevel_Req       = 1200 // 切换关卡
	Protocol_SyncGameLevel_Ntf         = 1201 // 同步关卡基础信息
	Protocol_SyncCurrentGameLevel_Ntf  = 1202 // 同步当前关卡信息
	Protocol_AdventureEvent_Req        = 1203 // 冒险事件
	Protocol_AdventureEvent_Resp       = 1204
	Protocol_OpenGameBox_Req           = 1205
	Protocol_OpenGameBox_Resp          = 1206
	Protocol_OpenGameLevel_Ntf         = 1207 // 解锁关卡通知
	Protocol_SyncCurrentGameLevel2_Ntf = 1208 // 同步当前关卡数据
	Protocol_GetFightCoolingTime_Req   = 1209 // 取得战败冷却时间
	Protocol_GetFightCoolingTime_Resp  = 1210
	Protocol_SpeedAdventure_Req        = 1211
	Protocol_SpeedAdventure_Resp       = 1212
	Protocol_SpeedAdventureReward_Ntf  = 1213

	// 背包相关
	Protocol_UseItem_Req           = 1300 // 使用物品
	Protocol_UseItem_Resp          = 1301
	Protocol_SyncItem_Ntf          = 1302 // 同步物品
	Protocol_SyncBag_Ntf           = 1303 // 同步背包
	Protocol_SyncAllResouce_Ntf    = 1304 // 同步所有的资源
	Protocol_SyncResouce_Ntf       = 1305 // 同步资源
	Protocol_AddItem_Req           = 1306 // 加道具
	Protocol_AddResource_Req       = 1307 // 加资源
	Protocol_AddItem_Resp          = 1308
	Protocol_AddResource_Resp      = 1309
	Protocol_BagNotEnough_Ntf      = 1340
	Protocol_UnlockBag_Req         = 1341 // 开启背包格子
	Protocol_UnlockBag_Resp        = 1342
	Protocol_GetUsedGameItems_Req  = 1343 // 取得已使用过的物品列表
	Protocol_GetUsedGameItems_Resp = 1344

	// 角色相关
	Protocol_EatFood_Req              = 1401 // 吃食物
	Protocol_EatFood_Resp             = 1402
	Protocol_UnLockMenu_Ntf           = 1403 // 解锁菜单
	Protocol_SyncStrength_Ntf         = 1404 // 同步饱足度
	Protocol_SyncWorkHeroTop_Ntf      = 1405 // 同步出站英雄上限
	Protocol_GetEatedFoods_Req        = 1410 // 同步已食用过的食物列表
	Protocol_GetEatedFoods_Resp       = 1411
	Protocol_SyncUnlockMenus_Ntf      = 1412 // 同步已解锁菜单列表
	Protocol_SyncUserGuidRecords_Ntf  = 1413 // 同步新手引导数据
	Protocol_UpdateUserGuidRecord_Req = 1414 // 更新玩家新手引导数据
	Protocol_SyncGameBoxTopNum_Ntf    = 1415 // 更新增加的宝箱上限数量

	// 神殿相关
	Protocol_SyncTemplateHeros_Req     = 1501 // 同步神殿英雄
	Protocol_SyncTemplateHeros_Resp    = 1502
	Protocol_ExchangeTemplateHero_Req  = 1503 // 兑换神殿英雄
	Protocol_ExchangeTemplateHero_Resp = 1504
	Protocol_UnlockTemple_Req          = 1505 // 解锁神殿
	Protocol_UnlockTemple_Resp         = 1506
	Protocol_RefreshTemple_Req         = 1507 // 刷新神殿
	Protocol_RefreshTemple_Resp        = 1508

	// 战斗相关
	Protocol_FightResult_Ntf      = 1600 // 战斗结果
	Protocol_FightRequest_Req     = 1601 // 人机对战请求
	Protocol_FightWithPlayer_Req  = 1602 // 玩家对战请求
	Protocol_FightWithPlayer_Resp = 1603

	// 奖励
	Protocol_RewardResult_Ntf = 1700 // 奖励物品接口

	// 竞技场
	Protocol_OpenArena_Req           = 1800 // 请求打开竞技场
	Protocol_SyncArena_Ntf           = 1801 // 同步竞技场
	Protocol_ArenaChallenge_Req      = 1802 // 发起挑战
	Protocol_ArenaChallenge_Resp     = 1803
	Protocol_RefreshArena_Req        = 1804 // 刷新竞技场
	Protocol_RefreshArena_Resp       = 1805
	Protocol_RecieveArenaReward_Req  = 1806 // 领取竞技场奖励
	Protocol_RecieveArenaReward_Resp = 1807

	// 游戏商城

	// 挖矿
	Protocol_DigRequest_Req                     = 2100 // 挖矿
	Protocol_DigRequest_Resp                    = 2101
	Protocol_PickOre_Req                        = 2102 // 拾取矿
	Protocol_PickOre_Resp                       = 2103
	Protocol_SyncMiningMap_Req                  = 2104 // 同步挖矿地图
	Protocol_SyncMiningMap_Resp                 = 2105
	Protocol_UseDetonator_Req                   = 2106 // 使用雷管
	Protocol_UseDetonator_Resp                  = 2107
	Protocol_UseMiningToolkit_Req               = 2108 // 使用挖矿工具包
	Protocol_UseMiningToolkit_Resp              = 2109
	Protocol_ResetMiningMap_Req                 = 2110 // 刷新矿区请求
	Protocol_ResetMiningMap_Resp                = 2111
	Protocol_UseDetonatorForRock_Req            = 2112 // 使用雷管炸岩石
	Protocol_UseDetonatorForRock_Resp           = 2113
	Protocol_FightMonster_Req                   = 2114 // 挑战怪物
	Protocol_RefreshBoss_Req                    = 2115 // 刷新BOSS
	Protocol_BuyDigTools_Req                    = 2116 // 购买挖矿辅助工具
	Protocol_BuyDigTools_Resp                   = 2117
	Protocol_UpdateTileStatus_Ntf               = 2118 // 更新地块状态
	Protocol_FightMonster_Resp                  = 2119
	Protocol_BuyMiningPickTop_Req               = 2120 // 购买矿镐耐久度上限
	Protocol_BuyMiningPickTop_Resp              = 2121
	Protocol_BuyDigQueue_Req                    = 2122 // 购买开采队列
	Protocol_BuyDigQueue_Resp                   = 2123
	Protocol_AddDiggingProxy_Req                = 2124 // 增加开采代理
	Protocol_AddDiggingProxy_Resp               = 2125
	Protocol_GetDiggingProxy_Req                = 2126 // 收获开采代理
	Protocol_PreDigTest_Req                     = 2128 // 预开采
	Protocol_PreDigTest_Resp                    = 2129
	Protocol_UpdateMiningMapBaseData_Ntf        = 2130 // 更新矿区基本信息
	Protocol_UpdateStatueLevelAndCount_Ntf      = 2131 // 更新巨魔雕像数量
	Protocol_RequestCoolingTime_Req             = 2132 // 获取战斗冷却时间
	Protocol_RequestCoolingTime_Resp            = 2133
	Protocol_SyncMiningMapDigAgentResource_Req  = 1234 // 请求同步挖矿代理可获得矿石种类及数量
	Protocol_SyncMiningMapDigAgentResource_Resp = 2135
	Protocol_BuyMiningToolkit_Req               = 2136 // 购买巨力工具包

	// 神器相关
	//Protocol_UnlockArtifactSeal_Req     = 1900 // 解锁神器封印
	//Protocol_UnlockArtifactSeal_Resp    = 1901
	Protocol_EquipArtifact_Req          = 1902 // 装备神器
	Protocol_EquipArtifact_Resp         = 1903
	Protocol_UpgradeArtifact_Req        = 1904 // 升级神器
	Protocol_UpgradeArtifact_Resp       = 1905
	Protocol_SyncArtifactStatus_Ntf     = 1906 // 同步神器状态
	Protocol_SyncArtifactSealStatus_Ntf = 1907 // 同步神器封印状态

	// 成就相关
	Protocol_GetAchievements_Req      = 2201 // 取得成就记录数据
	Protocol_GetAchievements_Resp     = 2202
	Protocol_RecieveAchievements_Req  = 2203 // 领取成就奖励
	Protocol_RecieveAchievements_Resp = 2204
	Protocol_UpdateAchievement_Ntf    = 2205 // 更新成就状态

	// 公告
	Protocol_SystemAnnouncement_Ntf     = 2701 // 系统公告
	Protocol_SystemAnnouncementRich_Ntf = 2702 // 富文本系统公告
	Protocol_SetPlayerBarrageConfig_Req = 2801 // 弹幕设置

)

///////////////////////////////////////////// 系统 ////////////////////////////////////////
type GetSystemTimeReq struct {
}

type GetSystemTimeResp struct {
	Time int64
}

type LoginServerResultNtf struct {
	Result         int32 // 0: Success
	IsCreatePlayer bool
}

type CreatePlayerReq struct {
	PlayerName     string
	HeroTemplateId int32
}

type CreatePlayerResp struct {
	Result int32 // 0: Success
}

type SyncLoginDataFinishNtf struct {
}

type LoginServerPlatformReq struct {
	Takon     string
	Version   int32
	ChannelID string
}

type SyncPlayerBaseInfoNtf struct {
	PlayerID           int32
	GameZoonID         int32 // 游戏分区ID
	IsSupperMan        bool  // 是否是GM
	PlatformType       int32 // 平台类型
	Viplevel           int32
	TotalRechargeIngot int32
}

type NameExistsReq struct {
	Name string
}

// 如果存在，则返回一个新名字，如果和传入的名字一样，则说明没有重名
type NameExistsResp struct {
	Name string
}

type GuildRecord struct {
	UserGuidTypes uint8
	TriggerCount  int32
}

type UpdateUserGuidRecordReq struct {
	UserGuidTypes uint8
}

type SyncUserGuidRecordsNtf struct {
	Records []*GuildRecord
}

type SyncStrengthNtf struct {
	Strength int32
}

///////////////////////////////////////////// 英雄 ////////////////////////////////////////
type SyncHeroWorkTopNtf struct {
	MaxWorker int32
}

type SyncUnlockMenuNtf struct {
	MenuID int32
}

type SyncUnlockMenusNtf struct {
	MenuStates []*MenuStatusItem
}

type SyncGameBoxTopNumNtf struct {
	AddNum int32
}

type SyncHeroNtf struct {
	SyncHeroType uint8
	Heros        []*Hero
}

const (
	AddHP_Type_Item        = iota // 物品添加
	AddHP_Type_Event              // 事件添加
	AddHP_Type_Achievement        // 成就添加
	AddHP_Type_Food               // 吃食物
	AddHP_Type_WinBattle          // 打副本胜利
)

type HeroHPAddNtf struct {
	Type   uint8
	HeroID int32
	AddHP  int32
}

type EmployReq struct {
	EmployType uint8
}

const (
	EmployRet_Success = iota
	EmployRet_Failed
	EmployRet_NotEnough
)

type EmployResp struct {
	Ret     uint8
	HeroIDs []int32
}

type UnEmployReq struct {
	HeroID int32
}
type UnEmployResp struct {
	Ret    uint8
	HeroID int32
}

type RewardResultNtf struct {
	IsRes   bool
	Rewards []*Reward
	Context string
}

type ResetHeroIndexReq struct {
	HeroIDs []int32
}

type WorkReq struct {
	HeroID int32
}

type WorkResp struct {
	Ret    uint8
	HeroID int32
}

type SomeWorkReq struct {
	HeroIDs []int32
}

type ResetResp struct {
	Ret    uint8
	HeroID int32
}

type ResetReq struct {
	HeroID int32
}

type SomeResetReq struct {
	HeroIDs []int32
}

type AwakeReq struct {
	HeroID int32
}

type AwakeResp struct {
	Ret    uint8
	HeroID int32
	AddHP  int32
}

type UpgradeWeaponReq struct {
	HeroID int32
	Ingot  int32
}

type UpgradeWeaponResp struct {
	Ret    uint8
	HeroID int32
	AddHP  int32
}

type SyncEmployReq struct {
}

type SyncEmployResp struct {
	Type                       []int32
	Cost                       []int32
	LeftSecond                 int32
	NextFreeIngotEmployLeftSec int32
}

///////////////////////////////////////////// 背包 ////////////////////////////////////////
type SyncAllResouceNtf struct {
	Money         int32
	Ingot         int32
	Fragment      int32
	Statue        int32
	Strength      int32
	Detonator     int32
	MiningToolkit int32
	Ors           []*IDNUM
	Foods         []*IDNUM
	Badges        []*IDNUM
	UnlockResIds  []int32
}

type SyncResourceNtf struct {
	ResID  int32
	ResNum int32
}

type SyncBagNtf struct {
	MaxCount    int32
	UnlockLevel int32
	Items       []*GameItem
}

const (
	SyncItem_Type_First = iota
	SyncItem_Type_Add
	SyncItem_Type_Remove
	SyncItem_Type_Update
)

type SyncItemNtf struct {
	Type  uint8
	Items []*GameItem
}

type GetUsedGameItemsResp struct {
	ItemIDs   []int32
	UserTimes []int64
}

type UseItemReq struct {
	ItemID int32
}

type UseItemResp struct {
	Ret uint8
}

type UnlockBagResp struct {
	Ret         uint8
	MaxCount    int32
	UnLockLevel int32
}

type EatFoodReq struct {
	FoodID int32
}

type EatFoodResp struct {
	Ret      uint8
	Strength int32
}

type GetEatedFoodsReq struct {
}

type GetEatedFoodsResp struct {
	FoodIDs   []int32
	EatedDate []int64
}

///////////////////////////////////////////// 冒险 ////////////////////////////////////////

type SyncGameLevelNtf struct {
	GameLevels      []*GameLevel
	CurLevelID      int32
	LastRefreshTime int64
	SpeedCount      int32 // 加速冒险的次数
}

type SyncCurrentGameLevelNtf struct {
	GameLevel *GameLevel
}

type SyncCurrentGameLevelNtf2 struct {
	GameLevel *GameLevel
}

type GetFightCoolingTimeReq struct {
	BattleID int32
}

type GetFightCoolingTimeResp struct {
	LeftTime int32
}

type SelectGameLevelReq struct {
	LevelID int32
}

type OpenGameBoxReq struct {
	Count int32
}

type OpenGameBoxResp struct {
	Ret     uint8
	Count   int32
	Rewards []*Reward
}

///////////////////////////////////////////// 冒险 ////////////////////////////////////////

type AdventureEventReq struct {
	EventID int32
}

type AdventureEventResp struct {
	Ret         uint8
	GameLevelID int32
	EventID     int32
}

type OpenGameLevelNtf struct {
	GameLevel *GameLevel
}

///////////////////////////////////////////// 神器 ////////////////////////////////////////
type SyncArtifactStatusNtf struct {
	SType  SyncType
	Status []*ArtifactStatus
}

type SyncArtifactSealStatusNtf struct {
	SType  SyncType
	Status []*ArtifactSealStatus
}

type UnlockArtifactSealReq struct { // 解锁一个神器封印
	sealID uint8
}

type UnlockArtifactSealResp struct {
	Ret uint8
}

type EquipArtifactReq struct { // 装备神器
	ArtifactID int32
}

type EquipArtifactResp struct {
	Ret        uint8
	ArtifactID int32
}

type UpgradeArtifactReq struct { // 神器升级
	Ingot int32
}
type UpgradeArtifactResp struct {
	Ret   uint8
	AddHP int32
}

///////////////////////////////////////////// 战斗 ////////////////////////////////////////
type FightResultNtf struct {
	FType  FightType
	Result *FightResult
}

///////////////////////////////////////// 成就相关 ////////////////////////////////////////
type GetAchievementsReq struct { // 取得成就记录数据
}

type GetAchievementsResp struct {
	Achievements          []*Achievement
	NextRefreshTimeDaily  int64 // 下一次刷新每日成就的时间
	NextRefreshTimeWeekly int64 // 下一次刷新每周成就的时间
}

type RecieveAchievementsReq struct{} // 领取成就奖励

type RecieveAchievementsResp struct{}

type UpdateAchievementNtf struct { // 更新成就状态
	Achievements []*Achievement
}

///////////////////////////////////////// 公告 ////////////////////////////////////////

type SystemAnnouncementNtf struct {
	Texts           []string
	Colors          []string
	LoopCount       int32
	IsLeftDirection bool
}

type SystemAnnouncementRichNtf struct {
	Text            string
	LoopCount       int32
	IsLeftDirection bool
}

///////////////////////////////////////// 竞技场 ////////////////////////////////////////

type OpenArenaReq struct {
}

type SyncArenaNtf struct {
	Targets         []*FightTarget // 可攻击目标
	LessRefreshTime int32          // 剩余的刷新时间
	ChallengeCount  int32          // 已经用掉的挑战次数
}

type ArenaChallengeReq struct { // 竞技场挑战
	PlayerID int32
}

type ArenaChallengeResp struct {
	Ret      uint8
	PlayerID int32
	IsWin    bool
}

type ArenaRefreshReq struct { // 刷新竞技场
	UserIngot bool
}

type ArenaRefreshResp struct {
	Ret uint8
}

type ArenaRecieveRewardReq struct { // 领取竞技场奖励
	StarIndex int
}

type ArenaRecieveRewardResp struct {
	RewardRecords []ArenaRwdStatus
}

///////////////////////////////////////// 神殿 ////////////////////////////////////////

type SyncTemplateHerosReq struct {
}
type SyncTemplateHerosResp struct {
	Heros        []*TempleHero
	LeftSecond   int32
	TradeCount   int32
	RefreshCount int32
}

type ExchangeTemplateHeroReq struct {
	HeroTemplateID int32
}
type ExchangeTemplateHeroResp struct {
	Ret            uint8
	HeroTemplateID int32
}

type UnlockTempleReq struct {
}
type UnlockTempleResp struct {
	Ret uint8
}

type RefreshTempleReq struct {
	Count int32
}
type RefreshTempleResp struct {
	Ret        uint8
	SplashGold bool
}

///////////////////////////////////////// 挖矿 ////////////////////////////////////////
type UpdateStatueLevelAndCountNtf struct {
	Level int32
	Count int32
}

type SyncMiningMapResp struct {
	NodeBlocks []*NodeList
	MapBase    *UserMineData
}

type ResetMiningMapReq struct {
	UserIgnot bool
}
type ResetMiningMapResp struct {
	Ret uint8
}

type DigReq struct {
	X                 int8
	Y                 int8
	IsOnlyExpandSight bool
}

type DigResp struct {
}
