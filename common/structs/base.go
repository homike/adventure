package structs

type IDNUM struct {
	ID  int32 `json:"id"`
	Num int32 `json:"num"`
}

type IDNUMARRAY struct {
	Datas []*IDNUM
}

func NewIDNUMARRAY(count int) *IDNUMARRAY {
	return &IDNUMARRAY{
		Datas: make([]*IDNUM, count),
	}
}

func (i *IDNUMARRAY) Add(id, num int32) {
	for k, v := range i.Datas {
		if v.ID == id {
			i.Datas[k].Num += num
			if i.Datas[k].Num < 0 {
				i.Datas[k].Num = 0
			}
			return
		}
	}

	i.Datas = append(i.Datas, &IDNUM{
		ID:  id,
		Num: num,
	})
}

func (i *IDNUMARRAY) Exist(id, num int32) bool {
	for _, v := range i.Datas {
		if v.ID == id {
			if v.Num >= num {
				return true
			}
			break
		}
	}
	return false
}

func (i *IDNUMARRAY) Count(id int32) int32 {
	for _, v := range i.Datas {
		if v.ID == id {
			return v.Num
		}
	}
	return 0
}

type INT32ARRAY struct {
	Datas []int32
}

func NewINT32ARRAY(count int) *INT32ARRAY {
	return &INT32ARRAY{
		Datas: make([]int32, count),
	}
}

func (i *INT32ARRAY) Add(data int32) {
	for _, v := range i.Datas {
		if v == data {
			return
		}
	}

	i.Datas = append(i.Datas, data)
}

