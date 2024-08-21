package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_Ant_Api(t *testing.T) {
	host := "http://mayi-bapp-api.hsk8s-dev.avlyun.org"
	path := "/task/RealTimeQuery/v1"
	method := http.MethodPost
	sk := "mayi-app-data"
	tpl := "mayi-app-data"

	apks := RequestData{
		Request: []RequestItem{
			{
				PkgName: "com.lvxing.datangeo",
			},
		},
	}
	body, err := json.Marshal(apks)
	if err != nil {
		t.Error(err)
		return
	}
	ts := time.Now().Unix()
	sg, err := calcSign256([]byte(sk), method, tpl, path, body, ts)
	if err != nil {
		t.Error(err)
		return
	}
	query := fmt.Sprintf("sign=%s&tpl=%s&ts=%d", sg, tpl, ts)

	req, err := http.NewRequest(method, host+path+"?"+query, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	curl, _ := GenerateCurlCommand(req)
	t.Log(curl)

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
	dst := bytes.NewBufferString("")
	err = json.Indent(dst, respBody, "", "  ")
	if err != nil {
		t.Error(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		t.Error(strconv.Itoa(resp.StatusCode) + "\n" + string(respBody))
		return
	}

	t.Log(dst.String())

}

type RequestData struct {
	Request []RequestItem `json:"request"`
}

type RequestItem struct {
	PkgName string `json:"pkg_name"`
}

func Test_RealtimeApk(t *testing.T) {
	host := "http://mayi-bapp-api.hsk8s-dev.avlyun.org"
	path := "/api/query/realtime/apk/v1"
	method := http.MethodPost
	sk := "mayi-app-data"
	tpl := "mayi-app-data"

	apks := RequestData{
		Request: []RequestItem{
			{
				PkgName: "com.tencent.android.qqdownloader",
			}, {
				PkgName: "com.samsung.android.hmt.vrsystem",
			}, {
				PkgName: "com.sing.client",
			},
		},
	}
	body, err := json.Marshal(apks)
	if err != nil {
		t.Error(err)
		return
	}
	ts := time.Now().Unix()
	sg, err := calcSign256([]byte(sk), method, tpl, path, body, ts)
	if err != nil {
		t.Error(err)
		return
	}
	query := fmt.Sprintf("sign=%s&tpl=%s&ts=%d", sg, tpl, ts)

	req, err := http.NewRequest(method, host+path+"?"+query, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	curl, _ := GenerateCurlCommand(req)
	t.Log(curl)
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
	dst := bytes.NewBufferString("")
	err = json.Indent(dst, respBody, "", "  ")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst.String())
}

type TaskReq struct {
	TaskId string `json:"taskId"`
}

func Test_RealtimeApkStatus(t *testing.T) {
	host := "http://mayi-bapp-api.hsk8s-dev.avlyun.org"
	path := "/api/query/realtime/apk/status/v1"
	method := http.MethodPost
	sk := "mayi-app-data"
	tpl := "mayi-app-data"

	apks := TaskReq{
		TaskId: "2d64e440-d90e-4107-9c8c-226d222bd8cc",
	}
	body, err := json.Marshal(apks)
	if err != nil {
		t.Error(err)
		return
	}
	ts := time.Now().Unix()
	sg, err := calcSign256([]byte(sk), method, tpl, path, body, ts)
	if err != nil {
		t.Error(err)
		return
	}
	query := fmt.Sprintf("sign=%s&tpl=%s&ts=%d", sg, tpl, ts)

	req, err := http.NewRequest(method, host+path+"?"+query, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	curl, _ := GenerateCurlCommand(req)
	t.Log(curl)
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
	dst := bytes.NewBufferString("")
	err = json.Indent(dst, respBody, "", "  ")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst.String())
}

// GenerateCurlCommand generates a curl command from an http.Request.
func GenerateCurlCommand(req *http.Request) (string, error) {
	var curlCmd strings.Builder

	// Start with the curl command and the HTTP method
	curlCmd.WriteString("curl -X ")
	curlCmd.WriteString(req.Method)
	curlCmd.WriteString(" ")

	// Add headers
	for name, values := range req.Header {
		for _, value := range values {
			curlCmd.WriteString(fmt.Sprintf("-H '%s: %s' ", name, value))
		}
	}

	// Add the URL
	curlCmd.WriteString(fmt.Sprintf("'%s' ", req.URL.String()))

	// Add the request body if it's a POST/PUT/PATCH request
	if req.Body != nil && (req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch) {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return "", err
		}
		req.Body = io.NopCloser(strings.NewReader(string(bodyBytes))) // Reset body
		if len(bodyBytes) > 0 {
			curlCmd.WriteString(fmt.Sprintf("-d '%s'", string(bodyBytes)))
		}
	}

	return curlCmd.String(), nil
}

func Test_RealtimeApkResult(t *testing.T) {
	host := "http://mayi-bapp-api.hsk8s-dev.avlyun.org"
	path := "/api/query/realtime/apk/result/v1"
	method := http.MethodPost
	sk := "mayi-app-data"
	tpl := "mayi-app-data"

	apks := TaskReq{
		TaskId: "2d64e440-d90e-4107-9c8c-226d222bd8cc",
	}
	body, err := json.Marshal(apks)
	if err != nil {
		t.Error(err)
		return
	}
	ts := time.Now().Unix()
	sg, err := calcSign256([]byte(sk), method, tpl, path, body, ts)
	if err != nil {
		t.Error(err)
		return
	}
	query := fmt.Sprintf("sign=%s&tpl=%s&ts=%d", sg, tpl, ts)

	req, err := http.NewRequest(method, host+path+"?"+query, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	curl, _ := GenerateCurlCommand(req)
	t.Log(curl)
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
	dst := bytes.NewBufferString("")
	err = json.Indent(dst, respBody, "", "  ")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dst.String())
}
