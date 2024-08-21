package test

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	strace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace/noop"
	"log"
	"os"
	"testing"
)

func Test_Logrus(t *testing.T) {
	// Instrument logrus.
	logrus.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	)))

	spanCtx := trace.SpanContext{}

	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)

	noop.NewTracerProvider().Tracer("").Start(ctx, "test")

	// Use ctx to pass the active span.
	logrus.WithContext(ctx).
		WithError(errors.New("hello world")).
		WithField("foo", "bar").
		Error("something failed")
}

var client = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport), Timeout: time.Second * 10}

func Test_Opentelemetry_Logrus(t *testing.T) {

	// 初始化 Logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.SetOutput(os.Stdout)
		logrus.Warn("Failed to log to file, using default stderr")
	} else {
		logrus.SetOutput(file)
	}

	// 初始化 OpenTelemetry
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to initialize stdout export pipeline: %v", err)
	}

	tp := strace.NewTracerProvider(
		strace.WithBatcher(exp),
		strace.WithResource(resource.Default()),
	)
	defer func() { _ = tp.Shutdown(context.Background()) }()

	otel.SetTracerProvider(tp)

	// 创建一个 Tracer
	tracer := otel.Tracer("example-tracer-1")

	// 开始一个新的 span
	ctx, span := tracer.Start(context.Background(), "main")
	defer span.End()

	// 在日志中记录 Trace ID
	traceID := span.SpanContext().TraceID().String()
	logrus.WithFields(logrus.Fields{
		"traceID": traceID,
		"event":   "start",
		"type":    "main",
	}).Info("Application starting")

	// 调用示例函数
	exampleFunction(ctx)

	// 在日志中记录 Trace ID
	logrus.WithFields(logrus.Fields{
		"traceID": traceID,
		"event":   "end",
		"type":    "main",
	}).Info("Application ending")
}

func exampleFunction(ctx context.Context) {
	tracer := otel.Tracer("example-tracer-2")
	_, span := tracer.Start(ctx, "exampleFunction")
	defer span.End()

	// 从上下文中提取 Trace ID
	traceID := span.SpanContext().TraceID().String()

	// 模拟一些工作并在日志中记录 Trace ID
	logrus.WithFields(logrus.Fields{
		"traceID": traceID,
		"event":   "work",
		"type":    "exampleFunction",
	}).Info("Doing some work in exampleFunction")
}
