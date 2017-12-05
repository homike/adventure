package service

import "Adventure/AdvServer/db/mysql"

var (
	PlayerDao *mysql.PlayerDao
	PlayerID  uint
)

func init() {
	PlayerDao = mysql.NewUserDao()
	PlayerID = 0
}

func IncrPlayerID() {
	PlayerID++
}
