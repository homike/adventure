package model

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"errors"
)

type Resource struct {
	Strength      int32               // 体力
	Ingot         int32               // 元宝(钻石)
	Money         int32               // 金钱
	Fragments     int32               // 碎片数
	Statue        int32               // 巨魔雕像数
	Detonator     int32               // 雷管
	MiningToolkit int32               // 挖矿工具包
	Ores          *structs.IDNUMARRAY // 矿石
	Foods         *structs.IDNUMARRAY // 食物
	Badges        *structs.IDNUMARRAY // 徽章
	UnlockResIDs  *structs.INT32ARRAY // 已经解锁的资源id列表
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
		Ores:          structs.NewIDNUMARRAY(0),
		Foods:         structs.NewIDNUMARRAY(0),
		Badges:        structs.NewIDNUMARRAY(0),
		UnlockResIDs:  structs.NewINT32ARRAY(0),
	}
}

func (r *Resource) OresChange(id, num int32) {
	r.Ores.Add(id, num)
	r.UnlockResIDs.Add(id)
}

func (r *Resource) StrengthChange(num int32) error {
	if num == 0 {
		return errors.New("StrengthChange error, num == 0")
	}

	r.Strength += num
	if r.Strength < 0 {
		r.Strength = 0
	}
	if r.Strength > gamedata.MaxStrength {
		r.Strength = gamedata.MaxStrength
	}

	return nil
}
