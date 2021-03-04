package controller

import (
	"fmt"
	"github.com/astaxie/beego/toolbox"
)
var task *toolbox.Task

func init() {

	 task = toolbox.NewTask("testChangeSpec","*/5 * * * * *",testChangeSpec)
	 toolbox.AddTask("testChangeSpec",task)
	 task.Run()

}

func testChangeSpec() error {
    var err error
	fmt.Println("test change spec")
    return err
}
