package redis

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"

	"adventure/advserver/config"
)

var (
	pool *redis.Pool
	conf *config.Config
	//logger *nlog.Logger
)

func Init() {
	//logger = log.GetLogger()
	conf = config.GetConfig()

	// init redis pool
	pool = &redis.Pool{
		MaxIdle:     conf.RedisIdle,
		IdleTimeout: conf.RedisTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.RedisAddr)
			if err != nil {
				return nil, err
			}
			if conf.RedisAuth == "" { //
				return c, nil
			}
			if _, err := c.Do("AUTH", conf.RedisAuth); err != nil {
				fmt.Printf("c.Do('AUTH', %v) failed(%v) \n", conf.RedisAuth, err)
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			//syslog.Printf("redis ping with %v conn and time %v and err %v \n", c, t, err)
			return err
		},
	}
}

func GetPool() *redis.Pool {
	return pool
}
