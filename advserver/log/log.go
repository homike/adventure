package log

import (
	"adventure/common/clog"
	"time"
)

var (
	nuanLog *clog.Logger
)

// InitLogger init logger
func Init(file string, lvl int) error {
	l := clog.New(file, lvl, clog.Rotate{Size: clog.GB, Expired: time.Hour * 24 * 7, Interval: time.Hour})
	nuanLog = l

	return nil
}

// GetLogger return a logger
func GetLogger() *clog.Logger {
	return nuanLog
}

func SetLevel(lvl int) {
	nuanLog.SetLevel(lvl)
}
