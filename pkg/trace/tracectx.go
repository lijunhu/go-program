package trace

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type traceKeyStruct struct {
}

func (t traceKeyStruct) String() string {
	return "trace"
}

var traceKey = traceKeyStruct{}

type traceInfo struct {
	traceID string
}

var (
	ErrTraceKeyNotFound     = errors.New("not found trace key")
	ErrTraceInfoTypeInvalid = errors.New("trace info type invalid")
)

func CopyTraceCtx(sCtx context.Context) (tCtx context.Context, err error) {
	var (
		traceID string
	)
	if traceID, err = GetTraceID(sCtx); err != nil {
		return nil, err
	}
	tCtx = context.Background()

	return context.WithValue(tCtx, traceKey, &traceInfo{traceID: traceID}), nil
}

func NewTraceCtx(ctx context.Context) context.Context {
	traceId := uuid.New().String()
	return context.WithValue(ctx, traceKey, &traceInfo{traceID: traceId})
}

func GetTraceIDIgnoreErr(ctx context.Context) string {
	var (
		t   *traceInfo
		tmp interface{}
		ok  bool
	)
	tmp = ctx.Value(traceKey)

	if tmp == nil {
		return ""
	}
	if t, ok = tmp.(*traceInfo); !ok {
		return ""
	}
	return t.traceID
}

func GetTraceID(ctx context.Context) (traceID string, err error) {
	var (
		t   *traceInfo
		tmp interface{}
		ok  bool
	)
	tmp = ctx.Value(traceKey)
	if tmp == nil {
		return "", ErrTraceKeyNotFound
	}
	if t, ok = tmp.(*traceInfo); !ok {
		return "", ErrTraceInfoTypeInvalid
	}
	return t.traceID, nil
}

func NewTraceCtxWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceKey, &traceInfo{traceID: traceID})
}

func GinCtxSetTraceID() gin.HandlerFunc {
	return func(gCtx *gin.Context) {
		var traceID string
		if traceID = gCtx.Request.Header.Get("X-Request-ID"); len(traceID) == 0 {
			traceID = uuid.New().String()
		}
		gCtx.Set(traceKey.String(), traceID)
	}
}

func GetTraceIDFromGinCtx(gCtx *gin.Context) string {
	var traceID string
	if traceID = gCtx.Request.Header.Get("X-Request-ID"); len(traceID) > 0 {
		return traceID
	}
	tmp, ok := gCtx.Get(traceKey.String())
	if !ok {
		return uuid.New().String()
	}
	traceID = tmp.(string)
	return traceID
}

const traceIDField = "traceId"

type Hook struct {
}

func NewTraceIdHook() logrus.Hook {
	return &Hook{}
}

func (hook *Hook) Fire(entry *logrus.Entry) error {

	ctx := entry.Context

	entry.Data[traceIDField] = GetTraceIDIgnoreErr(ctx)
	return nil
}

func (hook *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}
