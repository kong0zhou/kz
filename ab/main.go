package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"./common"
	"./http"
	"github.com/astaxie/beego/logs"
)

var wg sync.WaitGroup

var (
	csvResultName       = `./addRecord_result.csv` //输出结果的名字
	URL                 = `http://localhost:8083/export`
	Concurrency         = 1  //并发数
	timeMinute    int64 = 30 //持续时间（分钟）
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			return
		}
	}()
	err := common.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.Error(err)
		return
	}
	err = common.LoadConf(path + `/conf/conf.ini`)
	if err != nil {
		logs.Error(err)
		return
	}
	jsonData, err := common.NewJSONData(100000)
	if err != nil {
		logs.Error(err)
		return
	}
	dispatch, err := common.NewDispatch(http.HttpPostJson, Concurrency, 3*time.Second, jsonData.Data, URL)
	if err != nil {
		logs.Error(err)
		return
	}
	err = dispatch.Begin()
	if err != nil {
		logs.Error(err)
		return
	}
}
