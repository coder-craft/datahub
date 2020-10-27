package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"time"

	"./common/zlog"
	"./device"
	"./module"
	"./mysql"
	"./remoteuser"
	"./web"
)

// server
func Start() {
	defer func() {
		if err := recover(); err != nil {
			filename := fmt.Sprintf("im%v.dump", time.Now().Format("200601021504050700"))
			err := ioutil.WriteFile(filename, debug.Stack(), os.ModePerm)
			if err != nil {
				zlog.Error("Save dump file", zlog.String("Err",err.Error()))
			}
		}
	}()
	err := remoteuser.InitUser()
	if err != nil {
		zlog.Error("Init user",zlog.String("Err",err.Error()))
		return
	}
	_ = mysql.InitMysqlDB()
	web.StartHttpService()
	//module
	module.RegisterModule(device.DeviceMgr,time.Second)
	module.ModuleStart()
}
func Stop() {
	mysql.CloseMysqlDB()
	module.ModuleStop()
}
