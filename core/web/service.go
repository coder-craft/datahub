package web

import (
	"../common/zlog"
	"../conf"
	"../mysql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

//http://ip:port/QueryDeviceData?device=171CB73F333A
func init() {
	http.HandleFunc("/QueryDeviceData", QueryDeviceData)
}
func StartHttpService() {
	address, err := net.Listen(conf.Conf.HttpService.Network, conf.Conf.HttpService.Address)
	if err != nil {
		zlog.Error("Start http service", zlog.String("Err", err.Error()))
		return
	}
	go func() {
		err := http.Serve(address, nil)
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
	device = base64.StdEncoding.EncodeToString([]byte(device))
	deviceData := mysql.FindDeviceEvent(device)
	var buff string
	for _, value := range deviceData {
		str := fmt.Sprintf("设备ID:%v,采集时间:%v,", value[0], value[1])
		var subData []int32
		err := json.Unmarshal([]byte(value[3]), &subData)
		if err != nil {
			zlog.Error("Unmarshal data", zlog.String("Err", err.Error()))
			continue
		}
		if len(subData) != 12 {
			zlog.Error("QueryDeviceData", zlog.Int("data length", len(subData)))
			continue
		}
		str += fmt.Sprintf("温度:%v°C,", float32(subData[0])/10)
		str += fmt.Sprintf("湿度:%vRH,", float32(subData[1])/10)
		str += fmt.Sprintf("噪声:%vdB,", float32(subData[2])/10)
		str += fmt.Sprintf("大气压:%vKpa,", float32(subData[3])/10)
		str += fmt.Sprintf("二氧化碳:%vppm,", subData[4])
		str += fmt.Sprintf("负离子:%vc㎡,", subData[5])
		str += fmt.Sprintf("风速:%vm/s,", float32(subData[6])/10)
		str += fmt.Sprintf("风向:%v,", subData[7])
		str += fmt.Sprintf("紫外线:%vW/m2,", subData[8])
		str += fmt.Sprintf("信号强度:%vRSS,", subData[9])
		str += fmt.Sprintf("错误码:%v,", subData[10])
		str += fmt.Sprintf("版本号:%Version,", subData[11])
		buff += str
	}
	_, err := rw.Write([]byte(buff))
	if err != nil {
		zlog.Error("ResponseWriter", zlog.String("write err", err.Error()))
	}
}

//，湿度52.9%RH，，，，，，，，，，
type DeviceData struct {
	device_id    string `json:"设备ID"` //
	event_time   string `json:"采集时间"` //
	event_data1  string `json:"温度"`   //温度27.1°C
	event_data2  string `json:"温度"`   //湿度52.9%RH
	event_data3  string `json:"温度"`   //噪声51.9dB
	event_data4  string `json:"温度"`   //大气压100.9Kpa
	event_data5  string `json:"温度"`   //二氧化碳1057ppm
	event_data6  string `json:"温度"`   //负离子1233c㎡
	event_data7  string `json:"温度"`   //风速3.7m/s
	event_data8  string `json:"温度"`   //风向 西南
	event_data9  string `json:"温度"`   //紫外线0W/m2
	event_data10 string `json:"温度"`   //信号强度139RSS
	event_data11 string `json:"温度"`   //错误码 202
	event_data12 string `json:"温度"`   //版本号 2Version
}
