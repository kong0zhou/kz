package http

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

func HttpPostJson(url string, json string) (replyBody string, since time.Duration, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(`POST`, url, strings.NewReader(json))
	if err != nil {
		logs.Error(err)
		return ``, 0, err
	}

	req.Header.Set(`Content-Type`, `application/json`)
	start := time.Now()
	resp, err := client.Do(req)
	since = time.Since(start)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error(err)
		return ``, 0, err
	}
	replyBody = string(body)
	return replyBody, since, nil
}
