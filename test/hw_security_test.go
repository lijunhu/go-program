package test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"go-program/pkg/httpclient"
	"go-program/pkg/loader"
	"go-program/pkg/trace"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	TimeStampFMT     = "20060102T150405Z"
	TimeStampCredFMT = "20060102"

	AlgorithmKey = "X-Algorithm"
	Algorithm    = "HMAC-SHA256"

	CredentialKey  = "X-Credential" // #nosec G101
	SignKey        = "X-Sign"
	SignHeadersKey = "X-Sign-Headers"
	TimeStampKey   = "X-TimeStamp"
	semicolon      = ";"
)

var HeaderKeyArr = []string{TimeStampKey, AlgorithmKey, CredentialKey}

type SignItem struct {
	Method         string
	Path           string
	HeaderKeys     string
	Timestamp      string
	Algorithm      string
	Credential     string
	Body           []byte
	AccessID       string
	SecretKey      string
	ShortTimestamp string
	Antiy          string
}

func NewSignItem(ops ...Option) *SignItem {
	item := &SignItem{}

	for _, op := range ops {
		op(item)
	}
	return item
}

type Option func(item *SignItem)

func WithMethod(method string) Option {
	return func(item *SignItem) {
		item.Method = method
	}
}

func WithPath(path string) Option {
	return func(item *SignItem) {
		item.Path = path
	}
}

func WithHeaderKeys(headerKeys string) Option {
	return func(item *SignItem) {
		item.HeaderKeys = headerKeys
	}
}
func WithTimestamp(timestamp string) Option {
	return func(item *SignItem) {
		item.Timestamp = timestamp
	}
}
func WithAlgorithm(algorithm string) Option {
	return func(item *SignItem) {
		item.Algorithm = algorithm
	}
}
func WithCredential(credential string) Option {
	return func(item *SignItem) {
		item.Credential = credential
	}
}
func WithBody(body []byte) Option {
	return func(item *SignItem) {
		item.Body = body
	}
}
func WithAccessID(accessID string) Option {
	return func(item *SignItem) {
		item.AccessID = accessID
	}
}
func WithSecretKey(secretKey string) Option {
	return func(item *SignItem) {
		item.SecretKey = secretKey
	}
}
func WithShortTimestamp(shortTimestamp string) Option {
	return func(item *SignItem) {
		item.ShortTimestamp = shortTimestamp
	}
}

func WithAntiy(antiy string) Option {
	return func(item *SignItem) {
		item.Antiy = antiy
	}
}

func (s *SignItem) GetHwSecuritySign() (string, error) {
	h := sha256.New()
	_, err := h.Write(s.Body)
	if err != nil {
		return "", err
	}

	bodyHash := hex.EncodeToString(h.Sum(nil))

	tmpArr := []string{s.Method, s.Path}
	tmpArr = append(tmpArr, bodyHash)

	headerKeyArr := strings.Split(s.HeaderKeys, semicolon)
	if len(headerKeyArr) != 3 {
		return "", errors.New("headerKeys must eq 3")
	}

	for _, keyName := range headerKeyArr {
		switch keyName {
		case TimeStampKey:
			tmpArr = append(tmpArr, s.Timestamp)
		case AlgorithmKey:
			tmpArr = append(tmpArr, s.Algorithm)
		case CredentialKey:
			tmpArr = append(tmpArr, s.Credential)
		}
	}

	extractKey := hmac.New(sha256.New, []byte(s.SecretKey))
	_, err = extractKey.Write([]byte(s.AccessID))
	if err != nil {
		return "", err
	}

	signingKey := hmac.New(sha256.New, extractKey.Sum(nil))
	_, err = signingKey.Write([]byte(fmt.Sprintf("%s/%s", s.ShortTimestamp, s.Antiy)))
	if err != nil {
		return "", err
	}

	signStrHmac := hmac.New(sha256.New, signingKey.Sum(nil))
	_, err = signStrHmac.Write([]byte(strings.Join(tmpArr, "\n")))
	if err != nil {
		return "", err
	}

	sign = hex.EncodeToString(signStrHmac.Sum(nil))
	return sign, nil
}

func Test_Urls(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "e8776b63054bdad73421e299530ef986"
	accessKey := "768f7098ee56db1e899968105958dcae98533933ec9ea0efb39368dac44fc440"
	contentType := "application/json"
	method := http.MethodPost
	path := "/api/v1/urls"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.UTC().Format(TimeStampFMT)
	date := now.UTC().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{"request":{"url":"https://www.zhihu.com/"}}`
	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}

func Test_Apks(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "e8776b63054bdad73421e299530ef986"
	accessKey := "768f7098ee56db1e899968105958dcae98533933ec9ea0efb39368dac44fc440"
	contentType := "application/json"
	method := http.MethodPost
	path := "/api/v1/apks"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.Local().Format(TimeStampFMT)
	date := now.Local().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{"request":{"packageName": "com.taobao.idlefish"}}`
	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}

func Test_Detects_Apks(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "e8776b63054bdad73421e299530ef986"
	accessKey := "768f7098ee56db1e899968105958dcae98533933ec9ea0efb39368dac44fc440"
	contentType := "application/json"
	method := http.MethodPost
	path := "/api/v1/apks/detects"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.Local().Format(TimeStampFMT)
	date := now.Local().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{"requests":[{"packageName": "com.taobao.idlefish"},{"packageName": "zzy.devicetool"}]}`
	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}

func Test_Detects_Apks_Status(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "e8776b63054bdad73421e299530ef986"
	accessKey := "768f7098ee56db1e899968105958dcae98533933ec9ea0efb39368dac44fc440"
	contentType := "application/json"
	method := http.MethodPost
	path := "/api/v1/apks/detects/status"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.Local().Format(TimeStampFMT)
	date := now.Local().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{"request":{"taskId": "f7a4aaea-913a-44d4-96c4-24700b500624"}}`
	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}

