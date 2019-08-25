package common

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
)

type tFunc func(url string, json string) (replyBody string, since time.Duration, err error)

type Dispatch struct {
	HttpMethod  tFunc         //方法，调度器统一调用这个方法
	Concurrency int           //并发量
	AllTime     time.Duration // 压力测试的时间
	RepJsonData chan string   //一个对外的管道,由用户将自己要发送的json数据传过来
	URL         string

	handle *HandleResult //处理结果
	wg     sync.WaitGroup
}

func NewDispatch(f tFunc, concurrency int, allTime time.Duration, repJsonData chan string, URL string) (*Dispatch, error) {
	if concurrency <= 0 {
		err := fmt.Errorf(`concurrency is not right`)
		logs.Error(err)
		return nil, err
	}
	if allTime <= 0 {
		err := fmt.Errorf(`allTime is not right`)
		logs.Error(err)
		return nil, err
	}
	if repJsonData == nil {
		err := fmt.Errorf(`repJsonData is not right`)
		logs.Error(err)
		return nil, err
	}
	if URL == `` {
		err := fmt.Errorf(`url is null`)
		logs.Error(err)
		return nil, err
	}
	d := Dispatch{}
	d.HttpMethod = f
	d.Concurrency = concurrency
	d.RepJsonData = repJsonData
	d.AllTime = allTime
	d.handle = NewHandleResult(d.AllTime, d.Concurrency)
	d.URL = URL
	return &d, nil
}

func (t *Dispatch) Begin() (err error) {
	if t.Concurrency <= 0 {
		err = fmt.Errorf(`t.Concurrency is not right'`)
		logs.Error(err)
		return
	}
	if t.handle == nil {
		err = fmt.Errorf(`t.handle is null`)
		logs.Error(err)
		return
	}
	if t.RepJsonData == nil {
		err = fmt.Errorf(`t.RepJsonData is null`)
		logs.Error(err)
		return
	}
	if t.AllTime <= 0 {
		err = fmt.Errorf(`t.AllTime is not right`)
		logs.Error(err)
		return
	}
	if t.URL == `` {
		err = fmt.Errorf(`url is null`)
		logs.Error(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.AllTime)
	defer cancel()
	for i := 0; i < t.Concurrency; i++ {
		t.wg.Add(1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					logs.Error(err)
					return
				}
			}()
			defer t.wg.Done()
			for {
				var reply ReplyProto
				select {
				case <-ctx.Done():
					return
				case jsonData := <-t.RepJsonData:
					replyBody, since, err := t.HttpMethod(t.URL, jsonData)
					if err != nil {
						logs.Error(err)
						continue
					}
					if since <= 0 {
						err := fmt.Errorf(`since is not right`)
						logs.Error(err)
						return
					}
					if replyBody == `` {
						err := fmt.Errorf(`replyBody is null`)
						logs.Error(err)
						return
					}
					err = t.handle.InputFrontTime(since)
					if err != nil {
						logs.Error(err)
						return
					}
					result := make([]string, 4)
					result[0] = jsonData
					result[1] = replyBody
					result[2] = since.String()
					err = json.Unmarshal([]byte(replyBody), &reply)
					if err != nil {
						logs.Error(err)
						continue
					}
					result[3] = reply.BackSinceTime.String()
					err = t.handle.InputBackTime(reply.BackSinceTime)
					if err != nil {
						logs.Error(err)
						return
					}
					err = t.handle.InputResult(result)
					if err != nil {
						logs.Error(err)
						return
					}
				}
			}
		}()
	}
	err = t.handle.Handle()
	if err != nil {
		logs.Error(err)
		return
	}
	err = t.handle.LogPrint(500 * time.Millisecond)
	if err != nil {
		logs.Error(err)
		return
	}
	t.wg.Wait()
	err = t.handle.CloseChan()
	if err != nil {
		logs.Error(err)
		return
	}
	t.handle.Wait()
	err = t.handle.GenerateErrorCSV(`./error.csv`)
	if err != nil {
		logs.Error(err)
		return
	}
	err = t.handle.GenerateResultMD(`./result.md`)
	if err != nil {
		logs.Error(err)
		return
	}
	return
}
