package conf

import (
	"../common/zlog"
	"encoding/json"
	"flag"
	"io/ioutil"
)

var (
	confPath string
	logPath  string
	Conf     *Config
)

type Config struct {
	UserName      string
	UserPass      string
	UserApiKey    string
	LocalService  string
	RemoteService string
	MysqlAccounts string
	WhiteList     []string
	BlackList     []string
}

func init() {
	flag.StringVar(&confPath, "conf", "datahub.cfg", "default config path")
	flag.StringVar(&logPath, "log", "./msg_server.log", "default log path.")
}
func Init() error {
	zlog.InitDefaultZapLog(logPath)
	Conf = Default()
	buff, err := ioutil.ReadFile(confPath)
	if err != nil {
		zlog.Error("Read config file error", zlog.String("Err", err.Error()))
	} else {
		err = json.Unmarshal(buff, Conf)
		if err != nil {
			zlog.Error("Unmarshal config file error", zlog.String("Err", err.Error()))
			return err
		}
	}
	return nil
}
func Default() *Config {
	return &Config{
		LocalService:  "127.0.0.1:3500",
		RemoteService: "127.0.0.1:3501",
	}
}
