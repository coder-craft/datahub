package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"time"

	"../common/zlog"
)

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}

	return exist
}

func SafeCall(callback func(args ...interface{}) error, iargs ...interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			if err := recover(); err != nil {
				filename := fmt.Sprintf("im%v.dump", time.Now().Format("200601021504050700"))
				err := ioutil.WriteFile(filename, debug.Stack(), os.ModePerm)
				if err != nil {
					zlog.Error("Save dump file", zlog.String("Err", err.Error()))
				}
			}

		}
	}()
	return callback(iargs...)
}
