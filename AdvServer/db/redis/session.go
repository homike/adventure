package redis

import (
	"github.com/garyburd/redigo/redis"
)

func GetIncrPlayerID() (uint, error) {
	con := pool.Get()
	defer con.Close()

	playerIDKey := "playerid"
	newPlayerID, err := redis.Int(con.Do("INCR", playerIDKey))
	if err != nil {
		return 0, err
	}

	return uint(newPlayerID), nil
}
