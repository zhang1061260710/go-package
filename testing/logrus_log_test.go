package testing

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/go-package/logger/logger"
	"testing"
)

var log *logrus.Logger

func logrusLogTest(t *testing.T) {
	log = logger.NewLogrusLoggerClient().InitLogrus("logrus_test")
	log.Info("test logrus")
}
