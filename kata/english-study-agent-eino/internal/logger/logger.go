package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	var err error
	config := zap.NewDevelopmentConfig()
	Log, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func Get() *zap.Logger {
	return Log
}

