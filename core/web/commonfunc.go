package web

import (
	"../common/zlog"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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

func RedirectResponse(rw http.ResponseWriter, req *http.Request) error {
	zlog.Info("RedirectResponse", zlog.String("Url", req.URL.String()))
	client := &http.Client{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(req.Method, req.URL.String(), strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	for k, v := range req.Header {
		request.Header.Set(k, v[0])
	}
	respone, err := client.Do(request)
	if err != nil {
		return err
	}
	defer respone.Body.Close()
	for k, v := range respone.Header {
		rw.Header().Set(k, v[0])
	}
	io.Copy(rw, respone.Body)
	return nil
}
