package tools

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/astaxie/beego/logs"
)

// var wg sync.WaitGroup

type HandleResult struct {
	ErrorResult    [][]string    //每个请求的数据，只存放错误的数据   序号，前端输入数据，后端输出数据，前端响应时间，后端响应时间
	FrontAveTime   time.Duration //前端平均响应时间
	BackAveTime    time.Duration //后端平均响应时间
	FrontFrequency int           //前端请求次数
	BackFrequency  int           //后端处理请求次数
	ErrorFrequency int           //出现错误的次数

	AllTime     time.Duration //压力测试总时间，这个需要自己输入
	Concurrency int           //并发量

	// ===========可以推出来的数据=============
	Throughput       float64       //后端吞吐量，以分钟为单位
	FailureFrequency int           //失败次数
	NetworkDelay     time.Duration //网络延迟

	// 私有变量
	frontAllTime   time.Duration      //前端所有检测时间之和
	backAllTime    time.Duration      //后端所有检测时间之和
	frontChanTime  chan time.Duration //前端的时间，通过管道运输
	endChanTime    chan time.Duration //后端时间，通过管道运输
	resultChanData chan []string      //运输的内容为  前端输入数据，后端输出数据，前端响应时间，后端响应时间

	wg sync.WaitGroup
}

func NewHandleResult(alltime time.Duration, concurrency int) *HandleResult {
	h := HandleResult{}
	h.ErrorResult = make([][]string, 0)

	h.resultChanData = make(chan []string, 300)
	h.frontChanTime = make(chan time.Duration, 300)
	h.endChanTime = make(chan time.Duration, 300)

	h.AllTime = alltime
	h.Concurrency = concurrency
	return &h
}

func (t *HandleResult) InputFrontTime(time time.Duration) (err error) {
	if time <= 0 {
		err = fmt.Errorf(`time is not right`)
		logs.Error(err)
		return err
	}
	if t.frontChanTime == nil {
		err = fmt.Errorf(`t.frontChanTime is null`)
		logs.Error(err)
		return err
	}
	t.frontChanTime <- time
	return nil
}

func (t *HandleResult) InputBackTime(time time.Duration) (err error) {
	if time <= 0 {
		err = fmt.Errorf(`time is not right`)
		logs.Error(err)
		return err
	}
	if t.endChanTime == nil {
		err = fmt.Errorf(`t.endChanTime is null`)
		logs.Error(err)
		return err
	}
	t.endChanTime <- time
	return nil
}

func (t *HandleResult) InputResult(r []string) (err error) {
	if r == nil || len(r) == 0 {
		err = fmt.Errorf(`result is not right`)
		logs.Error(err)
		return
	}
	if t.resultChanData == nil {
		err = fmt.Errorf(`t.resultChanData is null`)
		logs.Error(err)
		return
	}
	t.resultChanData <- r
	return nil
}

func (t *HandleResult) Handle() (err error) {
	if t.resultChanData == nil {
		err = fmt.Errorf(`t.resultChanData is null`)
		logs.Error(err)
		return
	}
	if t.endChanTime == nil {
		err = fmt.Errorf(`t.endChanTime is null`)
		logs.Error(err)
		return err
	}
	if t.frontChanTime == nil {
		err = fmt.Errorf(`t.frontChanTime is null`)
		logs.Error(err)
		return err
	}
	// 处理Result
	t.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logs.Error(err)
				return
			}
		}()
		defer t.wg.Done()
		var reply ReplyProto
		for v := range t.resultChanData {
			err := json.Unmarshal([]byte(v[1]), &reply)
			if err != nil {
				logs.Info(v[1])
				logs.Error(err)
				continue
			}
			if reply.Status != 0 {
				t.ErrorFrequency++
				v = append([]string{strconv.Itoa(t.ErrorFrequency)}, v...)
				t.ErrorResult = append(t.ErrorResult, v)
				logs.Info(v)
			}
		}
	}()
	//
	t.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logs.Error(err)
				return
			}
		}()
		defer t.wg.Done()
		for ti := range t.frontChanTime {
			t.FrontFrequency++
			t.frontAllTime = t.frontAllTime + ti
		}
		if t.FrontFrequency == 0 {
			return
		}
		t.FrontAveTime = time.Duration(int64(t.frontAllTime) / int64(t.FrontFrequency))
	}()
	t.wg.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logs.Error(err)
				return
			}
		}()
		defer t.wg.Done()
		for ti := range t.endChanTime {
			t.BackFrequency++
			t.backAllTime = t.backAllTime + ti
		}
		if t.BackFrequency == 0 {
			return
		}
		t.BackAveTime = time.Duration(int64(t.backAllTime) / int64(t.BackFrequency))
	}()
	// t.wg.Wait()
	return
}

