package test

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/cache"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

var logger = logrus.New()

type testStruct struct {
	Id string `json:"id"`
	Data string `json:"data"`
}

func TestCache(t *testing.T) {


	fmt.Println(1<<15)

	testA :=[]testStruct{{"1","4"},{"2","3"},{"3","2"},{"4","1"},{"abc","bcd"}}
	data,err := json.Marshal(testA)
	fmt.Println(string(data))

	var testB []testStruct

	err = json.Unmarshal(data,&testB)


	localCache, err := cache.NewCache("memory", `{"interval":60}`)
	if err != nil {
		logger.Errorf("创建本地缓存失败！：%s", err.Error())
		return
	}
	key := "testCache"
	if err = localCache.Put(key, 0, 24*time.Hour);err != nil{
		logger.Errorf("本地缓存存入：key:%s-val:%s","1","2")
	}

}
