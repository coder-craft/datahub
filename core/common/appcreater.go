package common

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

func Md5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return string(h.Sum(nil))
}
func CreateApp() (string, string, string, int64) {
	appid := Md5([]byte(time.Now().String()))
	appidLen := len(appid)
	seedIndex := rand.Intn(appidLen - 5)
	seed := appid[seedIndex : seedIndex+5]
	ts := time.Now().UnixNano()
	appkey := Md5([]byte(fmt.Sprintf("%v,%v,%v", appid, seed, ts)))
	return appid, appkey, seed, ts
}
func CheckAppKey(appid, appkey, seed string, ts int64) bool {
	key := Md5([]byte(fmt.Sprintf("%v,%v,%v", appid, seed, ts)))
	return key == appkey
}
