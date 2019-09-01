package tools

import (
	"encoding/json"
	"log"

	"github.com/astaxie/beego/logs"
)

func InitLogger() error {

	// // 获取运行目录
	// dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	// if err != nil {
	// 	log.Fatal("获取运行目录失败 ", err)
	// 	return err
	// }
	// //在运行目录下创建日志目录
	// err = os.MkdirAll(dir+"/log/", os.ModePerm)
	// if err != nil {
	// 	log.Fatal("创建日志目录失败 ", err)
	// 	return err
	// }

	logs.SetLogFuncCallDepth(3)    //调用层级
	logs.EnableFuncCallDepth(true) //输出文件名和行号
	// logs.Async()                   //提升性能, 可以设置异步输出

	config := make(map[string]interface{})
	config["filename"] = `./logs/log.log`
	config["level"] = logs.LevelDebug
	config["perm"] = "0777"

	configStr, err := json.Marshal(config)
	if err != nil {
		log.Fatal("initLogger failed, marshal err:", err)
		return err
	}

	err = logs.SetLogger(logs.AdapterConsole, "")             //控制台输出
	err = logs.SetLogger(logs.AdapterFile, string(configStr)) //文件输出

	if err != nil {
		log.Fatal("SetLogger failed, err:", err)
		return err
	}
	return nil
}
