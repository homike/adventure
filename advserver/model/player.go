package model

import (
	"adventure/advserver/db/mysql"
	"adventure/advserver/db/redis"
	"adventure/advserver/gamedata"
	"adventure/advserver/log"
	"adventure/advserver/service"
	"adventure/common/clog"
	"adventure/common/structs"
	"errors"
	"fmt"
	"time"
)

var logger *clog.Logger

func Init() error {
	logger = log.GetLogger()

	return nil
}

type Player struct {
	AccountID         uint
	Name              string
	PlatformAccountID int
	PlatformTypes     int
	GameZoneID        int
	CreateTime        time.Time
	LastLoginTime     time.Time
	LastLogoffTime    time.Time
	BarrageSet        string
	VipLevel          int
	OnlineTime        int
	HeroTeam          *HeroTeams                // 玩家英雄
	Res               *Resource                 // 玩家资源
	PlayerGameLevel   *PlayerGameLevel          // 关卡数据
	Bag               *Bag                      // 背包
	UserGuidRecords   []*structs.GuildRecord    // 新手引导记录
	MenuStates        []*structs.MenuStatusItem // 菜单状态
	Artifact          *Artifact                 // 神器
	AddGameBoxCount   int32                     // 增加的宝箱上限数量
	ExtendData        *ExtendData
	MiningMap         string
}

func InitPlayer() *Player {
	player := &Player{}
	player.AccountID = 0
	player.ExtendData = NewExtendData()

	return player
}

func NewPlayer(name string, heroTemplateID int32) (*Player, error) {

	playID, err := redis.GetIncrPlayerID()
	if err != nil {
		fmt.Println("incr player id error :", err)
		return nil, err
	}
	fmt.Println("create player id : ", playID)

	player := InitPlayer()
	player.AccountID = playID
	player.Name = name
	// 初始化玩家英雄
	player.HeroTeam = NewHeroTeams()
	hero, _ := player.HeroTeam.AddHero(player.Name, true, heroTemplateID)
	player.HeroTeam.ReCalculateHeroLevelHp(hero)
	// 初始化玩家资源
	player.Res = NewResource()

	// 关卡数据初始化
	levelT, ok := gamedata.AllTemplates.GameLevelTemplates[1]
	if !ok {
		fmt.Println("GetGameLevelEvents(1) error")
		return nil, err
	}

	player.PlayerGameLevel = NewPlayerGameLevel()
	gameLevel := structs.GameLevel{
		GameLevelID:   1,
		IsUnlock:      true,
		CompleteEvent: make([]uint8, len(levelT.EvnetIDs)),
	}
	player.PlayerGameLevel.AddGameLevel(&gameLevel)

	// 背包数据初始化
	player.Bag = NewBag()

	// 新手引导状态初始化
	player.UserGuidRecords = make([]*structs.GuildRecord, 0, 10)

	// 神器初始化
	player.Artifact = NewArtifact()

	dbData := &mysql.PlayerDB{
		AccountID: player.AccountID,
		Name:      player.Name,
	}
	err = service.PlayerDao.CreatePlayer(dbData)
	if err != nil {
		fmt.Println("NewPlayer() error %v", err)
		return nil, err
	}

	return player, err
}

func (p *Player) UpdateGuidRecords(guidType uint8) {
	for k, v := range p.UserGuidRecords {
		if v.UserGuidTypes == guidType {
			p.UserGuidRecords[k].TriggerCount++
			return
		}
	}
	p.UserGuidRecords = append(p.UserGuidRecords, &structs.GuildRecord{
		UserGuidTypes: guidType,
		TriggerCount:  1,
	})
}

func (p *Player) GetFightTeamByHeroTeam() *structs.FightTeam {
	team := &structs.FightTeam{}
	heroTemplateIDs := []int32{}
	spellIDs := []int32{}

	fightHeros := p.HeroTeam.GetFightHeros()

	doCombinationSpellIDs := make(map[int32]struct{}) // 已经处理过的合作技
	for _, hero := range fightHeros {
		heroTemplate, ok := gamedata.AllTemplates.HeroTemplates[hero.HeroTemplateID]
		if !ok {
			continue
		}

		heroTemplateIDs = append(heroTemplateIDs, hero.HeroTemplateID)
		spell := &structs.SpellTemplate{}
		if hero.IsPlayer {
			useArtifact := p.Artifact.GetArtifactStatusUse()
			if useArtifact != nil {
				artfcatT, ok := gamedata.AllTemplates.ArtifactTemplates[useArtifact.ArtifactID]
				if ok {
					spellT, _ := gamedata.AllTemplates.SpellTemplates[artfcatT.SpellID]
					spell = &spellT
				}
			}
		} else {
			if len(heroTemplate.SkillID) > 0 {
				spellT, _ := gamedata.AllTemplates.SpellTemplates[heroTemplate.SkillID[0]]
				spell = &spellT
			}
		}

		// 技能
		if spell != nil {
			team.ShanBi += spell.DodgeProp
			team.XianGong += spell.FirstProp
			team.FangYu += spell.DefenceProp
			team.WangZhe += spell.KingProp

			if spell.AttackType != structs.AttackEffectType_None {
				spellIDs = append(spellIDs, spell.ID)
			}
		}

		// 武器进阶
		if hero.WeaponAdvanceLevel == 1 {
			team.ShanBi += heroTemplate.WeaponAdvance_ShanBi
			team.XianGong += heroTemplate.WeaponAdvance_XianGong
			team.FangYu += heroTemplate.WeaponAdvance_FangYu
			team.WangZhe += heroTemplate.WeaponAdvance_WangZhe
		}

		// 合作技
		if heroTemplate.CombinationSpllID > 0 {
			_, ok := doCombinationSpellIDs[heroTemplate.CombinationSpllID]
			if ok {
				continue
			}

			success := true
			cspell := gamedata.AllTemplates.CombinationSpells[heroTemplate.CombinationSpllID]
			for i := 0; i < len(cspell.HeroList); i++ {
				checkHero := cspell.HeroList[i]
				checkNum := cspell.HeroNumList[i]

				heroCnt := int32(0)
				for _, chero := range fightHeros {
					if chero.HeroTemplateID == checkHero {
						heroCnt++
					}
				}
				if heroCnt < checkNum {
					success = false
					break
				}
			}

			if success {
				spellT, _ := gamedata.AllTemplates.SpellTemplates[cspell.SpellId]
				spell = &spellT
				if spell != nil {
					team.ShanBi += spell.DodgeProp
					team.XianGong += spell.FirstProp
					team.FangYu += spell.DefenceProp
					team.WangZhe += spell.KingProp
					if spell.AttackType != structs.AttackEffectType_None {
						spellIDs = append(spellIDs, spell.ID)
					}
				}
			}

			doCombinationSpellIDs[heroTemplate.CombinationSpllID] = struct{}{}
		}
	}

	team.Models = heroTemplateIDs
	team.SpellIDs = spellIDs

	team.DefaultHP = p.HeroTeam.MaxHP()
	playerHero := p.HeroTeam.GetPlayerHero()
	if playerHero != nil {
		team.Name = playerHero.Name
	}

	return team
}

