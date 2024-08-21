package test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func Test_Mi_Apk_Api(t *testing.T) {
	host := "http://mi-antiy-server.hsk8s-dev.avlyun.org"
	path := "/v1/input/mi"
	method := http.MethodPost
	sk := "xiaomi"
	tpl := "xiaomi"

	apks := Apks{
		Tpl: tpl,
		Data: []*TMiTaierApk{
			/*			{
							ApkHash: "0e17b1c4f758b96a0f1de44addc3590b",
							ApkPath: "https://d2psj5hkq8ko3n.cloudfront.net/down/hls/hls_1.4.0_240613_6.apk",
						},
						{
							ApkHash: "f80c659058ba276302b97c0bca0e6f90",
							ApkPath: "https://app.likehuanxin.com/huanxin2.0.7.apk",
						},*/
			{
				ApkHash: "41880d6c65865abcd166ad8b989d0056",
				ApkPath: "https://okmtre.rr608ij5.top/zzxhy/zzxhy338292.apk",
			},
			{
				ApkPath: "https://dx18.635528.com/com.wisedu.cpdaily.nbu1006.apk",
				ApkHash: "dbecd67c89b67c1a6753675a24813a15",
			},
			{
				ApkHash: "45ae5ecaba3551010f243c8ad17f6eb5",
				ApkPath: "https://download.kfc.com.cn/KFC_Brand.apk",
			},
		},
	}
	body, err := json.Marshal(apks)
	if err != nil {
		t.Error(err)
		return
	}
	ts := time.Now().Unix()
	sg, err := calcSign256([]byte(sk), method, tpl, path, body, int64(ts))
	if err != nil {
		t.Error(err)
		return
	}
	query := fmt.Sprintf("sign=%s&tpl=%s&ts=%d", sg, tpl, ts)

	req, err := http.NewRequest(method, host+path+"?"+query, bytes.NewBuffer(body))

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

type Apks struct {
	Tpl  string         `json:"tpl"`
	Data []*TMiTaierApk `json:"data"`
}

// TMiTaierApk 根据方向C要求，表名为t_mi_taier_apk
type TMiTaierApk struct {
	ReqId                string `json:"reqId"`
	Tpl                  string `json:"tpl"`
	ApkHash              string `json:"apkHash"`
	ApkPath              string `json:"apkPath"`
	ApkSize              int64  `json:"apkSize"`
	AppInfo              string `json:"appInfo"`
	AppSourceValue       string `json:"appSourceValue"`
	AppSourcepackageName string `json:"appSourcepackageName"` //先前交付文档这么写的
	Fingerprint          string `json:"fingerprint"`
	InstallTime          int64  `json:"installTime"`
	Name                 string `json:"name"`
	PackageName          string `json:"packageName"`
	VersionCode          string `json:"versionCode"`
	VersionName          string `json:"versionName"`
	Permission           string `xorm:"text" json:"permission"`
	// 以下使用指针的项是因为其值 0和1都是有效值
	TouchInstallCnt int64         `json:"touchInstallCnt"`
	InstallCnt      int64         `json:"installCnt"`
	Permissions     []string      `xorm:"text" json:"permissions"`
	SftpPath        string        `json:"sftpPath"`
	IconPath        string        `json:"iconPath"`
	SdkCount        int64         `json:"sdkCount"`
	SdkList         []interface{} `xorm:"text" json:"sdkList"`
	CreatedAt       int64         `xorm:"created"`
	UpdatedAt       int64         `xorm:"updated"`
	SyncStatus      int64         `json:"syncStatus"`
	Date            string        `xorm:"-" json:"date"` // 生成文件目录日期  yyyymmdd
}

func calcSign256(sk []byte, method string, appId string, urlPath string, body []byte, ts int64) (string, error) {
	// timestamp UNIX 时间戳	携带在 url params 中，精确到秒
	// app_id	分配给用户的 app_id	携带在 url params 中
	// method	请求的 HTTP Method	HTTP Method，注意必须为大写
	// url_path	请求的 URL Path	URL 中的请求路径
	// content_hash	请求的 HTTP Body 的 SHA256 值	GET 请求时，可以不携带此字段； POST请求时，计算请求消息体的SHA256
	var signStr string
	var bodyHash string
	h := sha256.New()
	_, err := h.Write(body)
	if err != nil {
		return "", err
	}
	bodyHash = hex.EncodeToString(h.Sum(nil))
	switch method {
	case "GET":
		signStr = fmt.Sprintf("method=%s&tpl=%s&ts=%d&url_path=%s", method, appId, ts, urlPath)
	default:
		signStr = fmt.Sprintf("content_hash=%s&method=%s&tpl=%s&ts=%d&url_path=%s", bodyHash, method, appId, ts, urlPath)
	}
	h = hmac.New(sha256.New, sk)
	_, err = h.Write([]byte(signStr))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
