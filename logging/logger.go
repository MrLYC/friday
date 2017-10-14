package logging

import (
	"friday/config"
	"friday/utils"

	"github.com/sirupsen/logrus"
)

// GlobalLogger : global logger
var GlobalLogger *logrus.Logger

var (
	skipStackCount = 2
)

// Init : init GlobalLogger
func Init() {
	var (
		conf  = config.Configuration.Logging
		level logrus.Level
		err   error
	)
	GlobalLogger = logrus.New()
	level, err = logrus.ParseLevel(conf.Level)
	if err != nil {
		panic(err)
	}
	GlobalLogger.SetLevel(level)
}

// GetErrorLogFields : get error log stack info
func GetErrorLogFields(skipStack int) logrus.Fields {
	frame := utils.RunTimeStackFrame{}
	frame.InitWithSkip(skipStack + 1)
	return logrus.Fields{
		"file": frame.Name,
		"line": frame.Line,
	}
}

func Debugf(format string, args ...interface{}) {
	GlobalLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	GlobalLogger.Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	GlobalLogger.Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetErrorLogFields(skipStackCount)).Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetErrorLogFields(skipStackCount)).Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetErrorLogFields(skipStackCount)).Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetErrorLogFields(skipStackCount)).Fatalf(format, args...)
}
