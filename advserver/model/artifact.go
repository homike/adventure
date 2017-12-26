package model

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
)

type Artifact struct {
	Status     []*structs.ArtifactStatus     // 神器的状态
	SealStatus []*structs.ArtifactSealStatus // 封印解锁的状态
}

func NewArtifact() *Artifact {

	artifact := &Artifact{
		Status:     []*structs.ArtifactStatus{},
		SealStatus: []*structs.ArtifactSealStatus{},
	}

	for _, v := range gamedata.AllTemplates.SpellTemplates {
		artifact.Status = append(artifact.Status, &structs.ArtifactStatus{
			ArtifactID: v.ID,
			Status:     structs.UnLock,
		})
	}

	return artifact
}

func (a *Artifact) GetArtifactStatus(id int32) *structs.ArtifactStatus {
	for _, v := range a.Status {
		if v.ArtifactID == id {
			return v
		}
	}
	return nil
}

func (a *Artifact) GetArtifactStatusUse() *structs.ArtifactStatus {
	for _, v := range a.Status {
		if v.Status == structs.Use {
			return v
		}
	}
	return nil
}

func (a *Artifact) GetArtifactStatusUnLockCount() int32 {
	cnt := int32(0)
	for _, v := range a.Status {
		if v.Status != structs.Lock {
			cnt++
		}
	}
	return cnt
}

func (a *Artifact) ArtifactSealStatus(id int32) *structs.ArtifactSealStatus {
	for _, v := range a.SealStatus {
		if v.SeaID == id {
			return v
		}
	}
	return nil
}
