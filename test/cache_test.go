package test

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"testing"
	"time"
)

var logger = logrus.New()

type testStruct struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

func TestCache(t *testing.T) {

	start, _ := time.Parse(time.DateOnly, "2024-01-01")
	end, _ := time.Parse(time.DateOnly, "2024-04-18")
	var arr []string
	for tmp := start; tmp.Before(end); tmp = tmp.AddDate(0, 0, 1) {
		arr = append(arr, "'"+tmp.Format("20060102")+"'")
	}
	fmt.Println(strings.Join(arr, ","))
	_, err := regexp.Compile(`/[/]`)
	if err != nil {
		t.Error(err)
		return
	}

	// 示例 map
	var a interface{}
	data := []interface{}{"1", "2", "3"}
	a = []interface{}{"1", "2", "3"}

	// 创建表达式
	expression, err := govaluate.NewEvaluableExpression("('1','2','3') == data")
	if err != nil {
		panic(err)
	}

	parameters := map[string]interface{}{
		"data": data,
		"a":    a,
	}

	// 评估表达式
	result, err := expression.Evaluate(parameters)
	if err != nil {
		panic(err)
	}

	fmt.Println("Expression result:", result)

	fmt.Println(1 << 15)

}
