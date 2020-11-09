package model

type RemoteUser struct {
	UserName     string
	PassWord     string
	UserId       int64
	ApiKey       string
	Secret       string
	AccessToken  string
	TokenType    string
	RefreshToken string
	ClientId     string
	ClientSecret string
	ExpiresIn    int64
	Userstate    int
	NextFlush    int64
}

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

const RefreshToken = "/oauth/token?grant_type=refresh_token&refresh_token=%v&client_id=%v&client_secret=%v"

const DeviceData = "/api/device/getSingleDeviceDatas"

type DeviceDataReq struct {
	UserId   int64  `json:"userId"`
	DeviceId string `json:"deviceId"`
	DeviceNo string `json:"deviceNo"`
	CurrPage int64  `json:"currPage"`
	PageSize int64  `json:"pageSize"`
}

const SwitcherController = "/api/device/switcherController"

type SwitcherControllerReq struct {
	UserId   int64  `json:"userId"`
	DeviceNo string `json:"deviceNo"`
	Switcher int64  `json:"switcher"`
	SensorId int64  `json:"sensorId"`
}

const GetDevices = "api/device/getDevices"

type GetDevicesReq struct {
	UserId   int64 `json:"userId"`
	CurrPage int64 `json:"currPage"`
	PageSize int64 `json:"pageSize"`
}

const GetDeviceSensorDatas = "/api/device/getDeviceSensorDatas"

type GetDeviceSensorDatasReq struct {
	UserId   int64 `json:"userId"`
	CurrPage int64 `json:"currPage"`
	PageSize int64 `json:"pageSize"`
}

const GetSingleSensorDatas = "/api/device/getSingleSensorDatas"

type SingleSensorDatasReq struct {
	UserId   int64 `json:"userId"`
	SensorId int64  `json:"sensorId"`
}