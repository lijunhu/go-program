package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-program/pkg/loader"
	"go-program/pkg/trace"
	"os"
	"runtime"
	"strconv"
)

var ins = &Logger{}

type Logger struct {
	*logrus.Logger
}

type conf struct {
	logLevel string
}

func (l *Logger) GetName() string {
	return "Logger"
}

const (
	dateMillTimeFormat = "2006-01-02 15:04:05.000"
)

func (l *Logger) RunLoad() error {
	lvl := logrus.InfoLevel

	s := logrus.New()
	s.SetOutput(os.Stdout)
	s.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: dateMillTimeFormat,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, frame.File + ":" + strconv.Itoa(frame.Line)
		},
	})
	s.SetLevel(lvl)
	s.AddHook(trace.NewTraceIdHook())
	l.Logger = s
	return nil
}

func init() {
	loader.Register(ins)
}

func GetLoggerWithContext(ctx context.Context) *logrus.Entry {
	return ins.WithContext(ctx)

}
