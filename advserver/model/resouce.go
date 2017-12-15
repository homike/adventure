package model

import "adventure/common/structs"

type Resource struct {
	Strength      int32           // 体力
	Ingot         int32           // 元宝(钻石)
	Money         int32           // 金钱
	Fragments     int32           // 碎片数
	Statue        int32           // 巨魔雕像数
	Detonator     int32           // 雷管
	MiningToolkit int32           // 挖矿工具包
	Ores          []structs.IDNUM // 矿石
	Foods         []structs.IDNUM // 食物
	Badges        []structs.IDNUM // 徽章
	UnlockResIDs  []int32         // 已经解锁的资源id列表
}

func NewResource() *Resource {
	return &Resource{
		Strength:      100,
		Ingot:         100,
		Money:         100,
		Fragments:     100,
		Statue:        100,
		Detonator:     100,
		MiningToolkit: 100,
	}
}

func (r *Resource) OresChange(id, num int32) {
	for k, v := range r.Ores {
		if v.ID == id {
			r.Ores[k].Num += num
			if r.Ores[k].Num < 0 {
				r.Ores[k].Num = 0
			}
			return
		}
	}

	r.Ores = append(r.Ores, structs.IDNUM{
		ID:  id,
		Num: num,
	})

	for _, v := range r.UnlockResIDs {
		if v == id {
			return
		}
	}

	r.UnlockResIDs = append(r.UnlockResIDs, id)
}

func (r *Resource) FoodChange(id, num int32) {
	for k, v := range r.Foods {
		if v.ID == id {
			r.Foods[k].Num += num
			if r.Foods[k].Num < 0 {
				r.Foods[k].Num = 0
			}
			return
		}
	}

	r.Foods = append(r.Foods, structs.IDNUM{
		ID:  id,
		Num: num,
	})
}

func (r *Resource) BadgesChange(id, num int32) {
	for k, v := range r.Badges {
		if v.ID == id {
			r.Badges[k].Num += num
			if r.Badges[k].Num < 0 {
				r.Badges[k].Num = 0
			}
			return
		}
	}

	r.Badges = append(r.Badges, structs.IDNUM{
		ID:  id,
		Num: num,
	})
}
