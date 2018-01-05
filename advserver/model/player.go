package model

import (
	"adventure/advserver/db/mysql"
	"adventure/advserver/db/redis"
	"adventure/advserver/gamedata"
	"adventure/advserver/log"
	"adventure/common/clog"
	"adventure/common/structs"
	"fmt"
	"time"
)

var (
	logger    *clog.Logger
	playerDao *mysql.PlayerDao
)

func Init() error {
	logger = log.GetLogger()
	playerDao = mysql.NewUserDao()

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
	VipLevel          int32
	OnlineTime        int
	NextFreeIngotTime int64                     // 下一次免费招募时间
	LastEmployTime    int64                     // 最近一次刷新的招募时间
	AddGameBoxCount   int32                     // 增加的宝箱上限数量
	HeroTeam          *HeroTeams                // 玩家英雄
	Res               *Resource                 // 玩家资源
	PlayerGameLevel   *PlayerGameLevel          // 关卡数据
	Bag               *Bag                      // 背包
	UserGuidRecords   []*structs.GuildRecord    // 新手引导记录
	MenuStates        []*structs.MenuStatusItem // 菜单状态
	Artifact          *Artifact                 // 神器
	Achievement       *PlayerAchievenment       // 成就
	ExtendData        *ExtendData               // 扩展数据
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
	player.VipLevel = 10

	// 初始化玩家英雄
	player.HeroTeam = NewHeroTeams()
	hero, _ := player.HeroTeam.AddHero(player.Name, true, heroTemplateID)
	hero1, _ := player.HeroTeam.AddHero(player.Name, false, 10001)
	hero1.IsOutFight = true

	player.HeroTeam.ReCalculateHeroLevelHp(hero)
	player.HeroTeam.ReCalculateHeroLevelHp(hero1)

	// 初始化玩家资源
	player.Res = NewResource()
	for i := 101; i <= 115; i++ {
		player.Res.Ores.Add(int32(i), 10)
	}
	for i := 200; i <= 208; i++ {
		player.Res.Foods.Add(int32(i), 10)
	}

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
		CompleteEvent: make([]structs.AdventureEventStatus, len(levelT.EvnetIDs)),
	}
	player.PlayerGameLevel.AddGameLevel(&gameLevel)

	// 背包数据初始化
	player.Bag = NewBag()

	// 新手引导状态初始化
	player.UserGuidRecords = make([]*structs.GuildRecord, 0, 10)

	// 神器初始化
	player.Artifact = NewArtifact()

	// 成就初始化
	player.Achievement = NewPlayerAchievenment()

	// 菜单初始化
	for i := structs.MenuTypes_Temple; i <= structs.MenuTypes_Rift; i++ {
		player.MenuStates = append(player.MenuStates, &structs.MenuStatusItem{
			MenuID:     int32(i),
			MenuStatus: structs.MenuStatus_New,
		})
	}

	dbData := &mysql.PlayerDB{
		AccountID: player.AccountID,
		Name:      player.Name,
	}
	err = playerDao.CreatePlayer(dbData)
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

func (p *Player) GetFightTeamByHeroTemplateIDs(heroTemplateIDs []int32) *structs.FightTeam {
	team := &structs.FightTeam{}
	fightHeroTemplateIDs := []int32{}
	spellIDs := []int32{}

	doCombinationSpellIDs := make(map[int32]struct{}) // 已经处理过的合作技
	for _, heroTemplateID := range heroTemplateIDs {
		heroTemplate, ok := gamedata.AllTemplates.HeroTemplates[heroTemplateID]
		if !ok {
			continue
		}
		fightHeroTemplateIDs = append(fightHeroTemplateIDs, heroTemplateID)

		spell := &structs.SpellTemplate{}
		for _, spellID := range heroTemplate.SkillID {
			spellT, ok := gamedata.AllTemplates.SpellTemplates[spellID]
			fmt.Println("npc spellID: ", spellID, "spell ", spellT)
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
				for _, cheroID := range heroTemplateIDs {
					if cheroID == checkHero {
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
