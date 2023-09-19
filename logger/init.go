package logger

import (
	"github.com/zhang1061260710/go-package/logger/logrus"
	"github.com/zhang1061260710/go-package/logger/zap"
)

func NewLogrusLoggerClient() *logrus.LogrusLogger {
	log := &logrus.LogrusLogger{}
	log.Init()
	return log

}
func NewZapLoggerClient() *zap.ZapLogger {
	log := &zap.ZapLogger{}
	log.Init()
	return log

}
