package structs

type QualityType int32 // 英雄品质类型
const (
	None       QualityType = 0
	White      QualityType = 0
	Green      QualityType = 1
	Blue       QualityType = 2
	Purple     QualityType = 3
	Gold       QualityType = 4
	SplashGold QualityType = 5
)

type EmployType int32 // 雇佣类型
const (
	Money        EmployType = 0 // 用钱
	HunLuan      EmployType = 1 // 混乱之门
	HuiHuang     EmployType = 2 // 辉煌之门
	LvDong       EmployType = 3 // 律动之门
	Diamond      EmployType = 4 // 万象之门
	ManyDiamond  EmployType = 5 // 传奇之门(10连抽）
	ManyDiamond2 EmployType = 6 // 传奇之门(10连抽特殊，保证必须出一个紫色英雄）
	Exchange     EmployType = 7 // 碎片兑换
	Reward       EmployType = 8 // 系统奖励
)

type Hero struct {
	HeroId         int32       // 英雄id
	HeroTemplateId int32       // 英雄的模板id
	Name           string      // 英雄的名字
	Level          int32       // 当前等级
	IsOutFight     bool        // 出战状态
	IsPlayer       bool        // 是否是玩家
	Exp            int32       // 当前的经验
	MaxExp         int32       // 最大经验
	Quality        QualityType // 品质类型
	AwakeCount     int32       // 觉醒次数
	WeaponLevel    int32       // 武具等级
	Index          int32       // 排序索引
	LevelHP        int32       // 因为升级而改变的HP
	ItemHP         int32       // 因为物品而改变的HP
	HP             int32
	// {
	//     get
	//     {
	//         return LevelHP;
	//     }
	// }
}
