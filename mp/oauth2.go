package mp

import (
	"fmt"
	"net/url"
	"errors"
	"encoding/json"
	"github.com/qjw/go-wx-sdk/utils"
)

const (
	oauth2_authorize    = "https://open.weixin.qq.com/connect/oauth2/authorize"
	oauth2_access_token = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	oauth2_sns_userinfo = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
)

func (this WechatApi) authorizeUrl(redirectURI, scope, state string) string {
	return oauth2_authorize + "?appid=" + url.QueryEscape(this.Context.Config.AppID) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code&scope=" + url.QueryEscape(scope) +
		"&state=" + url.QueryEscape(state) +
		"#wechat_redirect"
}

func (this WechatApi) AuthorizeUserinfo(redirectURI, state string) string {
	return this.authorizeUrl(redirectURI, "snsapi_userinfo", state)
}

func (this WechatApi) AuthorizeBase(redirectURI, state string) string {
	return this.authorizeUrl(redirectURI, "snsapi_base", state)
}

type SnsUserinfo struct {
	utils.CommonError
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Unionid    string   `json:"unionid"`
	Privilege  []string `json:"privilege"`
}

func (this WechatApi) GetAuthorizeSnsUserinfo(access_token, openid string) (*SnsUserinfo, error) {
	var res SnsUserinfo
	uri := fmt.Sprintf(oauth2_sns_userinfo, access_token, openid)

	response, _, err := utils.HTTPGet(uri)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &res)
	if err != nil {
		return nil, errors.New(string(response))
	}
	return &res, nil
}

type AuthorizeAccessToken struct {
	utils.CommonError
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}

func (this WechatApi) GetAuthorizeAccessToken(code string) (*AuthorizeAccessToken, error) {
	var res AuthorizeAccessToken
	uri := fmt.Sprintf(oauth2_access_token,
		this.Context.Config.AppID,
		this.Context.Config.AppSecret,
		code,
	)

	response, _, err := utils.HTTPGet(uri)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(response, &res)
	if err != nil {
		return nil, errors.New(string(response))
	}
	return &res, nil
}
