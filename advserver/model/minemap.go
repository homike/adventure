package model

import (
	"adventure/common/structs"
)

type UserMineData struct {
	MineMap         *structs.MineMap      // 挖矿地图
	DigNode         *structs.DigBlockNode // 正在开采的node
	ExpandSight     int32                 // 扩展视野, VIP功能扩展
	MinePickLv      int32                 // 矿镐等级
	LvMinePickMax   int32                 // 挖矿次数上限
	BuyMinePickMax  int32                 // 购买挖掘次数上限
	DigCnt          int32                 // 挖矿次数
	DigDepthMax     int32                 // 最大挖掘深度
	LastResetDate   int32                 // 上次刷新地图时间
	LastRefreshDate int32                 // 上次刷新耐久度时间
	StatueLv        int32                 // 巨魔雕像等级
	StatueCnt       int32                 // 巨魔雕像数量
	DigQueueIDs     []int32               // 开采队列等级
	Boss            *structs.BossNode     // Boss信息
	DigProxys       []structs.DigProxy    // 挖矿代理列表
}

func NewUserMineData() *UserMineData {
	return &UserMineData{
		MineMap: &structs.MineMap{
			NodeList: []structs.NodeList{},
			BossIDs:  []int32{},
			DigCnt:   0,
		},
		ExpandSight: 0,
		StatueLv:    1,
		DigNode: &structs.DigBlockNode{
			NodeID: -1,
			X:      127,
		},
		MinePickLv:      1,
		LastResetDate:   0,
		LastRefreshDate: 0,
		Boss:            &structs.BossNode{},
		DigQueueIDs:     []int32{},
		DigProxys:       []structs.DigProxy{},
	}
}

func (u *UserMineData) ResetMineMap() {
	u.MineMap.NodeList = []structs.NodeList{}

}