func (p *Player) GetFightTeamByHeroTemplaetIds(heroTemplateIDs []int32) *structs.FightTeam {
	team := &structs.FightTeam{}
	fightHeroTemplateIDs := []int32{}
	spellIDs := []int32{}

	fightHeros := p.HeroTeam.GetFightHeros()

	doCombinationSpellIDs := make(map[int32]struct{}) // 已经处理过的合作技
	for _, hero := range fightHeros {
		heroTemplate, ok := gamedata.AllTemplates.HeroTemplates[hero.HeroTemplateID]
		if !ok {
			continue
		}
		fightHeroTemplateIDs = append(heroTemplateIDs, hero.HeroTemplateID)

		spell := &structs.SpellTemplate{}
		for _, spellID := range heroTemplate.SkillID {
			spellT, ok := gamedata.AllTemplates.SpellTemplates[spellID]
			if ok {
				spell = &spellT

				team.ShanBi += spell.DodgeProp
				team.XianGong += spell.FirstProp
				team.FangYu += spell.DefenceProp
				team.WangZhe += spell.KingProp

				if spell.AttackType != structs.AttackEffectType_None {
					spellIDs = append(spellIDs, spell.ID)
				}
			}
		}

		// 合作技
		if heroTemplate.CombinationSpllID > 0 {
			_, ok := doCombinationSpellIDs[heroTemplate.CombinationSpllID]
			if ok {
				continue
			}

			success := true
			cspell := gamedata.AllTemplates.CombinationSpells[heroTemplate.CombinationSpllID]
			for i := 0; i < len(cspell.HeroList); i++ {
				checkHero := cspell.HeroList[i]
				checkNum := cspell.HeroNumList[i]

				heroCnt := int32(0)
				for _, chero := range fightHeros {
					if chero.HeroTemplateID == checkHero {
						heroCnt++
					}
				}
				if heroCnt < checkNum {
					success = false
					break
				}
			}

			if success {
				spellT, _ := gamedata.AllTemplates.SpellTemplates[cspell.SpellId]
				spell = &spellT
				if spell != nil {
					team.ShanBi += spell.DodgeProp
					team.XianGong += spell.FirstProp
					team.FangYu += spell.DefenceProp
					team.WangZhe += spell.KingProp
					if spell.AttackType != structs.AttackEffectType_None {
						spellIDs = append(spellIDs, spell.ID)
					}
				}
			}

			doCombinationSpellIDs[heroTemplate.CombinationSpllID] = struct{}{}
		}
	}

	team.Models = fightHeroTemplateIDs
	team.SpellIDs = spellIDs

	return team
}

// 人机战斗模拟
func (p *Player) DoFightTest(battleFieldID int32) (*structs.FightResult, error) {
	battleFieldT, ok := gamedata.AllTemplates.Battlefields[battleFieldID]
	if !ok {
		return nil, errors.New("battleFieldID cannot find battle field")
	}

	playerTeam := p.GetFightTeamByHeroTeam()

	btTeam := p.GetFightTeamByHeroTemplaetIds(battleFieldT.NpcIDs)
	btTeam.DefaultHP = battleFieldT.HP
	btTeam.ShanBi = battleFieldT.ShanBi
	btTeam.XianGong = battleFieldT.XianGong
	btTeam.FangYu = battleFieldT.FangYu
	btTeam.WangZhe = battleFieldT.WangZhe
	btTeam.Name = battleFieldT.Name

	sim := NewFightSim(playerTeam, btTeam)
	fightRet := sim.Fight()
	fightRet.BackgroundID = battleFieldT.BackgroundID
	fightRet.ForegroundID = battleFieldT.ForegroundID

	return fightRet, nil
}
