package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	logger *logrus.Logger
	once   sync.Once
)

const (
	dateMillTimeFormat = "2006-01-02 15:04:05.000"
)

func GetLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}
	once.Do(
		func() {
			logger = logrus.New()
			infoLogPath := "./log/info.log"
			infoWriter, _ := rotatelogs.New(
				infoLogPath+".%Y-%m-%d-%H",
				rotatelogs.WithLinkName(infoLogPath),
				rotatelogs.WithMaxAge(time.Hour),
				rotatelogs.WithRotationTime(time.Hour),
			)

			fatalLogPath := "./log/fatal.log"
			fatalWriter, _ := rotatelogs.New(
				fatalLogPath+".%Y-%m-%d-%H",
				rotatelogs.WithLinkName(fatalLogPath),
				rotatelogs.WithMaxAge(time.Hour),
				rotatelogs.WithRotationTime(time.Hour),
			)

			debugLogPath := "./log/debug.log"
			debugWriter, _ := rotatelogs.New(
				debugLogPath+".%Y-%m-%d-%H",
				rotatelogs.WithLinkName(debugLogPath),
				rotatelogs.WithMaxAge(time.Hour),
				rotatelogs.WithRotationTime(time.Hour),
			)

			warnLogPath := "./log/warn.log"
			warnWriter, _ := rotatelogs.New(
				warnLogPath+".%Y-%m-%d-%H",
				rotatelogs.WithLinkName(warnLogPath),
				rotatelogs.WithMaxAge(time.Hour),
				rotatelogs.WithRotationTime(time.Hour),
			)

			errorLogPath := "./log/error.log"
			errorWriter, _ := rotatelogs.New(
				errorLogPath+".%Y-%m-%d-%H",
				rotatelogs.WithLinkName(errorLogPath),
				rotatelogs.WithMaxAge(time.Hour),
				rotatelogs.WithRotationTime(time.Hour),
			)

			panicLogPath := "./log/panic.log"
			panicWriter, _ := rotatelogs.New(
				panicLogPath+".%Y-%m-%d-%H",
				rotatelogs.WithLinkName(panicLogPath),
				rotatelogs.WithMaxAge(time.Hour),
				rotatelogs.WithRotationTime(time.Hour),
			)

			traceLogPath := "./log/trace.log"
			traceWriter, _ := rotatelogs.New(
				traceLogPath+".%Y-%m-%d-%H",
				rotatelogs.WithLinkName(traceLogPath),
				rotatelogs.WithMaxAge(time.Hour),
				rotatelogs.WithRotationTime(time.Hour),
			)

			writeMap := lfshook.WriterMap{
				logrus.InfoLevel:  infoWriter,
				logrus.FatalLevel: fatalWriter,
				logrus.DebugLevel: debugWriter,
				logrus.WarnLevel:  warnWriter,
				logrus.ErrorLevel: errorWriter,
				logrus.PanicLevel: panicWriter,
				logrus.TraceLevel: traceWriter,
			}
			logger.SetReportCaller(true)
			lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
				TimestampFormat:  dateMillTimeFormat,
				DisableTimestamp: false,
				CallerPrettyfier: nil,
			})
			logger.SetLevel(logrus.TraceLevel)
			logger.AddHook(lfHook)
		},
	)
	return logger
}
