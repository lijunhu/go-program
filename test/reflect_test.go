package test

import (
	"reflect"
	"testing"
)

type StructSql struct {
	A string
	B string
}

type As struct {
	A string
}

func TestReflect(t *testing.T) {
	Aa := As{A: "a"}
	Ab := As{A: "b"}
	Ac := As{A: "c"}
	arr := []As{Aa, Ab, Ac}
	arrB := make([]*As, 0, 3)
	for i := 0; i < len(arr); i++ {
		tmp := arr[i]

		arrB = append(arrB, &tmp)
	}

	arr[0] = As{A: "abc"}

	for i := range arrB {
		t.Log(arrB[i].A)
	}

}

func appendSql() {

}

func compareField(s, t reflect.Value) (ok bool) {

	if s.Kind() != t.Kind() {
		return false
	}

	switch s.Kind() {
	case reflect.Struct, reflect.Interface:
		{
			for i := 0; i < s.NumField(); i++ {
				s := s.Field(i)
				t := t.Field(i)
				if ok = compareField(s, t); !ok {
					t.Set(s)
					return
				}
			}
		}
	case reflect.Map:
		{
			iterS := s.MapRange()
			iterT := s.MapRange()
			for iterS.Next(); iterT.Next(); {
				if ok = compareField(iterS.Key(), iterT.Key()); !ok {
					return
				}
				if ok = compareField(iterS.Value(), iterT.Key()); !ok {
					return
				}
			}
		}
	case reflect.Array, reflect.Slice:
		{

		}

	case reflect.Ptr:
		{
			s = s.Elem()
			t = t.Elem()
			if ok = compareField(s, t); !ok {
				return
			}
		}
	case reflect.Invalid:
		{
			return false
		}

	default:
		return s == t
	}
	return
}

func traverseField(parent, curr reflect.Value, ) {
	for i := 0; i < curr.NumField(); i++ {
		field := curr.Field(i)
		switch curr.Kind() {
		case reflect.Struct:
			{
				parent = curr
				traverseField(parent, field)
			}
		case reflect.Map:
			{

			}
		case reflect.Array:
		case reflect.Interface:

		}
	}

}
