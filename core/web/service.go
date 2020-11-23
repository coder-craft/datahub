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
	http.HandleFunc("/QueryDevice", QueryDevice)
	http.HandleFunc("/QueryDeviceDatas", GetDeviceSensorDatas)
	http.HandleFunc("/SwitcherController", SwitcherController)
	http.HandleFunc("/GetSingleSensorDatas", GetSingleSensorDatas)
	//http.HandleFunc("/", RedirectResponse)
}
func StartHttpService() {
	go func() {
		err := http.ListenAndServe(conf.Conf.LocalService, nil)
		if err != nil {
			zlog.Error("Http server", zlog.String("Err", err.Error()))
		}
	}()
}

//http://ip:port/QueryDeviceData?device=SJ2D1EQ8A31E9O42&username=test&password=1111&curpage=1&pagesize=10
func QueryDeviceData(rw http.ResponseWriter, req *http.Request) {
	m := req.URL.Query()
	device := m.Get("device")
	if device == "" {
		zlog.Info("QueryDeviceData param", zlog.String(req.RemoteAddr, "device is empty"))
		return
	}
	userName := m.Get("username")
	passWord := m.Get("password")
	user := remoteuser.RemoteUserMgr.GetUser(userName + passWord)
	if user == nil {
		user = remoteuser.RemoteUserMgr.UserLogin(userName, passWord)
	}
	if user == nil {
		zlog.Info("QueryDeviceData param", zlog.String(req.RemoteAddr, "user error"))
		return
	}
	curpage, _ := strconv.Atoi(m.Get("curpage"))
	pagesize, _ := strconv.Atoi(m.Get("pagesize"))
	if curpage == 0 {
		curpage = 1
	}
	if pagesize == 0 {
		pagesize = 10
	}
	reqire, _ := json.Marshal(model.DeviceDataReq{
		UserId:   user.UserId,
		DeviceNo: device,
		CurrPage: int64(curpage),
		PageSize: int64(pagesize),
	})
	request, err := http.NewRequest("POST", conf.Conf.RemoteService+model.DeviceData,
		strings.NewReader(string(reqire)))
	if err != nil {
		zlog.Error("QueryDeviceData NewRequest", zlog.String("Err", err.Error()))
		return
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("tlinkAppId", user.ClientId)
	request.Header.Set("Authorization", "Bearer "+user.AccessToken)
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		zlog.Error("QueryDeviceData Do", zlog.String("Err", err.Error()))
		return
	}
	zlog.Info("QueryDeviceData status", zlog.Int("Code", respone.StatusCode))
	if respone.StatusCode != 200 {

		remoteuser.RemoteUserMgr.DelUser(userName + passWord)
	}
	defer respone.Body.Close()
	io.Copy(rw, respone.Body)
}

//http://ip:port/SwitcherController?device=SJ2D1EQ8A31E9O42&sensor=1&switcher=0&username=test&password=1111
func SwitcherController(rw http.ResponseWriter, req *http.Request) {
	m := req.URL.Query()
	device := m.Get("device")
	if device == "" {
		zlog.Info("SwitcherController param", zlog.String(req.RemoteAddr, "device is empty"))
		return
	}
	userName := m.Get("username")
	passWord := m.Get("password")
	user := remoteuser.RemoteUserMgr.GetUser(userName + passWord)
	if user == nil {
		user = remoteuser.RemoteUserMgr.UserLogin(userName, passWord)
	}
	if user == nil {
		zlog.Info("SwitcherController param", zlog.String(req.RemoteAddr, "user error"))
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
		UserId:   user.UserId,
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
	request.Header.Set("tlinkAppId", user.ClientId)
	request.Header.Set("Authorization", "Bearer "+user.AccessToken)
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		zlog.Error("SwitcherController Do", zlog.String("Err", err.Error()))
		return
	}
	zlog.Info("SwitcherController status", zlog.Int("Code", respone.StatusCode))
	if respone.StatusCode != 200 {
		remoteuser.RemoteUserMgr.DelUser(userName + passWord)
	}
	defer respone.Body.Close()
	io.Copy(rw, respone.Body)
}

//http://ip:port/QueryDevice?username=test&password=1111&curpage=1&pagesize=10
func QueryDevice(rw http.ResponseWriter, req *http.Request) {
	m := req.URL.Query()
	userName := m.Get("username")
	passWord := m.Get("password")
	if len(userName) == 0 || len(passWord) == 0 {
		zlog.Info("QueryDevice param", zlog.String("Err", "Param error"))
		return
	}
	user := remoteuser.RemoteUserMgr.GetUser(userName + passWord)
	if user == nil {
		user = remoteuser.RemoteUserMgr.UserLogin(userName, passWord)
	}
	if user == nil {
		zlog.Info("QueryDevice param", zlog.String(req.RemoteAddr, "user error"))
		return
	}
	curpage, _ := strconv.Atoi(m.Get("curpage"))
	pagesize, _ := strconv.Atoi(m.Get("pagesize"))
	if curpage == 0 {
		curpage = 1
	}
	if pagesize == 0 {
		pagesize = 10
	}
	reqire, _ := json.Marshal(model.GetDevicesReq{
		UserId:   user.UserId,
		CurrPage: int64(curpage),
		PageSize: int64(pagesize),
	})
	request, err := http.NewRequest("POST", conf.Conf.RemoteService+model.GetDevices,
		strings.NewReader(string(reqire)))
	if err != nil {
		zlog.Error("QueryDevice NewRequest", zlog.String("Err", err.Error()))
		return
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("tlinkAppId", user.ClientId)
	request.Header.Set("Authorization", "Bearer "+user.AccessToken)
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		zlog.Error("QueryDevice Do", zlog.String("Err", err.Error()))
		return
	}
	zlog.Info("QueryDevice status", zlog.Int("Code", respone.StatusCode))
	if respone.StatusCode != 200 {
		remoteuser.RemoteUserMgr.DelUser(userName + passWord)
	}
	defer respone.Body.Close()
	io.Copy(rw, respone.Body)
}

