package sessions

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"time"
)

func (sess *Session) CheckAchievements(condType structs.AchvCondType, condID, addCount int32) {

	sess.RefreshCircleAchievements()

	switch condType {
	case structs.AchvCondType_Collect: // 收集物品 / 资源
		sess.CheckCollect(condType, condID, addCount)

	case structs.AchvCondType_KillStatue: // 杀巨魔雕像
		sess.CheckCommonAchievement(condType, 0, 0, addCount)

	case structs.AchvCondType_KillLevelStatue: // 杀某等级巨魔雕像
		sess.CheckCommonAchievement(condType, 0, addCount, 0)

	case structs.AchvCondType_KillBoss: // 杀boss
		sess.CheckCommonAchievement(condType, condID, 0, addCount)

	case structs.AchvCondType_ChallengePlayer: // 挑战玩家
		sess.CheckCommonAchievement(condType, 0, 0, addCount)

	case structs.AchvCondType_WinArenaPlayer: // 战胜玩家
		sess.CheckCommonAchievement(condType, 0, 0, addCount)

	case structs.AchvCondType_CollectHero: // 收集英雄
		sess.CheckCommonAchievement(condType, condID, 0, addCount)

	case structs.AchvCondType_CollectPoint: // 收集点
		sess.CheckCommonAchievement(condType, condID, 0, addCount)

	case structs.AchvCondType_RecruitHeros: // 招募英雄
		sess.CheckCommonAchievement(condType, 0, 0, addCount)

	case structs.AchvCondType_RecruitHeroIngot: // 元宝招募英雄
		sess.CheckCommonAchievement(condType, 0, 0, addCount)

	case structs.AchvCondType_InvitationFriends: // 邀请好友
		sess.CheckCommonAchievement(condType, 0, 0, addCount)

	case structs.AchvCondType_FatalismWeapon: // 宿命武器
		sess.CheckCommonAchievement(condType, 0, addCount, 0)

	case structs.AchvCondType_OpenMenu: // 开启菜单
		sess.CheckCommonAchievement(condType, 0, addCount, 0)

	case structs.AchvCondType_PassGameLevel: // 通过游戏关卡
		sess.CheckCommonAchievement(condType, addCount, 0, 0)

	case structs.AchvCondType_PassRiftLevel: // 通过某个秘境
		sess.CheckCommonAchievement(condType, addCount, 0, 0)

	case structs.AchvCondType_MasterHeroLevel: // 主角英雄等级
		sess.CheckLevelAchievement(condType, addCount)

	case structs.AchvCondType_OpenGameLevel: // 开启关卡
		sess.CheckLevelAchievement(condType, addCount)
	}
}

// 刷新日成就 / 周成就
func (sess *Session) RefreshCircleAchievements() {

	usrAchv := sess.PlayerData.Achievement
	bChange := false
	// 每日成就刷新
	if time.Now().Unix() > usrAchv.NextRefreshTimeDaily {
		for _, achv := range usrAchv.Achievements {
			achvT, ok := gamedata.AllTemplates.AchievementTemplates[achv.TemplateID]
			if !ok {
				continue
			}
			if achvT.AchievementType == structs.AchvType_DayCircle && achv.Status != structs.AchvStatus_UnActive {
				achv.TotalCount = 0
				if achvT.PreID > 0 {
					achv.Status = structs.AchvStatus_UnActive
				} else {
					achv.Status = structs.AchvStatus_Active
				}
			}
		}
		bChange = true
		usrAchv.NextRefreshTimeDaily = time.Now().AddDate(0, 0, 1).Unix()
	}
	// 每周成就刷新
	if time.Now().Unix() > usrAchv.NextRefreshTimeWeekly {
		for _, achv := range usrAchv.Achievements {
			achvT, ok := gamedata.AllTemplates.AchievementTemplates[achv.TemplateID]
			if !ok {
				continue
			}
			if achvT.AchievementType == structs.AchvType_ManyDayCircle1 && achv.Status != structs.AchvStatus_UnActive {
				achv.TotalCount = 0
				if achvT.PreID > 0 {
					achv.Status = structs.AchvStatus_UnActive
				} else {
					achv.Status = structs.AchvStatus_Active
				}
			}
		}
		bChange = true
		usrAchv.NextRefreshTimeDaily = time.Now().AddDate(0, 0, 7-int(time.Now().Weekday())).Unix()
	}

	if bChange {
		ntf := &structs.GetAchievementsResp{
			Achievements:          usrAchv.GetAchieveMentsArray(),
			NextRefreshTimeDaily:  usrAchv.NextRefreshTimeDaily,
			NextRefreshTimeWeekly: usrAchv.NextRefreshTimeWeekly,
		}
		sess.Send(structs.Protocol_GetAchievements_Resp, ntf)
	}
}

