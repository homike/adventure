package sessions

import (
	"adventure/advserver/gamedata"
	"adventure/advserver/model"
	"adventure/common/structs"
	"errors"
	"fmt"
	"math"
	"time"
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
	sess.SyncArtifactStatus(structs.SyncType_First)
	sess.SyncEmployInfo()

	/****************同步背包数据********************/
	sess.SyncAllResources()   // 同步所有资源
	sess.SyncBag()            // 同步背包数据
	sess.SyncPlayerUsedItem() // 同步客户端已使用过的物品列表

	/****************同步冒险关卡********************/
	sess.SyncGameLevelNtf()
	sess.SyncCurrentGameLevelNtf()
}

func (sess *Session) SyncPlayerBaseInfo() {
	fmt.Println("czx@@@ SyncPlayerBaseInfo")

	resp := &structs.SyncPlayerBaseInfoNtf{
		PlayerID:           int32(sess.PlayerData.AccountID),
		GameZoonID:         1,
		IsSupperMan:        true,
		PlatformType:       1,
		Viplevel:           sess.PlayerData.VipLevel,
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
	//CZXDO: 新手引导功能
	records := []*structs.GuildRecord{}
	for i := 0; i < 24; i++ {
		records = append(records, &structs.GuildRecord{
			UserGuidTypes: uint8(i),
			TriggerCount:  int32(5),
		})
	}
	resp := &structs.SyncUserGuidRecordsNtf{
		Records: records, //sess.PlayerData.UserGuidRecords,
	}

	sess.Send(structs.Protocol_SyncUserGuidRecords_Ntf, resp)
}

func (sess *Session) HeroHPAdd(addType uint8, heroID, addHP int32) {

	resp := &structs.HeroHPAddNtf{
		Type:   addType,
		HeroID: heroID,
		AddHP:  addHP,
	}

	sess.Send(structs.Protocol_HeroHpAdd_Ntf, resp)
}

func (sess *Session) AddMainHeroHP(num int32, addType uint8) {
	hero, err := sess.PlayerData.HeroTeam.GetMainHero()
	if err != nil {
		logger.Error("AddMainHeroHP GetMainHero() Error %v", err)
		return
	}

	oldHP := hero.HP()

	hero.ItemHP += num
	err = sess.PlayerData.HeroTeam.ReCalculateHeroLevelHp(hero)
	if err != nil {
		logger.Error("AddMainHeroHP ReCalculateHeroLevelHp() Error %v", err)
		return
	}
	// 通知英雄变化
	sess.SyncHeroNtf(structs.SyncHeroType_Update, []*structs.Hero{hero})
	// 通知战力变化
	sess.HeroHPAdd(addType, hero.HeroID, (hero.HP() - oldHP))
}

func (sess *Session) SyncArtifactStatus(stype structs.SyncType) {
	resp := &structs.SyncArtifactStatusNtf{
		SType:  stype,
		Status: sess.PlayerData.Artifact.Status,
	}

	sess.Send(structs.Protocol_SyncArtifactStatus_Ntf, resp)
}

func (sess *Session) AddHero(heros []*structs.Hero, bNotify bool) {

	sess.SyncHeroNtf(structs.SyncHeroType_Add, heros)

	for _, hero := range heros {
		sess.CheckAchievements(structs.AchvCondType_CollectHero, int32(hero.Quality), 1)

		if bNotify && (hero.Quality == structs.QualityType_Gold || hero.Quality == structs.QualityType_SplashGold) {
			//CZXDO: 获得英雄公告
		}
	}
}

// 更新玩家经验
func (sess *Session) CalculateGetExp(exp int32) []int32 {
	p := sess.PlayerData
	if exp < 1 {
		return nil
	}

	updateLevelHeroID := []int32{}

	fightHeros := p.HeroTeam.GetFightHeros()
	if len(fightHeros) <= 0 {
		return nil
	}

	avgExp := exp / int32(len(fightHeros))
	if avgExp < 1 {
		avgExp = 1
	}

	for _, hero := range fightHeros {
		hero.Exp += avgExp
		if hero.Exp > math.MaxInt32 {
			hero.Exp = math.MaxInt32
		}

		if int(hero.Level) <= len(gamedata.AllTemplates.HeroLevelTemplates) {
			nextLevelExp, err := gamedata.GetHeroLevelExp(hero.Level, hero.AwakeCount)
			if err != nil {
				logger.Error("GetHeroLevelExp(%v, %v) error", hero.Level, hero.AwakeCount)
				continue
			}
			heroHasUpdate := false
			for hero.Exp >= nextLevelExp {
				if int(hero.Level) <= len(gamedata.AllTemplates.HeroLevelTemplates) {
					hero.Level++
					nextLevelExp, _ = gamedata.GetHeroLevelExp(hero.Level-1, hero.AwakeCount-1)
					p.HeroTeam.ReCalculateHeroLevelHp(hero)
					heroHasUpdate = true

					if hero.IsPlayer {
						sess.CheckAchievements(structs.AchvCondType_MasterHeroLevel, 0, hero.Level)
					}
				} else {
					break
				}
			}

			// 通知客户端，英雄升级
			if heroHasUpdate {
				updateLevelHeroID = append(updateLevelHeroID, hero.HeroID)
				sess.SyncHeroNtf(structs.SyncHeroType_Update, []*structs.Hero{hero})
			}
		}
	}

	return updateLevelHeroID
}

func (sess *Session) RefreshPlayerInfo(reward *structs.OfflineReward) {
	p := sess.PlayerData
	//CZXDO: 刷新挖矿地图

	// 刷新关卡地图
	gameLevelT, ok := gamedata.AllTemplates.GameLevelTemplates[p.PlayerGameLevel.CurrentGameLevelID]
	if !ok {
		return
	}

	lrefreshTime := time.Unix(p.PlayerGameLevel.LastRefreshTime, 0)
	sec := int32(time.Now().Sub(lrefreshTime).Seconds())
	if sec < 1 {
		return
	}

	if reward != nil {
		reward.OfflineTimeSec = sec
	}

	p.PlayerGameLevel.LastRefreshTime = time.Now().Unix()

	if p.HeroTeam.MaxHP() < gameLevelT.MinHP {
		return
	}

	fullSec := int32(0)
	halfSec := int32(0)

	if p.Res.Strength > sec {
		fullSec = sec
		p.Res.Strength -= fullSec
	} else {
		fullSec = p.Res.Strength
		halfSec = sec - p.Res.Strength
		p.Res.Strength = 0
	}
	if reward != nil {
		reward.HasStrength = p.Res.Strength > 0
	}

	// 事件的进度处理
	curGameLevel, err := p.PlayerGameLevel.GetCurGameLevelData()
	if err != nil {
		return
	}
	unActiveCnt := p.PlayerGameLevel.GetUnActiveEventCount()
	if unActiveCnt == 0 {
		curGameLevel.EventProgress = 0
	} else {
		curGameLevel.EventProgress += sec
		changeStatus := false
		for k, v := range curGameLevel.CompleteEvent {
			if v == structs.AdventureEventStatus_UnActive {
				eventT, ok := gamedata.AllTemplates.GameLevelEventTemplates[gameLevelT.EvnetIDs[k]]
				if !ok {
					continue
				}
				if curGameLevel.EventProgress > eventT.ActiveEventSec {
					curGameLevel.EventProgress -= eventT.ActiveEventSec
					curGameLevel.CompleteEvent[k] = structs.AdventureEventStatus_Active
				}
			}
		}
		if changeStatus && p.PlayerGameLevel.GetUnActiveEventCount() == 0 {
			curGameLevel.EventProgress = 0
		}
	}

	// 刷新宝箱进度
	if len(gameLevelT.GameBoxIDs) > 0 {
		//CZXDO: 宝箱最大数量
		if curGameLevel.BoxCount < 99 {
			curGameLevel.GameBoxProgress += sec
			for curGameLevel.GameBoxProgress > gameLevelT.ActiveGameBoxSec {
				curGameLevel.GameBoxProgress -= gameLevelT.ActiveGameBoxSec
				curGameLevel.BoxCount++
				if curGameLevel.BoxCount >= 99 {
					curGameLevel.GameBoxProgress = 0
					break
				}
			}
		}
	} else {
		curGameLevel.BoxCount = 0
	}

	// 刷新金钱
	money := int32(0)
	if fullSec > 0 {
		money += gameLevelT.MoneyPer * fullSec
	}
	if halfSec > 0 {
		money += gameLevelT.MoneyPer * halfSec / 2
	}
	p.Res.Money += money

	// 刷新英雄经验等级
	exp := int32(0)
	if fullSec > 0 {
		exp += gameLevelT.ExpPer * fullSec
	}
	if halfSec > 0 {
		exp += gameLevelT.ExpPer * halfSec / 2
	}

	if reward != nil {
		reward.Exp = exp
		reward.OfflineHP = p.HeroTeam.MaxHP()
		updateLevelHeroId := sess.CalculateGetExp(exp)
		reward.OnlineHP = p.HeroTeam.MaxHP()
		reward.UpLevelHero = updateLevelHeroId
	} else {
		sess.CalculateGetExp(exp)
	}
}

// 人机战斗模拟
func (sess *Session) DoFightTest(battleFieldID int32) (*structs.FightResult, error) {
	fmt.Println("battleFieldID ", battleFieldID)

	p := sess.PlayerData

	battleFieldT, ok := gamedata.AllTemplates.Battlefields[battleFieldID]
	if !ok {
		return nil, errors.New("battleFieldID cannot find battle field")
	}

	playerTeam := p.GetFightTeamByHeroTeam()

	fmt.Println("playerTeam ", playerTeam.DefaultHP, ",spellIDs ", playerTeam.SpellIDs)

	btTeam := model.GetFightTeamByHeroTemplateIDs(battleFieldT.NpcIDs)
	btTeam.DefaultHP = battleFieldT.HP
	btTeam.ShanBi = battleFieldT.ShanBi
	btTeam.XianGong = battleFieldT.XianGong
	btTeam.FangYu = battleFieldT.FangYu
	btTeam.WangZhe = battleFieldT.WangZhe
	btTeam.Name = battleFieldT.Name

	fmt.Println("btTeam ", btTeam.DefaultHP, ", spellIDs ", btTeam.SpellIDs)

	sim := model.NewFightSim(playerTeam, btTeam)
	fightRet := sim.Fight()
	fightRet.BackgroundID = battleFieldT.BackgroundID
	fightRet.ForegroundID = battleFieldT.ForegroundID

	return fightRet, nil
}

func (sess *Session) UnLockMenu(menuID int32) {
	menu := new(structs.MenuStatusItem)
	for _, v := range sess.PlayerData.MenuStates {
		if v.MenuID == menuID {
			menu = v
		}
	}

	if menu.MenuID == 0 {
		sess.PlayerData.MenuStates = append(sess.PlayerData.MenuStates, &structs.MenuStatusItem{
			MenuID:     menuID,
			MenuStatus: structs.MenuStatus_New,
		})
		sess.CheckAchievements(structs.AchvCondType_OpenMenu, 0, menuID)
	}
	menu.MenuStatus = structs.MenuStatus_New

	ntf := &structs.SyncUnlockMenuNtf{
		MenuID: menu.MenuID,
	}
	sess.Send(structs.Protocol_UnLockMenu_Ntf, ntf)

	sess.onUnlockMenu(structs.MenuTypes(menu.MenuID))
}

func (sess *Session) onUnlockMenu(mType structs.MenuTypes) {
	switch mType {
	case structs.MenuTypes_FS: // 封神之阶
	case structs.MenuTypes_Recruit: // 招募
		sess.PlayerData.HeroTeam.NextFreeIngotTime = time.Now().Add(time.Duration(gamedata.FreeIngotEmployFirstTimeSpan) * time.Second).Unix()
	case structs.MenuTypes_Rift: // 秘境
	case structs.MenuTypes_TradeHouse: // 商行
	case structs.MenuTypes_TradeTroop: // 贸易队
	}

}

func (sess *Session) SyncEatFoodList() {
	ids := []int32{}
	dates := []int64{}
	for k, v := range sess.PlayerData.ExtendData.EatedFoodRecords {
		ids = append(ids, k)
		dates = append(dates, v)
	}
	resp := &structs.GetEatedFoodsResp{
		FoodIDs:   ids,
		EatedDate: dates,
	}
	sess.Send(structs.Protocol_GetEatedFoods_Resp, resp)
}

func (sess *Session) SyncEmployInfo() {
	leftTime := time.Now().AddDate(0, 0, 1).Add(time.Duration(gamedata.SystemRefreshTime) * time.Second).Sub(time.Now()).Seconds()

	nextTime := int32(0)
	if sess.PlayerData.HeroTeam.NextFreeIngotTime > time.Now().Unix() {
		nextTime = int32(time.Unix(sess.PlayerData.HeroTeam.NextFreeIngotTime, 0).Sub(time.Now()).Seconds()) + 1
	}

	costList := make([]int32, 6, 6)
	typeList := []int32{}
	for i := structs.EmployType_Money; i <= structs.EmployType_ManyDiamond; i++ {
		typeList = append(typeList, int32(i))
		if i == structs.EmployType_Diamond || i == structs.EmployType_ManyDiamond {
			costList[i] = gamedata.InitialEmployCost[i]
		} else {
			costList[i] = int32(float64(gamedata.InitialEmployCost[i]) * math.Pow(2, float64(sess.PlayerData.HeroTeam.EmployRecord[i])))
		}
	}

	resp := &structs.SyncEmployResp{
		Type:                       typeList,
		Cost:                       costList,
		LeftSecond:                 int32(leftTime),
		NextFreeIngotEmployLeftSec: nextTime,
	}

	sess.Send(structs.Protocol_SyncEmploy_Resq, resp)
}

func (sess *Session) SyncTemplateHeros() {
	userTemple := sess.PlayerData.Temple

	userTemple.RefreshTemple()

	t1 := time.Unix(userTemple.NextRefreshTime, 0)
	leftTime := int32(time.Now().Sub(t1).Seconds())

	resp := &structs.SyncTemplateHerosResp{
		Heros:        userTemple.TempleHeros,
		LeftSecond:   leftTime,
		TradeCount:   userTemple.ToDayTradeCount,
		RefreshCount: userTemple.RefreshCount,
	}
	sess.Send(structs.Protocol_SyncTemplateHeros_Req, resp)
}
