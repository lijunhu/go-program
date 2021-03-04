package test

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"net/http"
	"testing"
)

var groupId = "5d807fed9725eb0007ce8431"
var janusToken = "5d07036b34572400072b4a70"
var janusHost = "http://janus.t.17usoft.com"

func TestExportAppUk(t *testing.T) {

	req := httplib.Get(janusHost + "/janus-api/api/server/" + groupId + "/list?pageSize=200&pageNum=1")
	req.Header("user-token", janusToken)
	resp, _ := req.Response()
	defer resp.Body.Close()
	body := resp.Body

	bodyBytes, _ := ioutil.ReadAll(body)
	var obj interface{}
	_ = json.Unmarshal(bodyBytes, &obj)
	result := obj.(map[string]interface{})["result"]
	pageData := result.(map[string]interface{})["pageData"]
	for _, app := range pageData.([]interface{}) {
		valueMap := app.(map[string]interface{})
		id := valueMap["id"].(string)
		nameServer := valueMap["nameServer"].(string)
		useNameServer := valueMap["useNameServer"].(bool)
		if useNameServer {

			req = httplib.Get("http://jean.corp.elong.com/ocean/api/node-search?name=" + nameServer + "&flag=false")
			cookie := &http.Cookie{Name: "jean", Domain: "jean.corp.elong.com", Value: "8abc4ed337cf4da0907ea5210d203191"}
			req.SetCookie(cookie)
			resp, _ = req.Response()
			body = resp.Body
			bodyBytes, _ = ioutil.ReadAll(body)
			_ = json.Unmarshal(bodyBytes, &obj)

			data := obj.(map[string]interface{})["data"]
			values := data.([]interface{})
			if len(values) > 0 {
				value := values[0].(map[string]interface{})

				sql := fmt.Sprintf("db.JanusApp.update({\"_id\":ObjectId(\"%s\"),\"groupId\":ObjectId(\"%s\")},{$set:{\"appUniqueKey\":\"%s\"}},false,true)", id, groupId, value["uk"])
				fmt.Println(sql)
			}
		}

	}
}
