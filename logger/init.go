package logger

import (
	"gitlab.com/go-package/logger/logger/logrus"
	"gitlab.com/go-package/logger/logger/zap"
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
