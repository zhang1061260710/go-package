package logger

import (
	"github.com/zhang1061260710/go-package/logger/logger/logrus"
	"github.com/zhang1061260710/go-package/logger/logger/zap"
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
