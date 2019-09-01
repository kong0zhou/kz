package tools

import (
	"github.com/astaxie/beego/logs"
	regen "github.com/zach-klippenstein/goregen"
)

type JOSNData struct {
	Data chan string
}

// n是指生成的UID的范围[1,n)
func NewJSONData(pattern string) (*JOSNData, error) {
	jo := JOSNData{}
	jo.Data = make(chan string, 300)
	go func() {
		for {
			str, err := regen.Generate(pattern)
			if err != nil {
				logs.Error(err)
				return
			}
			jo.Data <- str
		}
	}()
	return &jo, nil
}
