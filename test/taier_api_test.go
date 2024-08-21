package test

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_Taier_Api(t *testing.T) {

	key := "xiaomiTerminal_crakey"
	secret := "DJ^YL072TAAq1aTx"
	contentType := "application/json"
	method := http.MethodPost
	path := "/interface/data/crawler/verifyApp"
	host := "http://124.126.23.132:8051"

	body := `{"crawlerCode":"xiaomiTerminal","md5":["1e280a592cadd22f472146303427d590","1e280a592cadd22f472146303427d591","704290f07f26131d4c649a39a153274b"]}`
	date := time.Now().Format(time.DateTime)
	sg, err := createSign([]byte(body), secret, method, date, contentType, path)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(sg)

	req, err := http.NewRequest(method, host+path, bytes.NewBuffer([]byte(body)))

	if err != nil {
		t.Error(err)
		return
	}
	auth := fmt.Sprintf("CRAWLER %s:%s", key, sg)
	req.Header.Set("Authorization", auth)
	req.Header.Set("DATE", date)
	req.Header.Set("Content-Type", contentType)
	resp, err := http.DefaultClient.Do(req)
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

func createSign(body []byte, sk, method, date, contentType, path string) (sign string, err error) {
	bodyHash := strings.ToUpper(fmt.Sprintf("%x", md5.Sum(body))) // #nosec G401

	content := method + "\n" + bodyHash + "\n" + contentType + "\n" + date + "\n" + path

	s := hmac.New(sha256.New, []byte(sk))
	_, err = s.Write([]byte(content))
	if err != nil {
		return
	}

	return hex.EncodeToString(s.Sum(nil)), nil
}
