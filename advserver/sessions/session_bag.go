package sessions

import (
	"adventure/common/structs"
)

func (sess *Session) ResourceChange(resID int32, resNum int32, resChangeType int32) {
	switch {
	case resID == structs.ResourceType_Money: // 金币
		sess.MoneyChange(resNum, resChangeType)

	case resID == structs.ResourceType_Ingot: // 钻石
		sess.IgnotChange(resNum, resChangeType)

	case resID == structs.ResourceType_Fragment: // 碎片
		sess.FregmentsChange(resNum, resChangeType)

	case resID == structs.ResourceType_Strength: // 体力(饱食度)
		sess.StrengthChange(resNum, resChangeType)

	case resID == structs.ResourceType_Statue: // 巨魔雕像
		sess.StatueChange(resNum, resChangeType)

	case resID == structs.ResourceType_Detonator: // 雷管
		sess.DetonatorChange(resNum, resChangeType)

	case resID == structs.ResourceType_MiningToolkit: // 挖矿工具包
		sess.MiningToolkitChange(resNum, resChangeType)

	case resID == structs.ResourceType_MiningPickTop: // 耐久度上限
		//
	case resID >= structs.ResourceType_OriMin && resID <= structs.ResourceType_OriMax: // 矿石
		sess.PlayerData.Res.OresChange(resID, resNum)
		sess.SyncResource(resID, resNum)

	case resID >= structs.ResourceType_FoodMin && resID <= structs.ResourceType_FoodMax: // 食物
		sess.PlayerData.Res.FoodChange(resID, resNum)
		sess.SyncResource(resID, resNum)

	case resID >= structs.ResourceType_BadgesMin && resID <= structs.ResourceType_BadgesMax: // 徽章
		sess.PlayerData.Res.BadgesChange(resID, resNum)
		sess.SyncResource(resID, resNum)
	}
}

func (sess *Session) SyncAllResources() {
	resp := &structs.SyncAllResouceNtf{
		Money:         sess.PlayerData.Res.Money,
		Ingot:         sess.PlayerData.Res.Ingot,
		Fragment:      sess.PlayerData.Res.Fragments,
		Statue:        sess.PlayerData.Res.Statue,
		Strength:      sess.PlayerData.Res.Strength,
		Detonator:     sess.PlayerData.Res.Detonator,
		MiningToolkit: sess.PlayerData.Res.MiningToolkit,
		Ors:           sess.PlayerData.Res.Ores,
		Foods:         sess.PlayerData.Res.Foods,
		Badges:        sess.PlayerData.Res.Badges,
		UnlockResIds:  sess.PlayerData.Res.UnlockResIDs,
	}

	sess.Send(structs.Protocol_SyncAllResouce_Ntf, resp)
}

func (sess *Session) SyncResource(resID int32, resNum int32) {
	resp := &structs.SyncResourceNtf{
		ResID:  resID,
		ResNum: resNum,
	}
	sess.Send(structs.Protocol_SyncResouce_Ntf, resp)
}

func (sess *Session) SyncBag() {
	resp := &structs.SyncBagNtf{
		MaxCount:    sess.PlayerData.Bag.MaxCount,
		UnlockLevel: sess.PlayerData.Bag.UnlockLevel,
		Items:       sess.PlayerData.Bag.Items,
	}

	sess.Send(structs.Protocol_SyncBag_Ntf, resp)
}

func (sess *Session) SyncPlayerUsedItem() {

	itemIDS := []int32{}
	userTimes := []int64{}
	for _, v := range sess.PlayerData.Bag.UsedGameItems {
		itemIDS = append(itemIDS, v.TemplateID)
		userTimes = append(userTimes, v.LastUseTime)
	}

	resp := &structs.GetUsedGameItemsResp{
		ItemIDs:   itemIDS,
		UserTimes: userTimes,
	}
	sess.Send(structs.Protocol_GetUsedGameItems_Resp, resp)
}

func (sess *Session) MoneyChange(cnt int32, changeType int32) {
	if cnt == 0 {
		return
	}
	sess.PlayerData.Res.Money += cnt
	sess.SyncResource(structs.ResourceType_Money, sess.PlayerData.Res.Money)
}

func (sess *Session) CouponChange(cnt int32, changeType int32) {
	if cnt == 0 {
		return
	}

	sess.PlayerData.Res.Ingot += cnt
	sess.SyncResource(structs.ResourceType_Ingot, sess.PlayerData.Res.Ingot)
}

func (sess *Session) IgnotChange(cnt int32, changeType int32) {
	if cnt == 0 {
		return
	}

	sess.PlayerData.Res.Ingot += cnt
	sess.SyncResource(structs.ResourceType_Ingot, sess.PlayerData.Res.Ingot)
}

func (sess *Session) FregmentsChange(cnt int32, changeType int32) {
	if cnt == 0 {
		return
	}

	sess.PlayerData.Res.Fragments += cnt
	sess.SyncResource(structs.ResourceType_Fragment, sess.PlayerData.Res.Fragments)
}

func (sess *Session) StatueChange(cnt int32, changeType int32) {
	if cnt == 0 {
		return
	}

	sess.PlayerData.Res.Statue += cnt
	sess.SyncResource(structs.ResourceType_Fragment, sess.PlayerData.Res.Statue)
}

func (sess *Session) DetonatorChange(cnt int32, changeType int32) {
	if cnt == 0 {
		return
	}

	sess.PlayerData.Res.Detonator += cnt
	sess.SyncResource(structs.ResourceType_Detonator, sess.PlayerData.Res.Detonator)
}

func (sess *Session) MiningToolkitChange(cnt int32, changeType int32) {
	if cnt == 0 {
		return
	}

	sess.PlayerData.Res.MiningToolkit += cnt
	sess.SyncResource(structs.ResourceType_MiningToolkit, sess.PlayerData.Res.MiningToolkit)
}

func (sess *Session) StrengthChange(cnt int32, changeType int32) {
	if cnt == 0 {
		return
	}

	sess.PlayerData.Res.Strength += cnt

	ntf := &structs.SyncStrengthNtf{
		Strength: sess.PlayerData.Res.Strength,
	}
	sess.Send(structs.Protocol_SyncStrength_Ntf, ntf)
}