// 关闭所有管道，这是必须操作
func (t *HandleResult) CloseChan() error {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
			return
		}
	}()
	if t.resultChanData == nil {
		err := fmt.Errorf(`t.resultChanData is null`)
		logs.Error(err)
		return err
	}
	if t.frontChanTime == nil {
		err := fmt.Errorf(`t.frontChanTime is null`)
		logs.Error(err)
		return err
	}
	if t.endChanTime == nil {
		err := fmt.Errorf(`t.endChanTime is null`)
		logs.Error(err)
		return err
	}
	close(t.resultChanData)
	close(t.frontChanTime)
	close(t.endChanTime)
	return nil
}

// 等待所有线程完成
func (t *HandleResult) Wait() {
	t.wg.Wait()
}

func (t *HandleResult) GenerateErrorCSV(csvName string) (err error) {
	if csvName == `` {
		err = fmt.Errorf(`csvName is null`)
		logs.Error(err)
		return err
	}
	if t.ErrorResult == nil {
		err = fmt.Errorf(`t.ErrorResult is null`)
		logs.Error(err)
		return err
	}
	file, err := os.Create(csvName)
	if err != nil {
		logs.Error(err)
		return
	}
	write := csv.NewWriter(file)
	err = write.WriteAll(t.ErrorResult)
	if err != nil {
		logs.Error(err)
		return
	}
	return nil
}

func (t *HandleResult) GetResult() (err error) {
	if t.BackFrequency <= 0 {
		err = fmt.Errorf(`backFrequency is not right`)
		logs.Error(err)
		return
	}
	if t.AllTime <= 0 {
		err = fmt.Errorf(`t.AllTime is not right`)
		logs.Error(err)
		return
	}
	if t.FrontFrequency <= 0 {
		err = fmt.Errorf(`t.FrontFrequency is not right`)
		logs.Error(err)
		return
	}
	if t.FrontAveTime <= 0 {
		err = fmt.Errorf(`t.FrontAveTime is not right`)
		logs.Error(err)
		return
	}
	if t.BackAveTime <= 0 {
		err = fmt.Errorf(`t.BackAveTime is not right`)
		logs.Error(err)
		return
	}
	t.Throughput = float64(t.BackFrequency) / t.AllTime.Minutes()
	t.FailureFrequency = t.FrontFrequency - t.BackFrequency
	t.NetworkDelay = t.FrontAveTime - t.BackAveTime
	return nil
}

func (t *HandleResult) GenerateResultMD(mdName string) (err error) {
	if mdName == `` {
		err = fmt.Errorf(`mdName is null`)
		logs.Error(err)
		return err
	}
	err = t.GetResult()
	if err != nil {
		logs.Error(err)
		return err
	}
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logs.Error(err)
		return
	}
	temp, err := template.ParseFiles(path + `/addition_file/result.md`)
	if err != nil {
		logs.Error(err)
		return
	}
	file, err := os.Create(mdName)
	if err != nil {
		logs.Error(err)
		return
	}
	err = temp.Execute(file, &t)
	if err != nil {
		logs.Error(err)
		return
	}
	return nil
}

// 实时打印数据，ti用来设定打印间隔时间
func (t *HandleResult) LogPrint(ctx context.Context, ti time.Duration) (err error) {
	if ti <= 0 {
		err = fmt.Errorf(`ti is not right`)
		logs.Error(err)
		return err
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logs.Error(err)
				return
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(ti)
				logs.Info(`=========================`)
				logs.Info(`前端检测次数`, t.FrontFrequency)
				if t.FrontFrequency == 0 {
					logs.Info(`t.FrontFrequency is zero`)
					return
				}
				logs.Info(`前端平均检测时间`, time.Duration(int64(t.frontAllTime)/int64(t.FrontFrequency)))
				logs.Info(`后端响应次数`, t.BackFrequency)
				logs.Info(`后端平均检测时间`, time.Duration(int64(t.backAllTime)/int64(t.BackFrequency)))
				if t.BackFrequency == 0 {
					logs.Info(`t.BackFrequency is zero`)
					return
				}
				logs.Info(`=========================`)
			}
		}
	}()
	return nil
}
