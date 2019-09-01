package main

import (
	"fmt"
	"net/http"

	"./controllers"
	"./tools"
	"github.com/astaxie/beego/logs"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			return
		}
	}()
	err := tools.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/deploy", controllers.ErrorHandler(controllers.Deploy))
	logs.Info("http服务器启动，端口：8083")
	err = http.ListenAndServe(":8083", mux)
	if err != nil {
		logs.Error("启动失败", err)
	}
}
