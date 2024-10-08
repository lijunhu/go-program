package test

import (
	"net/http"
	"testing"
)
import "github.com/husobee/vestigo"



func TestRoute(t *testing.T) {
	router := vestigo.NewRouter()
	router.Add("GET", "/abc/efg", nil)
	router.Add("GET", "/efg", nil)
	router.Add("GET", "/ghn", nil)
	router.Add("GET", "/hhh", nil)

	request, _ := http.NewRequest("GET", "/abc/efg", nil)
	h := router.Find(request)
	if h==nil {
		t.Log("adasdasda")
	}
}
