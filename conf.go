package kz

import (
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	ConfigFile string `mapstructure:"configFile"` //配置本地的配置文件的文件路径，默认用这个
	Production bool   `mapstructure:"production"` //当前环境是不是生产环境
}

// Config 存储所有的配置
var Config config

// InitConf 初始化配置文件
func InitConf() (err error) {
	v := viper.New()

	// v.BindEnv(`PASSWORD`)
	// v.BindEnv(`DIRPATH`)
	// v.BindEnv(`UID`)
	// v.BindEnv(`SESSIONKEY`)

	pflag.StringP(`configFile`, `f`, ``, `please input your config file path`)
	pflag.BoolP(`production`, `p`, false, `please input true if is production`)
	pflag.Parse()
	err = v.BindPFlags(pflag.CommandLine)
	if err != nil {
		return err
	}
	configFile := v.GetString(`configFile`)

	// 从指定的文件中读取配置文件
	if configFile != `` {
		isExists, err := PathExists(configFile)
		if err != nil {
			return err
		}
		if !isExists {
			err = fmt.Errorf(`path is not exists`)
			return err
		}
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigFile(`./config.json`)
		logs.Info("使用当前目录下的config.json为配置文件")
	}

	err = v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&Config)
	if err != nil {
		return err
	}
	return nil
}
