package service

import "Adventure/AdvServer/db/mysql"
import "Adventure/AdvServer/db/redis"

var (
	PlayerDao *mysql.PlayerDao
)

func Init() error {
	PlayerDao = mysql.NewUserDao()

	redis.Init()

	return nil
}
