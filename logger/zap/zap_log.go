package zap

import (
	"context"
	"github.com/hashicorp/go-uuid"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"github.com/zhang1061260710/go-package/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

type ZapLogger struct{}

type zapLoggerConf struct {
	WithMaxAge       int64 //多久清理一次,单位 :H
	WithRotationTime int64 //多久切一次,单位 :H
	WithRotationSize int64 //超过多大切一次 单位：M
}

var conf zapLoggerConf

func (ZapLogger) Init() {
	viper.AddConfigPath("./etc/")
	viper.SetConfigName("logger")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		//log.Fatal("ZapLogger init err=", err)
		return
	}
	if err := viper.Unmarshal(&conf); err != nil {
		//log.Fatal("ZapLogger init err:=", err)
		return
	}
}

//logPath: 项目路径+文件名，比如：传参 /etc/test
//enableConsoleOutPut : 开启控制台输出
func (l *ZapLogger) InitZap(logPath string, enableConsoleOutPut bool) *zap.Logger {
	return globalLogZap(logPath, enableConsoleOutPut)
}

func (l *ZapLogger) WithContext(ctx context.Context, writeLog *zap.Logger) *zap.Logger {
	requestId := ctx.Value(common.TrackId)
	uid := ctx.Value("uid")
	if requestId == nil {
		requestId, _ = uuid.GenerateUUID()
	}
	logger := writeLog.WithOptions(zap.Fields(zap.String(common.TrackId, requestId.(string))))
	if uid != nil {
		logger = logger.WithOptions(zap.Fields(zap.Int64("uid", uid.(int64))))
	}
	return logger
}

func globalLogZap(logPath string, enableConsoleOutPut bool) *zap.Logger {
	infoPath := logPath + "." + "info.log"
	errorPath := logPath + "." + "error.log"
	encode := newEncoder()
	// 设置日志级别
	// 实现两个判断日志等级的interface (其实 zapcore.Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel
	})
	if enableConsoleOutPut {
		infoLevel = zap.LevelEnablerFunc(func(level zapcore.Level) bool {
			return level < zapcore.WarnLevel
		})
	}

	warnLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.WarnLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	core := zapcore.NewTee(
		zapcore.NewCore(encode, zapcore.AddSync(NewWriter(infoPath)), infoLevel),
		zapcore.NewCore(encode, zapcore.AddSync(NewWriter(errorPath)), warnLevel),
	)
	if enableConsoleOutPut {
		core = zapcore.NewTee(
			zapcore.NewCore(encode, zapcore.AddSync(io.Writer(os.Stdout)), debugLevel),
			core,
		)
	}
	// 构造日志 开启开发模式，堆栈跟踪 开启文件及行号 跳过当前调用方，防止将封装方法当作调用方
	return zap.New(core, zap.AddCaller(), zap.Development(), zap.AddCallerSkip(1))
}

func newEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logx",
		CallerKey:     "linenum",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(common.DateFullFormat))
		}, // 时间格式化
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器，用于显示Method
		EncodeName:     zapcore.FullNameEncoder,
	})
}

func NewWriter(logName string) io.Writer {
	withMaxAge := common.WithMaxAge
	withRotationTime := common.WithRotationTime
	withRotationSize := common.WithRotationSize
	if conf != (zapLoggerConf{}) {
		withMaxAge = conf.WithMaxAge
		withRotationTime = conf.WithRotationTime
		withRotationSize = conf.WithRotationSize
	}

	writer, err := rotatelogs.New(
		// 日志文件
		strings.Replace(logName, ".log", "", -1)+"-%Y%m%d.log",
		rotatelogs.WithLinkName(logName),
		// 日志周期(默认每86400秒/一天旋转一次)
		rotatelogs.WithRotationTime(time.Hour*time.Duration(withMaxAge)),
		// 清除历史 (WithMaxAge和WithRotationCount只能选其一)
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(withRotationTime)), //每30天清除下日志文件
		rotatelogs.WithRotationSize(withRotationSize*1024*1024),             //日志文件大小，单位byte
		//rotatelogs.WithRotationCount(10), //只保留最近的N个日志文件
	)
	if err != nil {
		panic(err)
	}
	return writer
}
