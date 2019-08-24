package common

import (
	"fmt"
	"os"

	"github.com/astaxie/beego/logs"
	"gopkg.in/ini.v1"
)

var Cfg *ini.File

func LoadConf(confPath string) error {
	isExist, err2 := PathExists(confPath)
	if err2 != nil {
		logs.Error(err2)
		return err2
	}
	if !isExist {

	}
	var err error
	Cfg, err = ini.Load(confPath)
	if err != nil || Cfg == nil {
		err := fmt.Errorf("can not load configure file:%s %s ", confPath, err)
		logs.Error(err)
		return err
	}
	return nil
}

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
