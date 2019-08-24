package common

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
)

type JOSNData struct {
	Data chan string
}

// n是指生成的UID的范围[1,n)
func NewJSONData(n int) (*JOSNData, error) {
	if n < 1 {
		err := fmt.Errorf(`n is not rightv %d`, n)
		logs.Error(err)
		return nil, err
	}
	jo := JOSNData{}
	jo.Data = make(chan string, 100)
	res := rand.New(rand.NewSource(time.Now().UnixNano()))
	go func() {
		for {
			l := res.Intn(n-1) + 1
			str := `{"action":"","data":{"meeting_id":2,"uid":` + strconv.Itoa(l) + `,"create_time":"0001-01-01T00:00:00Z","unit_name":"广大","name":"宙斯","job":"老师","mobile_phone":"15555555555","email":"777777777@qq.com","space":"11号天桥"},"sets":null,"orderBy":"","filter":"","page":0,"pageSize":0}`
			jo.Data <- str
		}
	}()
	return &jo, nil
}

func (t *JOSNData) GetJOSNData() (json string, err error) {
	if t.Data == nil {
		err = fmt.Errorf(`t.Data is null`)
		logs.Error(err)
		return "", err
	}
	json = <-t.Data
	return json, nil
}
