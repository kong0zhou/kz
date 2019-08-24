package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"./common"
	"github.com/astaxie/beego/logs"
)

var wg sync.WaitGroup

var (
	csvResultName       = `./addRecord_result.csv` //输出结果的名字
	URL                 = `http://cst.gzhu.edu.cn:6635/addRecord`
	Concurrency         = 125 //并发数
	timeMinute    int64 = 30  //持续时间（分钟）
)

func main() {
	err := common.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	err = common.LoadConf(path + `/conf/conf.ini`)
	if err != nil {
		logs.Error(err)
		return
	}
}
