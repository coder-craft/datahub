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

var RemoteUserData RemoteUser

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
	if err !=nil {
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
	RemoteUserData.Secret = ulr.Secret
	RemoteUserData.UserId = ulr.Id
	RemoteUserData.ClientId = ulr.ClientId
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
	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(RemoteUserData.ClientId+
		":"+RemoteUserData.Secret)))
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
		zlog.Error("Unmarshal userlogin respone", zlog.String("Err", err.Error()))
		return err
	}
	RemoteUserData.AccessToken = atr.Access_Token
	RemoteUserData.RefreshToken = atr.Refresh_Token
	RemoteUserData.TokenType = atr.Token_Type
	RemoteUserData.ClientId = atr.ClientId
	RemoteUserData.ClientSecret = atr.ClientSecret
	RemoteUserData.ExpiresIn = atr.Expires_In
	return nil
}
func FlushToken() error {
	return nil
}
