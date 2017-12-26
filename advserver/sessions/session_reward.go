package sessions

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"adventure/common/util"
	"errors"
)

func (sess *Session) DoReward(reward *structs.RewardTemplate, count int32) (*structs.Reward, error) {
	if reward == nil {
		return nil, errors.New("DoReward error, reward is nill")
	}

	rID := int32(0)
	rNum := int32(0)
	switch reward.Type {
	case structs.RewardType_Property:
		sess.ResourceChange(reward.Param1, reward.Param2*count, 0)
	case structs.RewardType_RandProperty:
		rID = reward.Param1
		rNum = (reward.Param2 + int32(util.RandNum(0, reward.Param3-reward.Param2))) * count
		sess.ResourceChange(rID, rNum, 0)
	case structs.RewardType_HP:
		rID = reward.Param1
		rNum := reward.Param2 * count
		sess.AddMainHeroHP(rNum, structs.AddHP_Type_WinBattle)
	case structs.RewardType_Item:
		rID = reward.Param1
		rNum := reward.Param2 * count
		item, isNew := sess.PlayerData.Bag.AddItem(rID, rNum)
		syncType := structs.SyncItem_Type_Update
		if isNew {
			syncType = structs.SyncItem_Type_Add
		}
		sess.SyncItems(uint8(syncType), []*structs.GameItem{item})
	//CZXDO: 奖励类型写完
	case structs.RewardType_Exp: // 经验奖励
	case structs.RewardType_UnlockGameLevel: // 解锁游戏关卡
	case structs.RewardType_Hero: // 英雄奖励
	case structs.RewardType_UnlockMenu: // 解锁菜单
	case structs.RewardType_AddHeroWorkTop: // 增加英雄出战数上限
	case structs.RewardType_AddMiningPickNumTop: // 增加挖掘次数上限
	case structs.RewardType_AddMiningPickLevel: // 增加矿镐等级
	case structs.RewardType_AddMiningPickNum: // 增加挖矿次数，无视上限
	case structs.RewardType_UnlockArtifact: // 解锁神器
	case structs.RewardType_AddGetGiftDayNum: // 增加好友中每日领取礼物次数
	case structs.RewardType_AddSendGiftDayNum: // 增加好友中每日送礼次数
	case structs.RewardType_TradeTaskReset: // 商会任务重置
	case structs.RewardType_AddGameBoxNumTop: // 增加宝箱上限
	}

	return &structs.Reward{RewardType: uint8(reward.Type), Param1: rID, Param2: rNum}, nil
}

func (sess *Session) DoSomeReward(itemTemplateID int32, num int32) error {
	//CZXDO: 动态掉落

	rewardIDs := gamedata.AllTemplates.ItemTemplates[itemTemplateID].RewardIDs
	_ = rewardIDs
	return nil
}

func (sess *Session) DoSomeRewards(rewardIDs []int32) error {
	//CZXDO: 动态掉落

	//rewardIDs := gamedata.AllTemplates.ItemTemplates[itemTemplateID].RewardIDs
	//_ = rewardIDs
	return nil
}

func (sess *Session) RewardResults(isRes bool, rewards []structs.Reward, context string) {
	ntf := &structs.RewardResultNtf{
		IsRes:   isRes,
		Rewards: rewards,
		Context: context,
	}

	sess.Send(structs.Protocol_RewardResult_Ntf, ntf)
}
