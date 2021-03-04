package service

import (
	"encoding/json"
	"fmt"
	"go-program/sword/servcie/reqAnalysis"
	"strings"
	"testing"
)

func Test_Query_All_Configs(t *testing.T) {
	groupId := "5db28fbeb1e3c11478896ae5"
	routeId := ""
	url := "http://10.160.92.88:6010/janus-api/api/uriconfig/" + groupId + "/config/list?pageSize=10000&pageNum=1&uri=&host=&status=0&using=&proxyPassType=&upstreamName=&id="
	headers := map[string]string{
		"user-token": "5e1be479fe33db383c0af207",
	}
	baseConfigs, err := reqAnalysis.QueryAllConfigs(url, headers)

	if err != nil {
		panic(err)
	}
	configIds, failedConfigIds := make(map[string]string, 16), make(map[string]string, 16)
	if len(baseConfigs) > 0 {

		var configDetail reqAnalysis.JanusRouteWeb
		for _, baseConfig := range baseConfigs {
			configIds[baseConfig.Id.Hex()] = "1"
			routeId = baseConfig.Id.Hex()
			url = "http://10.160.92.88:6010/janus-api/api/uriconfig/" + groupId + "/config/detail?routeId=" + routeId
			configDetail, err = reqAnalysis.QueryConfigDetail(url, headers)
			if err != nil {
				panic(err)
			}
			if configDetail.Decompress.Using && strings.ToLower(configDetail.Decompress.DecompressType) == "gzip" &&
				configDetail.Decrypt.Using && strings.ToLower(configDetail.Decrypt.DecryptType) == "aes" {
				configDetail.ReqAnalysis.Using = false
				configDetail.ReqAnalysis.Protocol = ""
				url = "http://10.160.92.88:6010/janus-api/api/uriconfig/" + groupId + "/config"
				err = reqAnalysis.UpdateConfig(url, headers, configDetail)
				if err != nil {
					failedConfigIds[baseConfig.Id.Hex()] = err.Error()
				}
			}
		}

		if len(configIds) > 0 {
			data, _ := json.Marshal(configIds)
			fmt.Println(string(data))
		}
		if len(failedConfigIds) > 0 {
			data, _ := json.Marshal(failedConfigIds)
			fmt.Println(string(data))
		}

	}
}
