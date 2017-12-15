package model

import (
	"adventure/advserver/db/mysql"
	"adventure/advserver/db/redis"
	"adventure/advserver/gamedata"
	"adventure/advserver/log"
	"adventure/advserver/service"
	"adventure/common/clog"
	"adventure/common/structs"
	"fmt"
	"time"
)

var logger *clog.Logger

func Init() {
	logger = log.GetLogger()
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
	AddGameBoxCount   int32                     // 增加的宝箱上限数量
	MiningMap         string
	ExtendData        string
}

func InitPlayer() *Player {
	player := &Player{}
	player.AccountID = 0

	return player
}

func NewPlayer(name string, heroTemplateID int32) (*Player, error) {

	playID, err := redis.GetIncrPlayerID()
	if err != nil {
		fmt.Println("incr player id error :", err)
		return nil, err
	}

	player := InitPlayer()
	player.AccountID = playID
	player.Name = name
	// 初始化玩家英雄
	player.HeroTeam = NewHeroTeams()
	player.HeroTeam.AddHero(player.Name, true, heroTemplateID)

	// 初始化玩家资源
	player.Res = NewResource()

	// 关卡数据初始化
	events, err := gamedata.GetGameLevelEvents(1)
	if err != nil {
		fmt.Println("GetGameLevelEvents(1) error")
		return nil, err
	}
	player.PlayerGameLevel = NewPlayerGameLevel()
	gameLevel := structs.GameLevel{
		GameLevelID:   1,
		IsUnlock:      true,
		CompleteEvent: make([]int32, len(events)),
	}
	player.PlayerGameLevel.AddGameLevel(&gameLevel)

	// 背包数据初始化
	player.Bag = NewBag()

	// 新手引导状态初始化
	player.UserGuidRecords = make([]*structs.GuildRecord, 0, 10)

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
