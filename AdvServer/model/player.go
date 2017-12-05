package model

import (
	"Adventure/AdvServer/db/mysql"
	"Adventure/AdvServer/service"
	"fmt"
	"time"
)

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
	HeroTeam          string
	PlayerGameLevel   string
	Bag               string
	MiningMap         string
	ExtendData        string
}

func InitPlayer() *Player {
	player := &Player{}
	player.AccountID = 0

	return player
}

func NewPlayer(name string, heroTemplateID int32) *Player {
	player := InitPlayer()

	player.AccountID = service.PlayerID
	player.Name = name

	service.IncrPlayerID()

	dbData := &mysql.PlayerDB{
		AccountID: player.AccountID,
		Name:      player.Name,
	}
	err := service.PlayerDao.CreatePlayer(dbData)
	if err != nil {
		fmt.Println("NewPlayer() error %v", err)
		return nil
	}

	return player
}
