package web

import (
	"../common/zlog"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var ApiDefaultTimeout = time.Second * 30

type HandlerWrapper func(*WebApiEvent, []byte) bool

type WebApiEvent struct {
	req      *http.Request
	path     string
	rawQuery string
	body     []byte
	h        HandlerWrapper
	res      chan map[string]interface{}
}

func (this *WebApiEvent) Done() error {
	this.h(this, this.body)
	return nil
}

func (this *WebApiEvent) Response(data map[string]interface{}) {
	this.res <- data
}

func webApiResponse(rw http.ResponseWriter, params map[string]interface{}) bool {
	data, err := json.Marshal(params)
	if err != nil {
		zlog.Error("webApiResponse Marshal", zlog.String("Err", err.Error()))
		return false
	}

	dataLen := len(data)
	rw.Header().Set("Content-Length", fmt.Sprintf("%v", dataLen))
	rw.WriteHeader(http.StatusOK)
	pos := 0
	for pos < dataLen {
		writeLen, err := rw.Write(data[pos:])
		if err != nil {
			zlog.Error("WebApiResponse SendData", zlog.String("Err", err.Error()),
				zlog.String("data", string(data[:])), zlog.Int("pos", pos),
				zlog.Int("writelen", writeLen), zlog.Int("dataLen", dataLen))
			return false
		}
		pos += writeLen
	}

	return true
}
