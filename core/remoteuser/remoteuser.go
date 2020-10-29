package remoteuser

import (
	"../common/zlog"
	"../conf"
	"../model"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	Userstate_Init   = 0
	Userstate_Login  = 1
	Userstate_Access = 2
	Userstate_Flush  = 3
)

type RemoteUser struct {
	UserName     string
	UserId       int64
	ApiKey       string
	Secret       string
	AccessToken  string
	TokenType    string
	RefreshToken string
	ClientId     string
	ClientSecret string
	ExpiresIn    int64
}

func InitUser() error {
	err := UserLogin()
	if err != nil {
		return err
	}
	err = AccessToken()
	if err != nil {
		return err
	}
	return nil
}
func UserLogin() error {
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, err := json.Marshal(&model.UserLoginReq{
		UserName:  conf.Conf.UserName,
		PassWorld: conf.Conf.UserPass,
		ApiKey:    conf.Conf.UserApiKey,
	})
	if err != nil {
		return err
	}
	resp, err := client.Post(conf.Conf.RemoteService+model.UserLoginUrl, "application/json",
		bytes.NewBuffer(jsonStr))
	if err != nil {
		zlog.Error("Require err", zlog.String("Url:", "UserLogin"), zlog.String("Err", err.Error()))
		return err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	ulr := &model.UserLoginRespone{}
	err = json.Unmarshal(result, ulr)
	if err != nil {
		zlog.Error("Unmarshal userlogin respone", zlog.String("Err", err.Error()))
		return err
	}
	UserMgr.Secret = ulr.Secret
	UserMgr.UserId = ulr.Id
	UserMgr.ClientId = ulr.ClientId
	UserMgr.userstate = Userstate_Login
	zlog.Info("UserLogin success.", zlog.String("ClientId", UserMgr.ClientId))
	return nil
}
func AccessToken() error {
	url := fmt.Sprintf(model.AccessToken, conf.Conf.UserName, conf.Conf.UserPass)
	request, err := http.NewRequest("POST", conf.Conf.RemoteService+url, strings.NewReader(""))
	if err != nil {
		return err
	}
	request.Header.Set("Content-type", "text/plain")
	//request.Header.Set("tlinkAppId", RemoteUserData.ClientId)
	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(UserMgr.ClientId+
		":"+UserMgr.Secret)))
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		return err
	}
	defer respone.Body.Close()

	result, _ := ioutil.ReadAll(respone.Body)
	atr := &model.AccessTokenRespone{}
	err = json.Unmarshal(result, atr)
	if err != nil {
		zlog.Error("Unmarshal AccessToken respone", zlog.String("Err", err.Error()))
		return err
	}
	UserMgr.AccessToken = atr.Access_Token
	UserMgr.RefreshToken = atr.Refresh_Token
	UserMgr.TokenType = atr.Token_Type
	UserMgr.ClientId = atr.ClientId
	UserMgr.ClientSecret = atr.ClientSecret
	UserMgr.ExpiresIn = atr.Expires_In

	UserMgr.userstate = Userstate_Access
	zlog.Info("AccessToken success.", zlog.String("AccessToken", UserMgr.AccessToken),
		zlog.Int64("ExpiresIn", UserMgr.ExpiresIn))
	return nil
}
func FlushToken() error {
	url := fmt.Sprintf(model.RefreshToken, UserMgr.RefreshToken, UserMgr.ClientId, UserMgr.Secret)
	request, err := http.NewRequest("POST", conf.Conf.RemoteService+url, strings.NewReader(""))
	if err != nil {
		return err
	}
	request.Header.Set("Content-type", "text/plain")
	request.Header.Set("cache-control", "no-cache")
	client := &http.Client{}
	respone, err := client.Do(request)
	if err != nil {
		return err
	}
	defer respone.Body.Close()

	result, _ := ioutil.ReadAll(respone.Body)
	atr := &model.AccessTokenRespone{}
	err = json.Unmarshal(result, atr)
	if err != nil {
		zlog.Error("Unmarshal FlushToken respone", zlog.String("Err", err.Error()))
		return err
	}
	UserMgr.AccessToken = atr.Access_Token
	UserMgr.RefreshToken = atr.Refresh_Token
	UserMgr.TokenType = atr.Token_Type
	UserMgr.ClientId = atr.ClientId
	UserMgr.ClientSecret = atr.ClientSecret
	UserMgr.ExpiresIn = atr.Expires_In

	UserMgr.userstate = Userstate_Flush
	zlog.Info("FlushToken success.", zlog.String("AccessToken", UserMgr.AccessToken),
		zlog.Int64("ExpiresIn", UserMgr.ExpiresIn))
	return nil
}

var UserMgr = &UserManager{}

type UserManager struct {
	RemoteUser
	userstate int
	nextFlush int64
}

func (dm *UserManager) Name() string {
	return "UserManager"
}
func (dm *UserManager) Init() bool {
	return true
}
func (dm *UserManager) Update() bool {
	switch dm.userstate {
	case Userstate_Init:
		err := UserLogin()
		if err != nil {
			zlog.Error("UserLogin in update.", zlog.String("Err", err.Error()))
		}
	case Userstate_Login:
		err := AccessToken()
		if err != nil {
			zlog.Error("AccessToken in update.", zlog.String("Err", err.Error()))
		}
		dm.nextFlush = time.Now().Add(time.Minute).Unix()
	case Userstate_Access:
		fallthrough
	case Userstate_Flush:
		if dm.nextFlush < time.Now().Unix() {
			err := FlushToken()
			if err != nil {
				dm.userstate = Userstate_Login
			} else {
				dm.nextFlush = time.Now().Add(time.Minute).Unix()
			}
		}
	}
	return true
}
func (dm *UserManager) End() bool {
	return true
}
