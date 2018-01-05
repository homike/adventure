package model

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"errors"
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
			Status:     structs.ArtifactStatusType_UnLock,
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
		if v.Status == structs.ArtifactStatusType_Use {
			return v
		}
	}
	return nil
}

func (a *Artifact) GetArtifactStatusUnLockCount() int32 {
	cnt := int32(0)
	for _, v := range a.Status {
		if v.Status != structs.ArtifactStatusType_Lock {
			cnt++
		}
	}
	return cnt
}

func (a *Artifact) UnlockArtifact(id int32) error {
	status := a.GetArtifactStatus(id)
	if status == nil {
		return errors.New("UnlockArtifact faield")
	}
	status.Status = structs.ArtifactStatusType_New
	return nil
}

func (a *Artifact) ArtifactSealStatus(id int32) *structs.ArtifactSealStatus {
	for _, v := range a.SealStatus {
		if v.SeaID == id {
			return v
		}
	}
	return nil
}
