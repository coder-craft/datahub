package model

const UserLoginUrl = "/api/user/userLogin"

type UserLoginReq struct {
	UserName  string `json:"userName"`
	PassWorld string `json:"password"`
	ApiKey    string `json:"apiKey"`
}
type UserLoginRespone struct {
	Msg      string
	Flag     string
	ClientId string
	Mobile   string
	Id       int64
	Secret   string
	Avatar   string
	UserName string
	Email    string
}

const AccessToken = "/oauth/token?grant_type=password&username=%v&password=%v"

type AccessTokenRespone struct {
	Access_Token  string
	Token_Type    string
	Refresh_Token string
	Expires_In    int64
	Scope         string
	ClientId      string
	ClientSecret  string
}

const RefreshToken = "/oauth/token?grant_type=refresh_token&refresh_token=%v"

const DeviceData = "/api/device/getSingleDeviceDatas"

type DeviceDataReq struct {
	UserId   int64  `json:"userId"`
	DeviceId string `json:"deviceId"`
	DeviceNo string `json:"deviceNo"`
	CurrPage int64  `json:"currPage"`
	PageSize int64  `json:"pageSize"`
}

const SwitcherController = "/api/device/switcherController"
