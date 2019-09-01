package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"
)

//后端响应数据通信协议
type ReplyProto struct {
	Status   int         `json:"status"` //状态 0正常，小于0出错，大于0可能有问题
	Msg      string      `json:"msg"`    //状态信息
	Data     interface{} `json:"data"`
	API      string      `json:"API"`      //api接口
	Method   string      `json:"method"`   //post,put,get,delete
	SN       string      `json:"SN"`       //标识符
	RowCount int         `json:"rowCount"` //Data若是数组，算其长度
}

func NewReplyProto(method string, API string) (*ReplyProto, error) {
	if method == `` {
		err := fmt.Errorf(`method is null`)
		logs.Error(err)
		return nil, err
	}
	if API == `` {
		err := fmt.Errorf(`API is null`)
		logs.Error(err)
		return nil, err
	}
	t := ReplyProto{}
	t.Method = method
	t.API = API
	return &t, nil
}

func (t *ReplyProto) ErrorResp(errMsg string, w http.ResponseWriter) (err error) {
	if t == nil {
		err = fmt.Errorf(`replyProto is null`)
		logs.Error(err)
		return
	}
	if w == nil {
		err = fmt.Errorf("arguments can not be a nil value")
		logs.Error(err)
		return err
	}
	t.Status = -1
	t.Msg = errMsg
	t.RowCount = 0
	response, err := json.Marshal(&t)
	if err != nil {
		logs.Error(err)
		return err
	}
	_, err = w.Write(response)
	if err != nil {
		logs.Error(err)
		return err
	}
	return nil
}

func (t *ReplyProto) SuccessResp(data interface{}, w http.ResponseWriter) (err error) {
	if t == nil {
		err = fmt.Errorf(`replyProto is null`)
		logs.Error(err)
		return
	}
	if w == nil {
		err = fmt.Errorf("arguments can not be a nil value")
		logs.Error(err)
		return err
	}
	t.Status = 0
	t.Data = data
	t.RowCount = 0
	response, err := json.Marshal(&t)
	if err != nil {
		return err
	}
	_, err = w.Write(response)
	if err != nil {
		return err
	}
	return nil
}

func (t *ReplyProto) DefinedResp(status int, msg string, data interface{}, w http.ResponseWriter) (err error) {
	if w == nil {
		err = fmt.Errorf("arguments can not be a nil value")
		logs.Error(err)
		return err
	}
	if t == nil {
		err = fmt.Errorf(`replyProto is null`)
		logs.Error(err)
		return
	}
	t.Status = status
	t.Msg = msg
	t.Data = data
	t.RowCount = 0
	response, err := json.Marshal(&t)
	if err != nil {
		return err
	}
	_, err = w.Write(response)
	if err != nil {
		return err
	}
	return nil
}

//前端请求数据通讯协议
type ReqProto struct {
	Action   string      `json:"action"` //请求类型GET/POST/PUT/DELETE
	Data     interface{} `json:"data"`   //请求数据
	Sets     []string    `json:"sets"`
	OrderBy  string      `json:"orderBy"`  //排序要求
	Filter   string      `json:"filter"`   //筛选条件
	Page     int         `json:"page"`     //分页
	PageSize int         `json:"pageSize"` //分页大小
}
