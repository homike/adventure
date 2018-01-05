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
	case structs.RewardType_Property: //奖励资源类
		sess.ResourceChange(reward.Param1, reward.Param2*count, 0)

	case structs.RewardType_RandProperty:
		rID = reward.Param1
		rNum = (reward.Param2 + int32(util.RandNum(0, reward.Param3-reward.Param2))) * count
		sess.ResourceChange(rID, rNum, 0)

	case structs.RewardType_HP: //奖励基础战力
		rID = reward.Param1
		rNum = reward.Param2 * count
		sess.AddMainHeroHP(rNum, structs.AddHP_Type_WinBattle)

	case structs.RewardType_Item: //奖励道具
		rID = reward.Param1
		rNum = reward.Param2 * count
		item, isNew := sess.PlayerData.Bag.AddItem(rID, rNum)
		syncType := structs.SyncItem_Type_Update
		if isNew {
			syncType = structs.SyncItem_Type_Add
		}
		sess.SyncItems(uint8(syncType), []*structs.GameItem{item})

	case structs.RewardType_Exp: // 经验奖励
		rNum = reward.Param2 * count
		sess.CalculateGetExp(rNum)

	case structs.RewardType_UnlockGameLevel: // 解锁游戏关卡
		rID = reward.Param1
		gameLevel, err := sess.PlayerData.PlayerGameLevel.UnLockGameLevel(rID)
		logger.Debug("czx@@@ UnlockGameLevel Reward: %v, Error: %v", rID, err)
		if err == nil {
			sess.Send(structs.Protocol_OpenGameLevel_Ntf, &structs.OpenGameLevelNtf{GameLevel: gameLevel})
		}

	case structs.RewardType_Hero: // 英雄奖励
		rID = reward.Param1
		hero, err := sess.PlayerData.HeroTeam.AddHero("", false, rID)
		if err == nil {
			sess.AddHero([]*structs.Hero{hero}, true)
		}

	case structs.RewardType_UnlockMenu: // 解锁菜单
		rID = reward.Param1
		sess.UnLockMenu(rID)

	case structs.RewardType_AddHeroWorkTop: // 增加英雄出战数上限
		rNum = reward.Param1 * count
		err := sess.PlayerData.HeroTeam.AddWorkHeroTop(rNum)
		if err == nil {
			sess.SyncHeroWorkTop()
		}

	case structs.RewardType_UnlockArtifact: // 解锁神器
		rID = reward.Param1
		err := sess.PlayerData.Artifact.UnlockArtifact(rID)
		if err == nil {
			sess.SyncArtifactStatus(structs.SyncType_Update)
			//CZXDO: 全服公告
		}

	case structs.RewardType_AddGetGiftDayNum: // 增加好友中每日领取礼物次数
	case structs.RewardType_AddSendGiftDayNum: // 增加好友中每日送礼次数
	case structs.RewardType_TradeTaskReset: // 商会任务重置
	case structs.RewardType_AddGameBoxNumTop: // 增加宝箱上限
	//CZXDO: 挖矿奖励
	case structs.RewardType_AddMiningPickNumTop: // 增加挖掘次数上限
	case structs.RewardType_AddMiningPickLevel: // 增加矿镐等级
	case structs.RewardType_AddMiningPickNum: // 增加挖矿次数，无视上限
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
	resRewads := []*structs.Reward{}
	resContent := ""
	otherRewads := []*structs.Reward{}
	otherContent := []string{}
	for _, rewardID := range rewardIDs {
		rewardT, ok := gamedata.AllTemplates.RewardTemplates[rewardID]
		if !ok {
			logger.Error("can not find reward %v template", rewardID)
			continue
		}
		reward, err := sess.DoReward(&rewardT, 1)
		if err != nil {
			logger.Error("DoReward(%v) error %v", rewardID, err)
			continue
		}
		if rewardT.Type == structs.RewardType_Property || rewardT.Type == structs.RewardType_Item {
			resRewads = append(resRewads, reward)
			resContent = rewardT.Context
		} else {
			otherRewads = append(otherRewads, reward)
			otherContent = append(otherContent, rewardT.Context)
		}
	}

	if len(resRewads) > 0 {
		sess.RewardResults(true, resRewads, resContent)
	} else {
		for k, v := range otherRewads {
			sess.RewardResults(false, []*structs.Reward{v}, otherContent[k])
		}
	}
	return nil
}

func (sess *Session) RewardResults(isRes bool, rewards []*structs.Reward, context string) {
	ntf := &structs.RewardResultNtf{
		IsRes:   isRes,
		Rewards: rewards,
		Context: context,
	}

	sess.Send(structs.Protocol_RewardResult_Ntf, ntf)
}
