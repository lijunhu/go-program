package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var server = gin.New()


var validate *validator.Validate


func main() {

	go Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-signals
}

func Run() {
	server.Use(cors.New(cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    []string{"http://", "https://"},
		AllowMethods:    []string{"GET", "PUT", "POST", "DELETE", "OPTIONS", "HEAD", "PATCH"},
		AllowHeaders:    []string{"Origin", "Authorization", "Content-Type", "user-token", "Set-Cookie", "X-Requested-With"},
		ExposeHeaders:   []string{},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge:           50 * time.Second,
		AllowCredentials: true,
	}))

	server.GET("/test/http/client", func(context *gin.Context) {
		time.Sleep(30 * time.Second)
		context.JSON(http.StatusOK, "{\"msg\":\"success\"}")
	})
	server.Any("/appproxy/v2", func(context *gin.Context) {

		args := context.Request.URL.Query()
		argBytes,_ := json.Marshal(args)
		headers := context.Request.Header

		headerBytes, _ := json.Marshal(headers)
		reqBody := context.Request.Body
		body, _ := ioutil.ReadAll(reqBody)
		defer reqBody.Close()
		params := context.Params
		paramsByte, _ := json.Marshal(params)
		cookies := context.Request.Cookies()
		cookiesByte, _ := json.Marshal(cookies)
		fmt.Printf("请求url:%s,\t请求方法：%s,\t请求headers：%s,请求args:%s\t请求cookie：%s,\t请求：params:%s,\t请求body：%s \n",
			context.Request.URL.Path, context.Request.Method, string(headerBytes),string(argBytes), string(cookiesByte), string(paramsByte), string(body))
		resp := struct {
			Code   int64
			Msg    string
			Result interface{}
		}{
			Code:   232232,
			Msg:    context.Param("id"),
			Result: []string{string(headerBytes), string(cookiesByte), string(paramsByte), string(body)},
		}

		context.JSON(http.StatusOK, resp)
	})
	if err := server.Run(fmt.Sprintf(":%d", 8080)); err != nil {
		os.Exit(0x1)
	}
}


