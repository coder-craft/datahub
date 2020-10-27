package web

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"../common/zlog"
	"../conf"
	"../model"
	"../remoteuser"
)

//http://ip:port/QueryDeviceData?device=SJ2D1EQ8A31E9O42
func init() {
	http.HandleFunc("/QueryDeviceData", QueryDeviceData)
}
func StartHttpService() {
	go func() {
		err := http.ListenAndServe(conf.Conf.LocalService, nil)
		if err != nil {
			zlog.Error("Http server", zlog.String("Err", err.Error()))
		}
	}()
}

func QueryDeviceData(rw http.ResponseWriter, req *http.Request) {
	m := req.URL.Query()
	device := m.Get("device")
	if device == "" {
		zlog.Info("WorldSrvApi param", zlog.String(req.RemoteAddr, "device is empty"))
		return
	}
	reqire, _ := json.Marshal(model.DeviceDataReq{
		UserId:   remoteuser.RemoteUserData.UserId,
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
	request.Header.Set("tlinkAppId", remoteuser.RemoteUserData.ClientId)
	request.Header.Set("Authorization", "Bearer "+remoteuser.RemoteUserData.AccessToken)
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		zlog.Error("QueryDeviceData Do", zlog.String("Err", err.Error()))
		return
	}
	zlog.Info("QueryDeviceData status",zlog.Int("Code",respone.StatusCode))
	defer respone.Body.Close()
	io.Copy(rw, respone.Body)
}
