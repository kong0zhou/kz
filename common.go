package kz

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// PathExists 查看相应的路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetClientIP 获取客户端的真实IP
func GetClientIP(r *http.Request) (ip string, err error) {
	if r == nil {
		err = fmt.Errorf(`*http.Request is nil`)
		log.Error(err)
		return ``, err
	}
	ip = r.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == `` {
		ip = r.Header.Get("X-real-ip")
	}
	if ip == `` {
		return "127.0.0.1", nil
	}
	return ip, nil
}

// ExecCommand 执行命令，执行过程中函数阻塞
func ExecCommand(commandName string, params []string) (result string, err error) {
	if commandName == `` {
		err = fmt.Errorf("commandName is null")
		log.Error(err)
		return ``, err
	}
	if params == nil {
		err = fmt.Errorf(`params is nil`)
		log.Error(err)
		return ``, err
	}
	cmd := exec.Command(commandName, params...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf(err.Error() + `:` + stderr.String())
		return ``, err
	}
	result = out.String()
	return
}

// Underline2hump 名字转化，将下划线转驼峰
func Underline2hump(s string) (r string) {
	if s == `` {
		return ``
	}
	if s == `id` {
		return `ID`
	}
	ss := strings.Split(s, `_`)
	for _, v := range ss {
		if v == `` {
			continue
		}
		if v == `id` {
			r = r + `ID`
			continue
		}
		r = r + strings.Title(v)
	}
	return r
}
