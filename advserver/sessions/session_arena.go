package sessions

import (
	"adventure/advserver/gamedata"
	"adventure/common/structs"
	"adventure/common/util"
	"fmt"

	"time"
)

func (sess *Session) RefreshArenaData() {
	if sess.PlayerData.Arena.NextRefrshTime > time.Now().Unix() {
		return
	}

	hp := sess.PlayerData.HeroTeam.MaxHP()
	maxPlayerNum := gamedata.AllTemplates.GlobalData.MaxChallengePlayerNum

	groupLevel := gamedata.AllTemplates.GlobalData.RandomGroup1
	if hp < gamedata.AllTemplates.GlobalData.RandPlayerLimitHP {
		groupLevel = gamedata.AllTemplates.GlobalData.RandomGroup2
	}

	playerGroup := gamedata.ArenaRobtosGroup[0]
	for _, v := range gamedata.ArenaRobtosGroup {
		if v.MinHP < hp && hp < v.MaxHP {
			playerGroup = v
		}
	}

	targets := []*structs.FightTarget{}
	for i := int32(0); i < maxPlayerNum; i++ {
		findGroupID := playerGroup.ID + groupLevel[i]
		if findGroupID < 1 {
			findGroupID = 1
		}

		robotGroup := gamedata.ArenaRobtosGroup[0]
		for _, v := range gamedata.ArenaRobtosGroup {
			if v.ID == findGroupID {
				robotGroup = v
			}
		}

		fmt.Println("robotGroup : ", robotGroup.Players)

		robot := robotGroup.Players[util.RandNum(int32(0), int32(len(robotGroup.Players)))]

		targetPlayer := structs.FightTarget{
			PlayerID:   robot.ID,
			IconID:     robot.IconID,
			PlayerName: robot.Name,
			HP:         robot.HP,
		}

		targets = append(targets, &targetPlayer)
	}

	sess.PlayerData.Arena.Targets = targets
	sess.PlayerData.Arena.NextRefrshTime = time.Now().Unix()
	sess.PlayerData.Arena.ChallengeCount = 0
	sess.PlayerData.Arena.RewardRecord = []structs.ArenaRwdStatus{}
}

func (sess *Session) SyncPlayerArena() {
	lessTime := time.Unix(sess.PlayerData.Arena.NextRefrshTime, 0).Sub(time.Now()).Seconds()

	respSyncArena := &structs.SyncArenaNtf{
		Targets:         sess.PlayerData.Arena.Targets,
		LessRefreshTime: int32(lessTime),
		ChallengeCount:  sess.PlayerData.Arena.ChallengeCount,
	}

	respRwd := &structs.ArenaRecieveRewardResp{
		RewardRecords: sess.PlayerData.Arena.RewardRecord,
	}

	sess.Send(structs.Protocol_SyncArena_Ntf, respSyncArena)

	sess.Send(structs.Protocol_RecieveArenaReward_Resp, respRwd)
}
