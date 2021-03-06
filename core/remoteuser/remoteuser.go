package remoteuser

import (
	"../common/zlog"
	"../conf"
	"../model"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	Userstate_Init   = 0
	Userstate_Login  = 1
	Userstate_Access = 2
	Userstate_Flush  = 3
)

var RemoteUserMgr = &RemoteUserManager{}

type RemoteUserManager struct {
	User sync.Map
}

func (rum *RemoteUserManager) Name() string {
	return "RemoteUserManager"
}
func (rum *RemoteUserManager) Init() bool {
	return true
}
func (rum *RemoteUserManager) Update() bool {
	deleteKeys := []string{}
	rum.User.Range(func(key, value interface{}) bool {
		user := value.(*model.RemoteUser)
		if user == nil {
			return false
		}
		switch user.Userstate {
		case Userstate_Init:
			user := userLogin(user.UserName, user.PassWord)
			if user == nil {
				zlog.Error("UserLogin failed in update.")
			} else {
				rum.StoreUser(user)
			}
		case Userstate_Login:
			err := AccessToken(user)
			if err != nil {
				zlog.Error("AccessToken in update.", zlog.String("Err", err.Error()))
			}
		case Userstate_Access:
			fallthrough
		case Userstate_Flush:
			if user.NextFlush < time.Now().Unix() {
				deleteKeys = append(deleteKeys, user.UserName+user.PassWord)
				//err := FlushToken(user)
				//if err != nil {
				//	user.Userstate = Userstate_Init
				//}
			}
		}
		return true
	})
	for _, key := range deleteKeys {
		rum.DelUser(key)
	}
	return true
}
func (rum *RemoteUserManager) End() bool {
	return true
}
func (rum *RemoteUserManager) DelUser(key string) {
	rum.User.Delete(key)
}
func (rum *RemoteUserManager) GetUser(key string) *model.RemoteUser {
	u, ok := rum.User.Load(key)
	if ok {
		user := u.(*model.RemoteUser)
		if user == nil {
			return nil
		}
		if user.NextFlush < time.Now().Unix() {
			return nil
		} else {
			return user
		}
	} else {
		return nil
	}
}
func (rum *RemoteUserManager) StoreUser(user *model.RemoteUser) {
	zlog.Info("StoreUser", zlog.String("User:", user.UserName), zlog.String("Pass:", user.PassWord))
	key := user.UserName + user.PassWord
	rum.User.Store(key, user)
}
func (rum *RemoteUserManager) UserLogin(userName, passWord string) *model.RemoteUser {
	user := userLogin(userName, passWord)
	if user == nil {
		return nil
	}
	err := AccessToken(user)
	if err != nil {
		zlog.Error("AccessToken", zlog.String("Err", err.Error()), zlog.String("UserName", userName))
		return nil
	}
	rum.StoreUser(user)
	return user
}
func userLogin(userName, passWord string) *model.RemoteUser {
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, err := json.Marshal(&model.UserLoginReq{
		UserName:  userName,
		PassWorld: passWord,
		ApiKey:    conf.Conf.UserApiKey,
	})
	if err != nil {
		return nil
	}
	resp, err := client.Post(conf.Conf.RemoteService+model.UserLoginUrl, "application/json",
		bytes.NewBuffer(jsonStr))
	if err != nil {
		zlog.Error("Require err", zlog.String("Url:", "UserLogin"), zlog.String("Err", err.Error()))
		return nil
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	ulr := &model.UserLoginRespone{}
	err = json.Unmarshal(result, ulr)
	if err != nil {
		zlog.Error("Unmarshal userlogin respone", zlog.String("Err", err.Error()))
		return nil
	}
	user := &model.RemoteUser{UserName: userName, PassWord: passWord}
	user.Secret = ulr.Secret
	user.UserId = ulr.Id
	user.ClientId = ulr.ClientId
	user.Userstate = Userstate_Login
	zlog.Info("UserLogin success.", zlog.String("ClientId", user.ClientId))
	return user
}
func AccessToken(user *model.RemoteUser) error {
	url := fmt.Sprintf(model.AccessToken, user.UserName, user.PassWord)
	request, err := http.NewRequest("POST", conf.Conf.RemoteService+url, strings.NewReader(""))
	if err != nil {
		return err
	}
	request.Header.Set("Content-type", "text/plain")
	//request.Header.Set("tlinkAppId", RemoteUserData.ClientId)
	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user.ClientId+
		":"+user.Secret)))
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
	if len(atr.Access_Token) == 0 || atr.Expires_In == 0 {
		return errors.New("AccessToken error.")
	}
	user.AccessToken = atr.Access_Token
	user.RefreshToken = atr.Refresh_Token
	user.TokenType = atr.Token_Type
	user.ClientId = atr.ClientId
	user.ClientSecret = atr.ClientSecret
	user.ExpiresIn = atr.Expires_In

	user.Userstate = Userstate_Access
	user.NextFlush = time.Now().Add(time.Second * time.Duration(user.ExpiresIn)).Unix()
	//user.NextFlush = time.Now().Add(time.Second * 60).Unix()
	zlog.Info("AccessToken success.", zlog.String("AccessToken", user.AccessToken),
		zlog.Int64("ExpiresIn", user.ExpiresIn))
	return nil
}
func FlushToken(user *model.RemoteUser) error {
	url := fmt.Sprintf(model.RefreshToken, user.RefreshToken, user.ClientId, user.Secret)
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
	if len(atr.Access_Token) == 0 || atr.Expires_In == 0 {
		return errors.New("FlushToken error.")
	}
	user.AccessToken = atr.Access_Token
	user.RefreshToken = atr.Refresh_Token
	user.TokenType = atr.Token_Type
	user.ClientId = atr.ClientId
	user.ClientSecret = atr.ClientSecret
	user.ExpiresIn = atr.Expires_In

	user.Userstate = Userstate_Flush
	if user.ExpiresIn == 0 {
		user.ExpiresIn = 3000
	}
	user.NextFlush = time.Now().Add(time.Second * time.Duration(user.ExpiresIn/2)).Unix()
	zlog.Info("FlushToken success.", zlog.String("AccessToken", user.AccessToken),
		zlog.Int64("ExpiresIn", user.ExpiresIn))
	return nil
}
