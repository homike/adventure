package model

import (
	"adventure/common/structs"
	"errors"
)

type Bag struct {
	MaxItemID     int32                   // 物品索引ID
	MaxCount      int32                   // 玩家的背包最大数量
	Items         []*structs.GameItem     // 玩家的物品列表
	UnlockLevel   int32                   // 背包的解锁等级
	UsedGameItems []*structs.UsedGameItem // 用过的物品列表
}

func NewBag() *Bag {
	bag := &Bag{
		MaxItemID:     1,
		MaxCount:      16,
		Items:         make([]*structs.GameItem, 0, 10),
		UsedGameItems: make([]*structs.UsedGameItem, 0, 10),
	}

	return bag
}

func (b *Bag) GetItem(id int32) (*structs.GameItem, error) {
	for _, v := range b.Items {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, errors.New("GetItem failed")
}

func (b *Bag) RemoveItem(id int32) error {
	for k, v := range b.Items {
		if v.ID == id {
			b.Items = append(b.Items[0:k], b.Items[k+1:]...)
			return nil
		}
	}
	return errors.New("RemoveItem failed")
}

func (b *Bag) GetUsedItem(templateID int32) (*structs.UsedGameItem, error) {
	for _, v := range b.UsedGameItems {
		if v.TemplateID == templateID {
			return v, nil
		}
	}
	return nil, errors.New("GetUsedItem failed")
}

func (b *Bag) AddUsedItem(item *structs.UsedGameItem) error {
	b.UsedGameItems = append(b.UsedGameItems, item)
	return nil
}
