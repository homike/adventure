package service

import (
	"adventure/advserver/db/mysql"
	"adventure/advserver/db/redis"
	"adventure/advserver/gamedata"
)

var (
	PlayerDao *mysql.PlayerDao
)

func Init() error {
	PlayerDao = mysql.NewUserDao()

	redis.Init()

	gamedata.Init()

	return nil
}
