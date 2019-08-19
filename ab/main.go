package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"./common"
	"./http"
	"github.com/astaxie/beego/logs"
)

var wg sync.WaitGroup

var (
	csvName       = `./addRecord.csv`
	csvResultName = `./addRecord_result.csv`
	URL           = `http://localhost:8083/addRecord`
)

func main() {
	err := common.InitLogger()
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Open(csvName)
	defer file.Close()
	if err != nil {
		logs.Error(err)
		return
	}
	if file == nil {
		err = fmt.Errorf(`file is null`)
		logs.Error(err)
		return
	}
	reader := csv.NewReader(file)
	datas, err := reader.ReadAll()
	if err != nil {
		logs.Error(err)
		return
	}
	if len(datas) == 0 || datas == nil {
		err = fmt.Errorf(`datas is null`)
		logs.Error(err)
		return
	}
	getData := make(chan []string, 300)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for i := 0; i < 100; i++ {
		data := datas[i+1]
		wg.Add(1)
		go func(d []string, index int, ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					replyBody, since, err := http.HttpPostJson(URL, d[1])
					if err != nil {
						logs.Error(err)
						return
					}
					d = append(d, replyBody)
					d = append(d, since.String())
					getData <- d
				}
			}
		}(data, i, ctx)
	}
	errChan, canStop := HandleResult(getData, csvResultName)
	wg.Wait()
	close(getData)
	select {
	case err = <-errChan:
		logs.Error(err)
	case <-canStop:
	}
}

// ResultCsvName包含路径名，如 ./sendEmailCode.csv
func HandleResult(data <-chan []string, ResultCsvName string) (errChan chan error, canStop chan bool) {
	canStop = make(chan bool)
	errChan = make(chan error)
	if data == nil {
		err := fmt.Errorf(`data is null`)
		logs.Error(err)
		errChan <- err
		return
	}
	if ResultCsvName == `` {
		err := fmt.Errorf(`ResultCsvName is null`)
		logs.Error(err)
		errChan <- err
		return
	}
	go func() {
		result := make([][]string, 0)
		for v := range data {
			result = append(result, v)
		}
		file, err := os.Create(ResultCsvName)
		if err != nil {
			logs.Error(err)
			errChan <- err
			return
		}
		write := csv.NewWriter(file)
		err = write.WriteAll(result)
		if err != nil {
			logs.Error(err)
			errChan <- err
			return
		}
		canStop <- true
	}()
	return
}
