package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"../https"
	"../tools"

	"github.com/astaxie/beego/logs"
)

func Deploy(w http.ResponseWriter, r *http.Request) {
	reply, err := NewReplyProto(`PUT`, `/deploy`)
	if err != nil {
		logs.Error(err)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		return
	}
	if body == nil || len(body) == 0 {
		err = fmt.Errorf(`body is null`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	var req ReqProto
	err = json.Unmarshal(body, &req)
	if err != nil {
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	data, ok := req.Data.(map[string]interface{})
	if !ok {
		err = fmt.Errorf(`req.Data is not map[string]interface {} do not exists`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	concurrent, ok := data[`concurrent`].(float64)
	if !ok {
		err = fmt.Errorf(`req.Data.concurrent is not float64 or do not exists`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	if concurrent <= 0 {
		err := fmt.Errorf(`concurrent is not right`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	method, ok := data[`method`].(string)
	if !ok {
		err = fmt.Errorf(`req.Data.method is not string or do not exists`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	if method == `` {
		err = fmt.Errorf(`req.Data.method is not right`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	reqData, ok := data[`reqData`].(string)
	if !ok {
		err = fmt.Errorf(`req.Data.reqData is not string or do not exists`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	if reqData == `` {
		err = fmt.Errorf(`req.Data.reqData is not right`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	filePath, ok := data[`filePath`].(string)
	if !ok {
		err = fmt.Errorf(`req.Data.filePath is not string or do not exists`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	if filePath == `` {
		err = fmt.Errorf(`req.Data.filePath is not right`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	ti, ok := data[`time`].(float64)
	if !ok {
		err = fmt.Errorf(`req.Data.time is not string or do not exists`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	if ti <= 0 {
		err = fmt.Errorf(`req.Data.time is not right`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	url, ok := data[`url`].(string)
	if !ok {
		err = fmt.Errorf(`req.Data.url is not string or do not exists`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}
	if url == `` {
		err = fmt.Errorf(`req.Data.url is not right`)
		logs.Error(err)
		err = reply.ErrorResp(err.Error(), w)
		if err != nil {
			logs.Error(err)
			return
		}
		return
	}

	jsonData, err := tools.NewJSONData(reqData)
	if err != nil {
		logs.Error(err)
		return
	}
	dispatch, err := tools.NewDispatch(https.HttpPostJson, int(concurrent), time.Duration(ti)*time.Second, jsonData.Data, url)
	if err != nil {
		logs.Error(err)
		return
	}
	err = dispatch.Begin()
	if err != nil {
		logs.Error(err)
		return
	}
	err = reply.SuccessResp(nil, w)
	if err != nil {
		logs.Error(err)
		return
	}
}