// 检测收集成就
func (sess *Session) CheckCollect(condType structs.AchvCondType, condID, addCount int32) {
	arrAchv, arrAchvT := sess.PlayerData.Achievement.GetAchievements(condType, condID, 0)
	if len(arrAchv) <= 0 || len(arrAchv) != len(arrAchvT) {
		return
	}

	for i := 0; i < len(arrAchv); i++ {
		if arrAchv[i].TotalCount > arrAchvT[i].ConditionCount {
			arrAchv[i].Status = structs.AchvStatus_Finish
			if !arrAchvT[i].IsConditinCountAddup() {
				arrAchv[i].TotalCount = arrAchvT[i].ConditionCount
			}
			//CZXDO: 如果成就关联了活动，则为玩家开启活动记录
		}
	}

	ntf := structs.UpdateAchievementNtf{
		Achievements: arrAchv,
	}
	sess.Send(structs.Protocol_UpdateAchievement_Ntf, ntf)
}

// 击杀巨魔雕像成就
// 检测杀boss成就
// 检测杀boss成就
// 检测挑战玩家成就
// 检测战胜玩家的成就
// 检测普通收集英雄成就
// CZXDO: 坚持收集点成就
// 检测主角英雄等级成就
// 坚持累计普通／元宝招募英雄次数成就
func (sess *Session) CheckCommonAchievement(condType structs.AchvCondType, condID, condCount, addCount int32) {
	arrAchv, arrAchvT := sess.PlayerData.Achievement.GetAchievements(condType, condID, condCount)
	if len(arrAchv) <= 0 || len(arrAchv) != len(arrAchvT) {
		return
	}

	for i := 0; i < len(arrAchv); i++ {
		arrAchv[i].TotalCount += addCount
		if arrAchv[i].TotalCount >= arrAchvT[i].ConditionCount {
			arrAchv[i].Status = structs.AchvStatus_Finish

			if !arrAchvT[i].IsConditinCountAddup() {
				arrAchv[i].TotalCount = arrAchvT[i].ConditionCount
			}
		}
	}
	ntf := structs.UpdateAchievementNtf{
		Achievements: arrAchv,
	}
	sess.Send(structs.Protocol_UpdateAchievement_Ntf, ntf)
}

// 检测杀某等级巨魔雕像成就
// func (sess *Session) CheckKillLevelStatue(condType structs.AchvCondType, level int32) {
// 	arrAchv, arrAchvT := sess.PlayerData.Achievement.GetAchievements(condType, 0, level)
// 	if len(arrAchv) <= 0 || len(arrAchv) != len(arrAchvT) {
// 		return
// 	}

// 	for i := 0; i < len(arrAchv); i++ {
// 		arrAchv[i].Status = structs.AchvStatus_Finish
// 	}
// 	ntf := structs.UpdateAchievementNtf{
// 		Achievements: arrAchv,
// 	}
// 	sess.Send(structs.Protocol_UpdateAchievement_Ntf, ntf)
// }

// 检测开启关卡等级成就
func (sess *Session) CheckLevelAchievement(condType structs.AchvCondType, level int32) {
	arrAchv, arrAchvT := sess.PlayerData.Achievement.GetAchievements(condType, 0, 0)
	if len(arrAchv) <= 0 || len(arrAchv) != len(arrAchvT) {
		return
	}

	for i := 0; i < len(arrAchv); i++ {
		if arrAchvT[i].ConditionCount <= level {
			arrAchv[i].Status = structs.AchvStatus_Finish
		} else {
			arrAchv[i].TotalCount = level
		}
	}
	ntf := structs.UpdateAchievementNtf{
		Achievements: arrAchv,
	}
	sess.Send(structs.Protocol_UpdateAchievement_Ntf, ntf)
}

//
