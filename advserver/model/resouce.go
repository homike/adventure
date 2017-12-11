package model

import "adventure/common/structs"

type Resource struct {
	Strength           int32           // 体力
	Money              int32           // 金钱
	Fragments          int32           // 碎片数
	Statue             int32           // 巨魔雕像数
	Detonator          int32           // 雷管
	MiningToolkit      int32           // 挖矿工具包
	Ores               []structs.IDNUM // 矿石
	Foods              []structs.IDNUM // 食物
	Badges             []structs.IDNUM // 徽章
	BUnlockResIdsadges []int32         // 已经解锁的资源id列表
}

func NewResource() *Resource {
	return &Resource{}
}
