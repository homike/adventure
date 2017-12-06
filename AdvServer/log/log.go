package log

import (
	"nuanv3/shared/nlog"
	"time"
)

var (
	nuanLog *nlog.Logger
)

// InitLogger init logger
func Init(file string, lvl int) error {
	l := nlog.New(file, lvl, nlog.Rotate{Size: nlog.GB, Expired: time.Hour * 24 * 7, Interval: time.Hour})
	nuanLog = l

	return nil
}

// GetLogger return a logger
func GetLogger() *nlog.Logger {
	return nuanLog
}

func SetLevel(lvl int) {
	nuanLog.SetLevel(lvl)
}
