package main

import (
	"fmt"
	"os"
	"strconv"

	"encoding/csv"

	"./common"
	"github.com/astaxie/beego/logs"
)

var (
	csvName = `./addRecord.csv`
)

func main() {
	err := common.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Create(csvName)
	if err != nil {
		logs.Error(err)
		return
	}
	if file == nil {
		logs.Error(err)
		return
	}
	write := csv.NewWriter(file)

	csvData := make([][]string, 0)
	ctitle := make([]string, 2)
	ctitle[0] = `序号`
	ctitle[1] = `请求数据`
	csvData = append(csvData, ctitle)
	for i := 0; i < 1000; i++ {
		c := make([]string, 2)
		c[0] = strconv.Itoa(i + 1)
		c[1] = `{"action":"","data":{"meeting_id":7,"uid":` + strconv.Itoa(i+1) + `,"create_time":"0001-01-01T00:00:00Z","unit_name":"广大","name":"宙斯","job":"老师","mobile_phone":"15555555555","email":"777777777@qq.com","space":"11号天桥"},"sets":null,"orderBy":"","filter":"","page":0,"pageSize":0}`
		csvData = append(csvData, c)
	}
	err = write.WriteAll(csvData)
	if err != nil {
		logs.Error(err)
		return
	}
}