func (i *INT32ARRAY) Exist(data int32) bool {
	for _, v := range i.Datas {
		if v == data {
			return true
		}
	}
	return false
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

type MenuStatus uint8

const (
	MenuStatus_Close MenuStatus = iota // 关闭
	MenuStatus_New                     // 刚解锁
	MenuStatus_Open                    // 已开放
)

type MenuStatusItem struct {
	MenuID     int32      `json:"mid"`     // 菜单ID
	MenuStatus MenuStatus `json:"mstatus"` // 菜单状态
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
	SyncType_First SyncType = iota
	SyncType_Add
	SyncType_Remove
	SyncType_Update
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

type Achievement struct {
	TemplateID int32 // 成就模板ID
	// ConditionType
	// ConditionId
	TotalCount int32      // 累计数
	Status     AchvStatus // 状态
}

type MenuTypes uint8

const (
	MenuTypes_None           MenuTypes = iota
	MenuTypes_Temple                   // 神殿
	MenuTypes_TradeHouse               // 商行
	MenuTypes_MiningWorkshop           // 挖矿-工坊
	MenuTypes_MiningMap                // 挖矿
	MenuTypes_Recruit                  // 招募
	MenuTypes_BaseRoot                 // 基地
	MenuTypes_Arena                    // 竞技场
	MenuTypes_Shop                     // 商城
	MenuTypes_Awaken                   // 觉醒
	MenuTypes_WeaponUpgrade            // 武具升级
	MenuTypes_Artifact                 // 神器
	MenuTypes_UnEmploy                 // 解雇
	MenuTypes_HeroIndex                // 英雄顺序调整
	MenuTypes_FS                       // 封神之阶
	MenuTypes_TradeTroop               // 贸易队
	MenuTypes_Rift                     // 秘境功能开启
)

type ArenaRwdStatus uint8

const (
	ArenaRwdStatus_UnAtive ArenaRwdStatus = iota // 未激活
	ArenaRwdStatus_Ative                         // 已激活
	ArenaRwdStatus_Recieve                       // 已领取
)

type HeroPostion struct {
	HeroTemplateID int32
	HeroIconID     int32
	HeroPosition   int32 // 英雄所在位置
}

type PlayerBaseInfo struct {
	ID            int32  // 玩家ID
	Name          string // 玩家名字
	GameZoneID    int32  // 游戏分区ID
	IconID        int32  // 玩家图标ID
	HP            int32  // 战力
	MaxHP         int32  // 历史最大战力(竞技场)
	Level         int32  // 等级
	LastLoginDate int64  // 最后登录时间
	IsOnline      bool   // 是否在线

	// 战斗信息
	XianGong          int32
	ShanBi            int32
	FangYu            int32
	WangZhe           int32
	ArtfcatWeaponID   int32          // 宿命武器ID
	ArtfcatWeaponName int32          // 宿命武器名字
	OutFightHerosInfo int32          // 出战人数信息
	HeroIDs           []*HeroPostion // 出战英雄
}

func (p *PlayerBaseInfo) GetHeroIDs() []int32 {
	ids := []int32{}
	for _, v := range p.HeroIDs {
		ids = append(ids, v.HeroTemplateID)
	}
	return ids
}

type PlayerGroup struct {
	ID      int32             // 数据库ID
	MinHP   int32             // 最小战斗力
	MaxHP   int32             // 最大战斗力
	Players []*PlayerBaseInfo // 玩家列表
}

type TempleHero struct {
	Quality        QualityType // 品质
	HeroTemplateID int32       // 英雄模板ID
	Cost           int32       // 消耗碎片数量
	IsTrade        bool        // 是否兑换过
}

// 挖矿地图数据
type DigBlockNode struct {
	NodeID    int32 // 地块ID
	X         int8
	Y         int8
	StartTime int64 // 开始时间
}

type BossStatus uint8

const (
	BossStatus_NoAppear BossStatus = 0 // 未出现
	BossStatus_Appear                  // 已出现
	BossStatus_Killed                  // 已击杀
	BossStatus_OverTime                // 已超时
)

type BossNode struct {
	BossID     int32 // BOSS ID
	NodeID     int32 // 地块ID
	X          int8
	Y          int8
	Status     BossStatus
	AppearTime int64 // 出现时间
}

type DigProxy struct {
	ProjectID int32 // 工程ID
	StartTime int64 // 开工时间
}

type BlockNode struct {
	X         int8
	Y         int8
	IsVisible bool
}

type NodeList struct {
	NodeID int32
	Nodes  []*BlockNode
}

func (n *NodeList) AddBlockNode(x, y int8, isVisible bool) {
	n.Nodes = append(n.Nodes, &BlockNode{
		X:         x,
		Y:         y,
		IsVisible: isVisible,
	})
}

type MineMap struct {
	NodeList []*NodeList // 视野内的资源地块列表
	BossIDs  []int32     // 出现过的boss id
	DigCnt   int32       // 挖掘点数
}

type UserMineData struct {
	DigNode         *DigBlockNode // 正在开采的node
	ExpandSight     int32         // 扩展视野, VIP功能扩展
	MinePickLv      int32         // 矿镐等级
	LvMinePickMax   int32         // 挖矿次数上限
	BuyMinePickMax  int32         // 购买挖掘次数上限
	DigCnt          int32         // 挖矿次数
	DigDepthMax     int32         // 最大挖掘深度
	LastResetDate   int64         // 上次刷新地图时间
	LastRefreshDate int64         // 上次刷新耐久度时间
	StatueLv        int32         // 巨魔雕像等级
	StatueCnt       int32         // 巨魔雕像数量
	DigQueueIDs     []int32       // 开采队列等级
	Boss            *BossNode     // Boss信息
	DigProxys       []DigProxy    // 挖矿代理列表
}

type ShopEnum int16

const (
	ShopEnum_DigQueue2       ShopEnum = 1001 // 2级开采队列
	ShopEnum_DigQueue3       ShopEnum = 1002 // 3级开采队列
	ShopEnum_DigQueue4       ShopEnum = 1003 // 4级开采队列
	ShopEnum_DigQueue5       ShopEnum = 1004 // 5级开采队列
	ShopEnum_TradeTaskReset  ShopEnum = 1005 // 商会刷新
	ShopEnum_MiningToolkitEx ShopEnum = 1006 // 神力矿工包，不在商城中显示
	ShopEnum_MiningMapReset  ShopEnum = 1007 // 挖矿地图刷新
	ShopEnum_Chocolate       ShopEnum = 5001 // 巧克力
	ShopEnum_FriedRice       ShopEnum = 5002 // 蛋炒饭
	ShopEnum_Detonator       ShopEnum = 5003 // 雷管
	ShopEnum_MiningToolkit   ShopEnum = 5004 // 神力矿工包
	ShopEnum_MiningPickTop   ShopEnum = 5005 // 矿镐耐久度上限(属性类的)
	ShopEnum_AddGameBoxTop   ShopEnum = 5006 // 增加宝箱上限数量
)
