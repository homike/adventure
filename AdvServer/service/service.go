package service

import (
	"Adventure/AdvServer/db/mysql"
	"Adventure/AdvServer/db/redis"
	"Adventure/AdvServer/gamedata"
	"fmt"
)

var (
	PlayerDao *mysql.PlayerDao
)

func Init() error {
	PlayerDao = mysql.NewUserDao()

	redis.Init()

	gamedata.Init()

	name, err := gamedata.AllTemplates.HeroTemplate.HeroName(10109)
	skill, _ := gamedata.AllTemplates.HeroTemplate.SkillID(10109)

	fmt.Println("heroName: ", name, "skill ", skill, " error :", err)

	return nil
}