//http://ip:port/GetDeviceSensorDatas?username=test&password=1111&curpage=1&pagesize=10
func GetDeviceSensorDatas(rw http.ResponseWriter, req *http.Request) {
	m := req.URL.Query()
	userName := m.Get("username")
	passWord := m.Get("password")
	if len(userName) == 0 || len(passWord) == 0 {
		zlog.Info("GetDeviceSensorDatas param", zlog.String("Err", "Param error"))
		return
	}
	user := remoteuser.RemoteUserMgr.GetUser(userName + passWord)
	if user == nil {
		user = remoteuser.RemoteUserMgr.UserLogin(userName, passWord)
	}
	if user == nil {
		zlog.Info("GetDeviceSensorDatas param", zlog.String(req.RemoteAddr, "user error"))
		return
	}
	curpage, _ := strconv.Atoi(m.Get("curpage"))
	pagesize, _ := strconv.Atoi(m.Get("pagesize"))
	if curpage == 0 {
		curpage = 1
	}
	if pagesize == 0 {
		pagesize = 10
	}
	reqire, _ := json.Marshal(model.GetDeviceSensorDatasReq{
		UserId:   user.UserId,
		CurrPage: int64(curpage),
		PageSize: int64(pagesize),
	})
	request, err := http.NewRequest("POST", conf.Conf.RemoteService+model.GetDeviceSensorDatas,
		strings.NewReader(string(reqire)))
	if err != nil {
		zlog.Error("GetDeviceSensorDatas NewRequest", zlog.String("Err", err.Error()))
		return
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("tlinkAppId", user.ClientId)
	request.Header.Set("Authorization", "Bearer "+user.AccessToken)
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		zlog.Error("GetDeviceSensorDatas Do", zlog.String("Err", err.Error()))
		return
	}
	zlog.Info("GetDeviceSensorDatas status", zlog.Int("Code", respone.StatusCode))
	if respone.StatusCode != 200 {
		remoteuser.RemoteUserMgr.DelUser(userName + passWord)
	}
	defer respone.Body.Close()
	io.Copy(rw, respone.Body)
}

//http://ip:port/GetSingleSensorDatas?username=test&password=1111&sensor=1&curpage=1&pagesize=10
func GetSingleSensorDatas(rw http.ResponseWriter, req *http.Request) {
	m := req.URL.Query()
	userName := m.Get("username")
	passWord := m.Get("password")
	if len(userName) == 0 || len(passWord) == 0 {
		zlog.Info("GetSingleSensorDatas param", zlog.String("Err", "Param error"))
		return
	}
	user := remoteuser.RemoteUserMgr.GetUser(userName + passWord)
	if user == nil {
		user = remoteuser.RemoteUserMgr.UserLogin(userName, passWord)
	}
	if user == nil {
		zlog.Info("GetSingleSensorDatas param", zlog.String(req.RemoteAddr, "user error"))
		return
	}
	curpage, _ := strconv.Atoi(m.Get("curpage"))
	pagesize, _ := strconv.Atoi(m.Get("pagesize"))
	if curpage == 0 {
		curpage = 1
	}
	if pagesize == 0 {
		pagesize = 10
	}
	sensorStr := m.Get("sensor")
	if sensorStr == "" {
		zlog.Info("GetSingleSensorDatas param", zlog.String(req.RemoteAddr, "sensor is empty"))
		return
	}
	sensorId, _ := strconv.Atoi(sensorStr)
	reqire, _ := json.Marshal(model.SingleSensorDatasReq{
		UserId:   user.UserId,
		SensorId: int64(sensorId),
	})
	request, err := http.NewRequest("POST", conf.Conf.RemoteService+model.GetSingleSensorDatas,
		strings.NewReader(string(reqire)))
	if err != nil {
		zlog.Error("GetSingleSensorDatas NewRequest", zlog.String("Err", err.Error()))
		return
	}
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("tlinkAppId", user.ClientId)
	request.Header.Set("Authorization", "Bearer "+user.AccessToken)
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		zlog.Error("GetSingleSensorDatas Do", zlog.String("Err", err.Error()))
		return
	}
	zlog.Info("GetSingleSensorDatas status", zlog.Int("Code", respone.StatusCode))
	if respone.StatusCode != 200 {
		remoteuser.RemoteUserMgr.DelUser(userName + passWord)
	}
	defer respone.Body.Close()
	io.Copy(rw, respone.Body)
}
