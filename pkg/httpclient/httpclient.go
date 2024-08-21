package httpclient

import (
	"bufio"
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"go-program/pkg/logger"
	"go-program/pkg/trace"
	"io"
	"net/http"
	"strings"
	"time"
)

type ClientWrapper struct {
	client *http.Client
}

var client = http.DefaultClient

const (
	XRequestID                 = "X-Request-Id"
	ContentTypeHeader          = "content-type"
	ApplicationJsonContentType = "application/json"
	UrlFormEncodeContentType   = "application/x-form-urlencode"
)

func RequestWithContext(ctx context.Context, req *http.Request, timeout int64) (resp *http.Response, err error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout))
	defer cancel()

	var traceID string
	traceID, err = trace.GetTraceID(ctx)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)
	req.Header.Set(XRequestID, traceID)
	start := time.Now()
	resp, err = client.Do(req)

	var reqBodyBytes []byte
	// 请求content-type为application/json和application/x-form-urlencode格式时读取请求body
	contentType := req.Header.Get(ContentTypeHeader)
	if strings.HasPrefix(contentType, ApplicationJsonContentType) || strings.HasPrefix(contentType, UrlFormEncodeContentType) {
		var reqBody io.ReadCloser
		reqBody, err = req.GetBody()
		if err != nil {
			return nil, err
		}
		reqBodyBytes, err = readBody(reqBody)
	}

	var respBodyBytes []byte
	respBodyBytes, err = readBody(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(respBodyBytes))

	logger.GetLoggerWithContext(ctx).WithFields(logrus.Fields{
		"method":   req.Method,
		"path":     req.URL.Path,
		"host":     req.URL.Host,
		"query":    req.URL.RawQuery,
		"body":     string(reqBodyBytes),
		"header":   req.Header,
		"status":   resp.Status,
		"err":      err,
		"respBody": string(respBodyBytes),
		"latency":  time.Since(start).Milliseconds(),
	})
	return
}

func readBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	var err error
	reader := bufio.NewReaderSize(body, 1024*1024*10)

	var sum []byte
	for {
		var (
			text   []byte
			prefix bool
		)
		text, prefix, err = reader.ReadLine()
		sum = append(sum, text...)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return []byte{}, err
		}
		var total, tmp []byte
		for prefix {
			tmp, prefix, err = reader.ReadLine()
			if err == io.EOF {
				err = nil
				break
			}
			if err != nil {
				return []byte{}, err
			}
			total = append(total, tmp...)
		}
		sum = append(sum, total...)
	}
	return sum, nil
}
