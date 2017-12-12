package mysql

import (
	"database/sql"
	"fmt"
	"time"
)

type PlayerDao struct {
	db *sql.DB
}

func NewUserDao() *PlayerDao {
	return &PlayerDao{db: AdvDB}
}

type PlayerDB struct {
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

const (
	GetAllPlayerDataSQL = "SELECT AccountId, Name, PlatformAccountId, PlatformTypes, GameZoneId, CreateTime, LastLoginTime," +
		" LastLogoffTime, BarrageSet, VipLevel, OnlineTime, HeroTeam, PlayerGameLevel, Bag, MiningMap, ExtendData from player" +
		" where player.AccountId = %d "

	CreateUserSQL = "INSERT INTO player(AccountId, Name, PlatformAccountId, PlatformTypes, GameZoneId, CreateTime, LastLoginTime," +
		" LastLogoffTime, BarrageSet, VipLevel, OnlineTime, HeroTeam, PlayerGameLevel, Bag, MiningMap, ExtendData) " +
		" VALUES(%d, '%s', %d, %d, %d, now(), now(), now(), '%s', %d, %d, '%s', '%s', '%s', '%s', '%s')"
)

func (dao *PlayerDao) GetPlayerAllData(accountid uint) (*PlayerDB, error) {
	user := new(PlayerDB)

	row := dao.db.QueryRow(fmt.Sprintf(GetAllPlayerDataSQL, accountid))
	switch err := row.Scan(&user.AccountID, &user.Name, &user.PlatformAccountID,
		&user.PlatformTypes, &user.GameZoneID, &user.CreateTime, &user.LastLoginTime, &user.LastLogoffTime, &user.BarrageSet, &user.VipLevel, &user.OnlineTime,
		&user.HeroTeam, &user.PlayerGameLevel, &user.Bag, &user.MiningMap, &user.ExtendData); {
	case err == nil:
	case err == sql.ErrNoRows:
		//logger.Info("no user %d in db", accountid)
	default:
		//logger.Error("scan failed for user %d for err: %s", accountid, err.Error())
		return nil, err
	}
	return nil, nil
}

func (dao *PlayerDao) CreatePlayer(player *PlayerDB) error {
	sql := fmt.Sprintf(CreateUserSQL, player.AccountID, player.Name, player.PlatformAccountID,
		player.PlatformTypes, player.GameZoneID, player.BarrageSet, player.VipLevel, player.OnlineTime,
		player.HeroTeam, player.PlayerGameLevel, player.Bag, player.MiningMap, player.ExtendData)
	//fmt.Println("sql :", sql)
	_, err := dao.db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
