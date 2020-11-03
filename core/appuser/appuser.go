package appuser

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"../common"
	"../common/zlog"
	"../conf"
	"../model"
)

var AppUserMgr = &AppUserManager{}

type AppUserManager struct {
	users sync.Map
	dirty bool
}

func (aum *AppUserManager) Name() string {
	return "AppUserManager"
}
func (aum *AppUserManager) Init() bool {
	aum.readFile()
	sm := http.NewServeMux()
	if len(conf.Conf.AppService) > 0 {
		return false
	}
	sm.HandleFunc("/registeapp", aum.registeapp)
	l, err := net.Listen("tcp", conf.Conf.AppService)
	if err != nil {
		zlog.Error("AppUserManager Listen", zlog.String("Err", err.Error()))
	}
	go func() {
		err := http.Serve(l, sm)
		if err != nil {
			zlog.Error("AppUserManager server", zlog.String("Err", err.Error()))
		}
	}()

	return true
}
func (aum *AppUserManager) Update() bool {
	if aum.dirty {
		aum.writeFile()
	}
	return true
}
func (aum *AppUserManager) End() bool {
	return true
}
func (aum *AppUserManager) readFile() {
	file, err := os.Open("app.data")
	if err != nil {
		zlog.Error("AppUserManager write data", zlog.String("Err", err.Error()))
	} else {
		data := []*model.AppUser{}
		var buff []byte
		_, err = file.Read(buff)
		err := json.Unmarshal(buff, data)
		if err != nil {
			zlog.Error("Marshal appuser", zlog.String("Err", err.Error()))
		} else {
			for _, value := range data {
				aum.users.Store(value.AppId, value)
			}
		}
		err = file.Close()
		if err != nil {
			zlog.Error("Close appuser data", zlog.String("Err", err.Error()))
		}
	}
}
func (aum *AppUserManager) writeFile() {
	file, err := os.Open("app.data")
	if err != nil {
		zlog.Error("AppUserManager write data", zlog.String("Err", err.Error()))
	} else {
		data := []*model.AppUser{}
		aum.users.Range(func(key, value interface{}) bool {
			data = append(data, value.(*model.AppUser))
			return true
		})
		buff, err := json.Marshal(data)
		if err != nil {
			zlog.Error("Marshal appuser", zlog.String("Err", err.Error()))
		} else {
			_, err = file.Write(buff)
			if err != nil {
				zlog.Error("Write appuser", zlog.String("Err", err.Error()))
			}
		}
		err = file.Close()
		if err != nil {
			zlog.Error("Close appuser data", zlog.String("Err", err.Error()))
		}
	}
}
func (aum *AppUserManager) registeapp(rw http.ResponseWriter, r *http.Request) {
	appid, appkey, seed, ts := common.CreateApp()
	aum.users.Store(appid, &model.AppUser{
		AppId:  appid,
		AppKey: appkey,
		Seed:   seed,
		Ts:     ts,
	})
	rw.Write([]byte(fmt.Sprintf("{AppId:%s,AppKey:%s}", appid, appkey)))
	aum.dirty = true
}
