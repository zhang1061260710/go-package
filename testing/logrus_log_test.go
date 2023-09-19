package testing

import (
	"github.com/sirupsen/logrus"
	"github.com/zhang1061260710/go-package/logger/logger"
	"testing"
)

var log *logrus.Logger

func logrusLogTest(t *testing.T) {
	log = logger.NewLogrusLoggerClient().InitLogrus("logrus_test")
	log.Info("test logrus")
}
