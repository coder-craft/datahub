package web

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"../common/zlog"
	"../conf"
	"../model"
	"../remoteuser"
)

func init() {
	http.HandleFunc("/QueryDeviceData", QueryDeviceData)
	http.HandleFunc("/SwitcherController", SwitcherController)
	http.HandleFunc("/", RedirectResponse)
}
func StartHttpService() {
	go func() {
		err := http.ListenAndServe(conf.Conf.LocalService, nil)
		if err != nil {
			zlog.Error("Http server", zlog.String("Err", err.Error()))
		}
	}()
}

//http://ip:port/QueryDeviceData?device=SJ2D1EQ8A31E9O42
func QueryDeviceData(rw http.ResponseWriter, req *http.Request) {
	m := req.URL.Query()
	device := m.Get("device")
	if device == "" {
		zlog.Info("QueryDeviceData param", zlog.String(req.RemoteAddr, "device is empty"))
		return
	}
	reqire, _ := json.Marshal(model.DeviceDataReq{
		UserId:   remoteuser.UserMgr.UserId,
		DeviceNo: device,
		CurrPage: 1,
		PageSize: 10,
	})
	request, err := http.NewRequest("POST", conf.Conf.RemoteService+model.DeviceData,
		strings.NewReader(string(reqire)))
	if err != nil {
		zlog.Error("QueryDeviceData NewRequest", zlog.String("Err", err.Error()))
		return
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("tlinkAppId", remoteuser.UserMgr.ClientId)
	request.Header.Set("Authorization", "Bearer "+remoteuser.UserMgr.AccessToken)
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		zlog.Error("QueryDeviceData Do", zlog.String("Err", err.Error()))
		return
	}
	zlog.Info("QueryDeviceData status", zlog.Int("Code", respone.StatusCode))
	defer respone.Body.Close()
	io.Copy(rw, respone.Body)
}

//http://ip:port/QueryDeviceData?device=SJ2D1EQ8A31E9O42&sensor=1&switcher=0
func SwitcherController(rw http.ResponseWriter, req *http.Request) {
	m := req.URL.Query()
	device := m.Get("device")
	if device == "" {
		zlog.Info("SwitcherController param", zlog.String(req.RemoteAddr, "device is empty"))
		return
	}
	sensorStr := m.Get("sensor")
	if sensorStr == "" {
		zlog.Info("SwitcherController param", zlog.String(req.RemoteAddr, "sensor is empty"))
		return
	}
	switcherStr := m.Get("switcher")
	if switcherStr == "" {
		zlog.Info("SwitcherController param", zlog.String(req.RemoteAddr, "switcher is empty"))
		return
	}
	sensor, _ := strconv.Atoi(sensorStr)
	switcher, _ := strconv.Atoi(switcherStr)
	reqire, _ := json.Marshal(model.SwitcherControllerReq{
		UserId:   remoteuser.UserMgr.UserId,
		DeviceNo: device,
		SensorId: int64(sensor),
		Switcher: int64(switcher),
	})
	request, err := http.NewRequest("POST", conf.Conf.RemoteService+model.SwitcherController,
		strings.NewReader(string(reqire)))
	if err != nil {
		zlog.Error("SwitcherController NewRequest", zlog.String("Err", err.Error()))
		return
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("tlinkAppId", remoteuser.UserMgr.ClientId)
	request.Header.Set("Authorization", "Bearer "+remoteuser.UserMgr.AccessToken)
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		zlog.Error("SwitcherController Do", zlog.String("Err", err.Error()))
		return
	}
	zlog.Info("SwitcherController status", zlog.Int("Code", respone.StatusCode))
	defer respone.Body.Close()
	io.Copy(rw, respone.Body)
}
