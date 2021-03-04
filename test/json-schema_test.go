package test

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestJsonSchema(t *testing.T) {

	fmt.Println(1 << 30 * 10)
	fmt.Println(1024 * 1024 * 1024)

	fmt.Println(fmt.Sprintf("%05d", 200))

	id := bson.ObjectIdHex("5ea121dfb1e3c13c3508bc69")
	t.Logf("%s\n", id.Hex())

	a := []int{0, 1, 2, 3, 4, 5}
	t.Log(a[:4])

	json := "{\"age\":15,\"name\":\"TigerLee\",\"sex\":\"male\",\"score\":89}"
	schemaJson := "{\"title\":\"Example Schema\",\"type\":\"object\",\"properties\":{\"age\":{\"type\":\"integer\",\"minimum\":0},\"name\":{\"type\":\"string\"},\"sex\":{\"description\":\"\",\"type\":\"string\"},\"score\":{\"description\":\"分数\",\"type\":\"number\"}},\"required\":[\"name\",\"age\"]}"
	schemaLoader := gojsonschema.NewBytesLoader([]byte(schemaJson))
	jsonLoader := gojsonschema.NewBytesLoader([]byte(json))

	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		t.Error(err)
	}
	ret, err := schema.Validate(jsonLoader)

	if err != nil {
		t.Error(err)
	}

	t.Log(ret.Valid())

}
