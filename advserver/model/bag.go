package model

import (
	"adventure/common/structs"
)

type Bag struct {
	MaxItemID     int32                  // 物品索引ID
	MaxCount      int32                  // 玩家的背包最大数量
	Items         []structs.GameItem     // 玩家的物品列表
	UnlockLevel   int32                  // 背包的解锁等级
	UsedGameItems []structs.UsedGameItem // 用过的物品列表
}

func NewBag() *Bag {
	bag := &Bag{
		MaxItemID:     0,
		MaxCount:      16,
		Items:         make([]structs.GameItem, 0, 10),
		UsedGameItems: make([]structs.UsedGameItem, 0, 10),
	}

	return bag
}
