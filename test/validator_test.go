package test

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/stretchr/testify/assert"

	"reflect"
	"testing"
)

var validate *validator.Validate
var uni *ut.UniversalTranslator

func TestValidator(t *testing.T) {

	zh := zh.New()
	uni = ut.New(zh, zh)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("zh")
	validate = validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		value, ok := field.Tag.Lookup("desc")
		if ok {
			return value
		}
		return ""
	})
	_ = zh_translations.RegisterDefaultTranslations(validate, trans)

	//name := new(string)
	//age:=new(int)
	//score := new(int)

	persion := person{
		Name:  nil,
		Age:   nil,
		Score: nil,
	}

	rt := reflect.TypeOf(persion)

	for i := 0; i < rt.NumField(); i++ {
		value, ok := rt.Field(i).Tag.Lookup("validate")
		if ok {
			t.Logf("tag validate value :%s", value)
		}
	}

	err := validate.Struct(persion)

	errs := err.(validator.ValidationErrors)
	for _, errMsg := range errs.Translate(trans) {
		t.Log(errMsg+"\n")
	}

	assert.NoError(t, err, "失败")
}

type person struct {
	Name  *string `json:"name" validate:"required" desc:"abc"`
	Age   *int    `json:"age" validate:"required"`
	Score *int    `json:"score" validate:"required"`
}