func Test_Detects_Apks_Result(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "e8776b63054bdad73421e299530ef986"
	accessKey := "768f7098ee56db1e899968105958dcae98533933ec9ea0efb39368dac44fc440"
	contentType := "application/json"
	method := http.MethodPost
	path := "/api/v1/apks/detects/results"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.Local().Format(TimeStampFMT)
	date := now.Local().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{"request":{"taskId": "f7a4aaea-913a-44d4-96c4-24700b500624"}}`
	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}

func Test_Developer_Apps(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "e8776b63054bdad73421e299530ef986"
	accessKey := "768f7098ee56db1e899968105958dcae98533933ec9ea0efb39368dac44fc440"
	contentType := "application/json"
	method := http.MethodPost
	path := "/api/v1/developers/apps"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.Local().Format(TimeStampFMT)
	date := now.Local().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{"request":{"pubKey": "2486AD8B36184232BA76590CD5CF330"}}`
	//body := `{"request":{"signMd5": "847E659F0F1D8B6D375EEAC073B93B79"}}`

	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}

func Test_Verify_Detects(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "e8776b63054bdad73421e299530ef986"
	accessKey := "768f7098ee56db1e899968105958dcae98533933ec9ea0efb39368dac44fc440"
	contentType := "application/json"
	method := http.MethodPost
	path := "/api/v1/verify/detects"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.Local().Format(TimeStampFMT)
	date := now.Local().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{
        "requests": [
            {
                "md5": "",
                "packageName": "com.taobao.taobao",
                "version": "5.22.1",
                "priority": 0,
                "clues": [
                    {
                        "key": "",
                        "value": ""
                    },
                    {
                        "key": "iconHash",
                        "value": "a1b2c3d4e5f6"
                    }
                ]
            },
            {
                "md5": "22633572CFEFEA1C229185B7248368A4",
                "packageName": "",
                "version": "",
                "priority": 0,
                "clues": [
                    {
                        "key": "iconHash",
                        "value": "aaaaaa"
                    },
                    {
                        "key": "iconHash",
                        "value": "bbbbbbbb"
                    }
                ]
            }
        ]
    }`
	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}

func Test_Verify_Status(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "e8776b63054bdad73421e299530ef986"
	accessKey := "768f7098ee56db1e899968105958dcae98533933ec9ea0efb39368dac44fc440"
	contentType := "application/json"
	method := http.MethodPost
	path := "/api/v1/verify/detects/status"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.Local().Format(TimeStampFMT)
	fmt.Println(timestamp)
	date := now.Local().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{"request":{"taskId":"963ae09d-41c7-48d3-b499-720dc5331889"}}`
	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}

func Test_Verify_Result(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "e8776b63054bdad73421e299530ef986"
	accessKey := "768f7098ee56db1e899968105958dcae98533933ec9ea0efb39368dac44fc440"
	contentType := "application/json"
	method := http.MethodPost
	path := "/api/v1/verify/detects/results"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.UTC().Format(TimeStampFMT)
	ti := now.UTC().Format(TimeStampFMT)

	fmt.Println(timestamp + "------------------" + ti)
	date := now.UTC().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{"request":{"taskId":"963ae09d-41c7-48d3-b499-720dc5331889"}}`
	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}

func Test_Verify_Operate_User(t *testing.T) {
	if err := loader.Load(); err != nil {
		t.Error(err)
		return
	}

	now := time.Now()

	accessID := "0801d3ad4dd5beb778d5361ff0867353"
	accessKey := "881534077ed11b4ae1c0b46d1fa943e1c8a7f966eabe758ff7e18cb750a9bbe6"
	contentType := "application/json"
	method := http.MethodPost
	path := "/inner/v1/user"
	host := "http://security-service-api.hsk8s-dev.avlyun.org"
	//host := "http://127.0.0.1:8080"

	timestamp := now.Local().Format(TimeStampFMT)
	date := now.Local().Format(TimeStampCredFMT)
	headerKeys := strings.Join(HeaderKeyArr, semicolon)
	credential := fmt.Sprintf("%s/%s/antiy", accessID, date)
	antiy := "antiy"

	body := `{"opt":"getall"}`
	item := NewSignItem(
		WithMethod(method),
		WithPath(path),
		WithAccessID(accessID),
		WithSecretKey(accessKey),
		WithAlgorithm(Algorithm),
		WithTimestamp(timestamp),
		WithShortTimestamp(date),
		WithBody([]byte(body)),
		WithHeaderKeys(headerKeys),
		WithAntiy(antiy),
		WithCredential(credential),
	)
	securitySign, err := item.GetHwSecuritySign()
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	req.Header.Set(AlgorithmKey, Algorithm)
	req.Header.Set(TimeStampKey, timestamp)
	req.Header.Set(CredentialKey, credential)
	req.Header.Set(SignHeadersKey, headerKeys)
	req.Header.Set(SignKey, securitySign)

	req.Header.Set("Content-Type", contentType)
	ctx := trace.NewTraceCtx(context.Background())
	resp, err := httpclient.RequestWithContext(ctx, req, 10)
	if err != nil {
		t.Error(err)
		return
	}
	respBody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}
	t.Log(string(respBody))

}
