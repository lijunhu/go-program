package logger

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-program/pkg/trace"
	"io"
	"strings"
	"time"
)

// AccessLog returns a gin.HandlerFunc
func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		uri := c.Request.RequestURI
		content := GetPostRequestBody(c)

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		ctx := c.Request.Context()
		ctx = trace.NewTraceCtxWithTraceID(ctx, trace.GetTraceIDFromGinCtx(c))
		log := GetLoggerWithContext(ctx)
		l := log.WithFields(logrus.Fields{
			"status": c.Writer.Status(),
			"method": c.Request.Method,
			"uri":    uri,
			// "ip":         c.ClientIP(),
			"user-agent": c.Request.UserAgent(),
			"latency":    fmt.Sprintf("%dms", latency.Milliseconds()),
			"body":       content,
		})
		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				l.Error(e)
			}
		} else {
			l.Info()
		}

	}
}

func GetPostRequestBody(c *gin.Context) string {
	var (
		content []byte
		err     error
	)
	if c.Request != nil && strings.Contains(c.GetHeader("Content-Type"), "application/json") {
		content, err = io.ReadAll(c.Request.Body)
		if err != nil {
			GetLoggerWithContext(c.Request.Context()).Errorf("read body err:%s", err)
			return ""
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(content))
	}
	return string(content)
}
