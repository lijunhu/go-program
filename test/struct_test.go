package test

import (
	"reflect"
	"testing"
)

type A struct {
	a  string  `json:"a"`
	b  string  `json:"b"`
	c  int     `json:"c"`
	d  *string `json:"d"`
	ss *SS     `json:"ss"`
}

type SS struct {
	aa *string `json:"aa"`
	bb *string `json:"bb"`
	cc *int    `json:"cc"`
}

func TestStruct(t *testing.T) {

	num := 1
	str := "a"
	aA := A{
		a: str,
		b: str,
		c: num,
		d: &str,
		ss: &SS{
			aa: &str,
			bb: &str,
			cc: &num,
		},
	}

	ret := traverseStruct(reflect.ValueOf(aA))
	t.Log(ret)
}

func traverseStruct(v reflect.Value) (ret map[string]interface{}) {

	ret = make(map[string]interface{})
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		tmp := v.Field(i)
	PTR:
		switch tmp.Kind() {
		case reflect.Ptr:
			tmp = tmp.Elem()
			goto PTR
		case reflect.Struct:
			{
				ret[t.Field(i).Tag.Get("json")] = traverseStruct(tmp)
			}
		case reflect.String:
			ret[t.Field(i).Tag.Get("json")] = tmp.String()
		case reflect.Int:
			ret[t.Field(i).Tag.Get("json")] = tmp.Int()
		case reflect.Bool:
			ret[t.Field(i).Tag.Get("json")] = tmp.Bool()
		case reflect.Interface:
			ret[t.Field(i).Tag.Get("json")] = tmp.Interface()
		default:
			{
			}
		}

	}

	return
}
