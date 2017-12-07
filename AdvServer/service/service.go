package service

import (
	"Adventure/AdvServer/db/mysql"
	"Adventure/AdvServer/db/redis"
	"Adventure/AdvServer/gamedata"
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
