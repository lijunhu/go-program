package test

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestTypeChange(t *testing.T) {
	var err error

	startTime, err := time.Parse("2006-01-02 15:04:05", "2019-09-27 23:59:00")
	if err ==nil{
		fmt.Println(startTime)
	}


	test :=MethodValueToMethods(10)

	fmt.Println(test)

	fmt.Println(math.MaxInt32)

	data := map[string]interface{}{
		"Qos": 1,
	}

	var byts []byte
	for _, v := range data {
		fmt.Println("type:", reflect.TypeOf(v).String())
	}

	fmt.Println("after do marshal and unmarshal...........")

	byts, err = json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	marshalData := make(map[string]interface{})

	err = json.Unmarshal(byts, &marshalData)

	for _, v := range marshalData {
		fmt.Println("type:", reflect.TypeOf(v).String())
	}
}


func MethodValueToMethods(value int) []string {
	sm := []int{1, 2, 4, 8}
	var ms []string
	for _, v := range sm {
		if value&v > 0 {
			method,_ := HttpIntToSting(v)
			ms = append(ms, method)
		}
	}
	return ms
}


func HttpIntToSting(intValue int) (httpMethod string, err error) {
	switch intValue {
	case 1:
		{
			httpMethod = "GET"
			break
		}
	case 2:
		{
			httpMethod = "POST"
			break
		}
	case 4:
		{
			httpMethod = "PUT"
			break
		}
	case 8:
		{
			httpMethod = "DELETE"
			break
		}
	default:
		{
			err = fmt.Errorf("%dHTTP方法不符合规范", intValue)
			break
		}
	}
	return
}