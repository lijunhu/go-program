package test

import (
	"fmt"
	"runtime"
	"testing"
)

func TestPanic(t *testing.T) {
	ProtectFunc(func() {
		panic("测试panic")
	})
}

func ProtectFunc(target func()) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case runtime.Error:
				{
					fmt.Println("运行时错误", err)
				}
			default:
				{
					fmt.Println("error:", err)
				}
			}
		}
	}()
	target()
}
