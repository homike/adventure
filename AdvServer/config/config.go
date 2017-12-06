package config

import (
	"flag"
	"time"

	"github.com/Terry-Mao/goconf"
)

var (
	Conf     *Config
	confFile string
)

func Init() error {
	var err error

	flag.StringVar(&confFile, "conf", "./server.conf", " set login config file path")
	flag.Parse()
	Conf, err = newConfig(confFile)
	if err != nil {
		return err
	}

	return nil
}

type Config struct {
	LogPath      string        `goconf:"log:path"`
	LogLevel     int           `goconf:"log:level"`
	RedisIdle    int           `goconf:"redis:idle"`
	RedisTimeout time.Duration `goconf:"redis:timeout:time"`
	RedisAddr    string        `goconf:"redis:addr"`
	RedisAuth    string        `goconf:"redis:auth"`
	RedisSwitch  int8          `goconf:"redis:switch"`
}

func newConfig(file string) (*Config, error) {
	gconf := goconf.New()
	if err := gconf.Parse(file); err != nil {
		return nil, err
	}

	// Default config
	conf := &Config{
		RedisIdle:    50,
		RedisTimeout: 240 * time.Second,
		RedisAddr:    "localhost:6379",
		RedisAuth:    "",
		RedisSwitch:  0,
	}

	if err := gconf.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func ReloadCfg() error {
	gconf := goconf.New()
	if err := gconf.Parse(confFile); err != nil {
		return err
	}

	if err := gconf.Unmarshal(Conf); err != nil {
		return err
	}

	return nil
}

func GetConfig() *Config {
	return Conf
}
