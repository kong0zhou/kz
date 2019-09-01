package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"
)

func ErrorHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logs.Error(err)
				var reply ReplyProto
				reply.Status = -1
				reply.Msg = fmt.Sprint(err)
				response, err := json.Marshal(reply)
				if err != nil {
					logs.Error(err)
					return
				}
				_, err = w.Write(response)
				if err != nil {
					logs.Error(err)
					return
				}
			}
		}()
		h(w, r)
	}
}
