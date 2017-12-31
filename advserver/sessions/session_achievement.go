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
	case structs.AchvCondType_KillStatue: // 杀巨魔雕像
	case structs.AchvCondType_KillLevelStatue: // 杀某等级巨魔雕像
	case structs.AchvCondType_KillBoss: // 杀boss
	case structs.AchvCondType_ChallengePlayer: // 挑战玩家
	case structs.AchvCondType_WinArenaPlayer: // 战胜玩家
	case structs.AchvCondType_CollectHero: // 收集英雄
	case structs.AchvCondType_CollectPoint: // 收集点
	case structs.AchvCondType_MasterHeroLevel: // 主角英雄等级
	case structs.AchvCondType_RecruitHeros: // 招募英雄
	case structs.AchvCondType_RecruitHeroIngot: // 元宝招募英雄
	case structs.AchvCondType_InvitationFriends: // 邀请好友
	case structs.AchvCondType_OpenGameLevel: // 开启关卡
	case structs.AchvCondType_FatalismWeapon: // 宿命武器
	case structs.AchvCondType_OpenMenu: // 开启菜单
	case structs.AchvCondType_PassGameLevel: // 通过游戏关卡
	case structs.AchvCondType_PassRiftLevel: // 通过某个秘境
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

func (sess *Session) CheckCollect(condType structs.AchvCondType, condID, addCount int32) {
	arrAchv, arrAchvT := sess.PlayerData.Achievement.GetAchievements(condType, condID)
	if len(arrAchv) <= 0 || len(arrAchv) != len(arrAchvT) {
		return
	}

	for i := 0; i < len(arrAchv); i++ {
		if arrAchv[i].TotalCount > arrAchvT[i].ConditionCount {
			arrAchv[i].Status = structs.AchvStatus_Finish
			if arrAchvT[i].
		}
	}
}
