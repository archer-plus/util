package openplatform

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/archer-plus/util/config"

	"github.com/archer-plus/util"
	"github.com/archer-plus/util/logx"
)

const (
	qrConnectURL          = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	userInfoURL           = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=%s"
	accessTokenURL        = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	refreshAccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	checkAccessTokenURL   = "https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s"
)

var instance *Oauth
var conf = Config{}
var once sync.Once

type Config struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
}

type Oauth struct {
	AppID     string
	AppSecret string
}

type CommonError struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type AccessToken struct {
	CommonError
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
}

//UserInfo 用户授权获取到用户信息
type UserInfo struct {
	CommonError

	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int32    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

func Init() {
	err := config.UnmarshalKey("wechat", &conf)
	if err != nil {
		logx.Sugar.Warnf("未发现wechat配置信息, %v", err)
		return
	}
	once.Do(func() {
		instance = &Oauth{AppID: conf.AppID, AppSecret: conf.AppSecret}
	})
}

func NewOauth() *Oauth {
	return instance
}

// GetQRConnectURL 获取登录二维码地址
func (c *Oauth) GetQRConnectURL(redirectURL, scope, state string) string {
	urlstr := url.QueryEscape(redirectURL)
	return fmt.Sprintf(qrConnectURL, c.AppID, urlstr, scope, state)
}

// GetUserAccessToken 根据code获取acess token
func (c *Oauth) GetUserAccessToken(code string) (result AccessToken, err error) {
	urlstr := fmt.Sprintf(accessTokenURL, c.AppID, c.AppSecret, code)
	var response []byte
	response, err = util.HTTPGet(urlstr)
	if err != nil {
		err = fmt.Errorf("GetUserAccessToken error: %v", err.Error())
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("GetUserAccessToken json error: %v", err.Error())
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

// RefreshAccessToken 刷新access_token
func (c *Oauth) RefreshAccessToken(refreshToken string) (result AccessToken, err error) {
	urlstr := fmt.Sprintf(refreshAccessTokenURL, c.AppID, refreshToken)
	var response []byte
	response, err = util.HTTPGet(urlstr)
	if err != nil {
		err = fmt.Errorf("RefreshAccessToken error: %v", err.Error())
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("RefreshAccessToken json error: %v", err.Error())
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("RefreshAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

// CheckAccessToken 检查access token是否有效
func (c *Oauth) CheckAccessToken(accessToken, openID string) (b bool, err error) {
	urlstr := fmt.Sprintf(checkAccessTokenURL, accessToken, c.AppID)
	var response []byte
	response, err = util.HTTPGet(urlstr)
	if err != nil {
		err = fmt.Errorf("CheckAccessToken error: %v", err.Error())
		return
	}
	var result CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		b = false
		err = fmt.Errorf("CheckAccessToken json error: %v", err.Error())
		return
	}
	if result.ErrCode != 0 {
		b = false
		err = fmt.Errorf("CheckAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	b = true
	return
}

// GetUserInfo 获取用户信息
func (c *Oauth) GetUserInfo(accessToken, openID, lang string) (result UserInfo, err error) {
	if lang == "" {
		lang = "zh_CN"
	}
	urlstr := fmt.Sprintf(userInfoURL, accessToken, openID, lang)
	var response []byte
	response, err = util.HTTPGet(urlstr)
	if err != nil {
		err = fmt.Errorf("GetUserInfo error: %v", err.Error())
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		err = fmt.Errorf("GetUserInfo json error: %v", err.Error())
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserInfo error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}
