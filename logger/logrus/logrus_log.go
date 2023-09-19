package logrus

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zhang1061260710/go-package/common"
	"io/ioutil"
	"time"
)

type LogrusLogger struct{}

type logrusLoggerConf struct {
	WithMaxAge       int64 //多久清理一次,单位 :H
	WithRotationTime int64 //多久切一次,单位 :H
	WithRotationSize int64 //超过多大切一次 单位：M
}

var conf logrusLoggerConf

func (l *LogrusLogger) Init() {
	viper.AddConfigPath("./etc/")
	viper.SetConfigName("logger")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		//log.Fatal("LogrusLogger init ReadInConfig err=", err)
		return
	}
	if err := viper.Unmarshal(&conf); err != nil {
		//log.Fatal("LogrusLogger init Unmarshal err=", err)
		return
	}
}

//logPath: 项目路径+文件名，比如：传参 /etc/test
func (l *LogrusLogger) InitLogrus(logPath string) *logrus.Logger {
	log := logrus.New()

	writerInfo, writerError, writerWarn := getWriter(logPath)
	hook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  writerInfo,
		logrus.ErrorLevel: writerError,
		logrus.WarnLevel:  writerWarn,
	}, &DiyFormatter{})
	log.Hooks.Add(hook)
	log.Out = ioutil.Discard
	return log
}
func (l *LogrusLogger) InitLogrusWithFormatter(logPath string, formatter logrus.Formatter) *logrus.Logger {
	log := logrus.New()

	writerInfo, writerError, writerWarn := getWriter(logPath)
	hook := lfshook.NewHook(lfshook.WriterMap{
		logrus.InfoLevel:  writerInfo,
		logrus.ErrorLevel: writerError,
		logrus.WarnLevel:  writerWarn,
	}, formatter)
	log.Hooks.Add(hook)
	log.Out = ioutil.Discard
	return log
}
func getWriter(logPath string) (*rotatelogs.RotateLogs, *rotatelogs.RotateLogs, *rotatelogs.RotateLogs) {

	withMaxAge := common.WithMaxAge
	withRotationTime := common.WithRotationTime
	withRotationSize := common.WithRotationSize

	if conf != (logrusLoggerConf{}) {
		withMaxAge = conf.WithMaxAge
		withRotationTime = conf.WithRotationTime
		withRotationSize = conf.WithRotationSize
	}

	writerInfo, _ := rotatelogs.New(
		logPath+".info.log-%Y%m%d",
		rotatelogs.WithLinkName(logPath+".info.log"),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(withMaxAge)),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(withRotationTime)),
		rotatelogs.WithRotationSize(withRotationSize*1024*1024),
	)
	writerError, _ := rotatelogs.New(
		logPath+".error.log-%Y%m%d",
		rotatelogs.WithLinkName(logPath+".error.log"),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(withMaxAge)),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(withRotationTime)),
		rotatelogs.WithRotationSize(withRotationSize*1024*1024),
	)
	writerWarn, _ := rotatelogs.New(
		logPath+".warn.log-%Y%m%d",
		rotatelogs.WithLinkName(logPath+".warn.log"),
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(withMaxAge)),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(withRotationTime)),
		rotatelogs.WithRotationSize(withRotationSize*1024*1024),
	)
	return writerInfo, writerError, writerWarn
}

type DiyFormatter struct {
}

func (m *DiyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string
	newLog = fmt.Sprintf("[%s] [%s] %s - ", timestamp, entry.Level, entry.Message)
	b.WriteString(newLog)
	for k, v := range entry.Data {
		if m, ok := v.(error); ok {
			b.WriteString(fmt.Sprintf("%s:%s\t", k, m.Error()))
		} else {
			data, err := jsoniter.MarshalToString(v)
			if err != nil {
				println(err)
			}
			b.WriteString(fmt.Sprintf("%s:%s\t", k, data))
		}
	}
	b.WriteString("\n")
	return b.Bytes(), nil
}
