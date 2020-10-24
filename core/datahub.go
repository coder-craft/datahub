package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"time"

	"./common/zlog"
	"./module"
	"./mysql"
	"./web"
)

// server
func Start() {
	defer func() {
		if err := recover(); err != nil {
			filename := fmt.Sprintf("im%v.dump", time.Now().Format("200601021504050700"))
			err := ioutil.WriteFile(filename, debug.Stack(), os.ModePerm)
			if err != nil {
				zlog.Errorf("Save dump file err:", err.Error())
			}
		}
	}()
	_ = mysql.InitMysqlDB()
	web.StartHttpService()
	//module
	module.ModuleStart()
}
func Stop() {
	mysql.CloseMysqlDB()
	module.ModuleStop()
}
