package log

import (
	"adventure/common/clog"
	"time"
)

var (
	xLog *clog.Logger
)

// InitLogger init logger
func Init(file string, lvl int) error {
	l := clog.New(file, lvl, clog.Rotate{Size: clog.GB, Expired: time.Hour * 24 * 7, Interval: time.Hour})
	xLog = l

	return nil
}

// GetLogger return a logger
func GetLogger() *clog.Logger {
	return xLog
}

func SetLevel(lvl int) {
	xLog.SetLevel(lvl)
}
