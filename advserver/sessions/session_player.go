package sessions

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"fmt"
)

func (sess *Session) OnEnterGame() {

	/****************同步基础数据********************/
	sess.SyncPlayerBaseInfo()
	sess.SyncStrength()        // 同步饱食度
	sess.SyncUserGuidRecords() // 同步新手引导进度
	//CZXDO: 同步客户端已食用过的食物列表
	//CZXDO: 同步购买商城物品记录
	sess.SyncUnlockMenus()      //同步客户端已解锁菜单列表
	sess.SyncGameBoxTopNumNtf() //同步客户端附加的宝箱上限数量

	/****************同步英雄数据********************/
	sess.SyncHeroWorkTop()                                                       // 同步最大出战英雄数
	sess.SyncHeroNtf(structs.SyncHeroType_First, sess.PlayerData.HeroTeam.Heros) // 同步英雄信息

	/****************同步背包数据********************/
	sess.SyncAllResources()   // 同步所有资源
	sess.SyncBag()            // 同步背包数据
	sess.SyncPlayerUsedItem() // 同步客户端已使用过的物品列表
}

// 1008
func (sess *Session) SyncPlayerBaseInfo() {
	fmt.Println("czx@@@ SyncPlayerBaseInfo")

	resp := &structs.SyncPlayerBaseInfoNtf{
		PlayerID:           1,
		GameZoonID:         1,
		IsSupperMan:        true,
		PlatformType:       1,
		Viplevel:           1,
		TotalRechargeIngot: 1,
	}
	sess.Send(structs.Protocol_SyncPlayerBaseInfo_Ntf, resp)
}

func (sess *Session) SyncHeroNtf(syncType uint8, heros []*structs.Hero) {
	fmt.Println("SyncHeroNtf heros num : ", len(heros))
	resp := &structs.SyncHeroNtf{
		SyncHeroType: syncType,
		Heros:        heros,
	}
	sess.Send(structs.Protocol_SyncHero_Ntf, resp)
}

func (sess *Session) SyncStrength() {
	resp := &structs.SyncStrengthNtf{
		Strength: sess.PlayerData.Res.Strength,
	}
	sess.Send(structs.Protocol_SyncStrength_Ntf, resp)
}

func (sess *Session) SyncHeroWorkTop() {
	resp := &structs.SyncHeroWorkTopNtf{
		MaxWorker: sess.PlayerData.HeroTeam.MaxWorker,
	}
	sess.Send(structs.Protocol_SyncWorkHeroTop_Ntf, resp)
}

func (sess *Session) SyncUnlockMenus() {
	resp := &structs.SyncUnlockMenusNtf{
		MenuStates: sess.PlayerData.MenuStates,
	}
	sess.Send(structs.Protocol_SyncUnlockMenus_Ntf, resp)
}

func (sess *Session) SyncGameBoxTopNumNtf() {
	resp := &structs.SyncGameBoxTopNumNtf{
		AddNum: sess.PlayerData.AddGameBoxCount,
	}
	sess.Send(structs.Protocol_SyncGameBoxTopNum_Ntf, resp)
}

func (sess *Session) SyncUserGuidRecords() {
	records := []structs.GuildRecord{}
	for i := 0; i < 24; i++ {
		records = append(records, structs.GuildRecord{
			UserGuidTypes: uint8(i),
			TriggerCount:  int32(5),
		})
	}
	resp := &structs.SyncUserGuidRecordsNtf{
		Records: records,
	}

	sess.Send(structs.Protocol_SyncUserGuidRecords_Ntf, resp)
}

func (sess *Session) DoSomeRewards(itemTemplateID int32, num int32) error {
	//CZXDO: 动态掉落
	rewardIDs, err := gamedata.AllTemplates.ItemTemplate.RewardIDs(itemTemplateID)
	if err != nil {
		return err
	}
	_ = rewardIDs
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
