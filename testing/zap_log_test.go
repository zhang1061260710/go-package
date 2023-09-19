package testing

import (
	"github.com/zhang1061260710/go-package/logger/logger"
	"go.uber.org/zap"
	"testing"
)

var zapLog *zap.Logger

func zapLogTest(t *testing.T) {
	zapLog = logger.NewZapLoggerClient().InitZap("test", false)

	zapLog.Info("test")
}
