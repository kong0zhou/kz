package kz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// ErrorHandlerFunc 错误处理中间件,panic时调用
func ErrorHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Panic(err)
				var reply ReplyProto
				reply.Status = -1
				reply.Msg = fmt.Sprint(err)
				response, err := json.Marshal(reply)
				if err != nil {
					log.Error(err)
					return
				}
				_, err = w.Write(response)
				if err != nil {
					log.Error(err)
					return
				}
			}
		}()
		h(w, r)
	}
}

// ErrorHandler 错误处理中间件,panic时调用
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Panic(err)
				var reply ReplyProto
				reply.Status = -1
				reply.Msg = "后端出现panic,请及时和后端程序员反映这个错误：" + fmt.Sprint(err)
				response, err := json.Marshal(reply)
				if err != nil {
					log.Error(err)
					return
				}
				_, err = w.Write(response)
				if err != nil {
					log.Error(err)
					return
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Logger 日志中间件，当有人访问时打印其信息
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		endTime := time.Now()
		//执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := r.Method
		// 请求路由
		reqURI := r.RequestURI
		clientIP, err := GetClientIP(r)

		if err != nil {
			log.Error(err)
			var reply ReplyProto
			reply.Status = -1
			reply.Msg = "logger中间件出现错误：" + fmt.Sprint(err)
			response, err := json.Marshal(reply)
			if err != nil {
				log.Error(err)
				return
			}
			_, err = w.Write(response)
			if err != nil {
				log.Error(err)
				return
			}
		}
		log.WithFields(logrus.Fields{
			"latencyTime": latencyTime,
			"requMethod":  r.Method,
			"reqURI":      reqURI,
			"reqMethod":   reqMethod,
			"clientIP":    clientIP,
		}).Info("")
	})
}
