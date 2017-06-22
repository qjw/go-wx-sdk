package corp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"net/url"
	"strconv"
)

const (
	oauth2_authorize     = "https://open.weixin.qq.com/connect/oauth2/authorize"
	oauth2_getuserinfo   = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s"
	oauth2_getuserdetail = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserdetail?access_token=%s"
)

func (this CorpApi) authorizeUrl(redirectURI, scope, state string, agentID *int64) string {
	url := oauth2_authorize + "?appid=" + url.QueryEscape(this.Context.Config.CorpID) +
		"&redirect_uri=" + url.QueryEscape(redirectURI) +
		"&response_type=code&scope=" + url.QueryEscape(scope) +
		"&state=" + url.QueryEscape(state)
	if agentID != nil {
		url = url + "&agentid=" + strconv.FormatInt(*agentID, 10)
	}
	return url + "#wechat_redirect"
}

func (this CorpApi) AuthorizeUserinfo(redirectURI, state string, agentID int64) string {
	return this.authorizeUrl(redirectURI, "snsapi_userinfo", state, &agentID)
}

func (this CorpApi) AuthorizeBase(redirectURI, state string) string {
	return this.authorizeUrl(redirectURI, "snsapi_base", state, nil)
}

func (this CorpApi) AuthorizePrivateInfo(redirectURI, state string, agentID int64) string {
	return this.authorizeUrl(redirectURI, "snsapi_privateinfo", state, &agentID)
}

type OauthUserDetail struct {
	utils.CommonError
	UserID     string  `json:"userid,omitempty"`
	Name       string  `json:"name,omitempty"`
	Gender     string  `json:"gender"`
	Department []int64 `json:"department,omitempty"`
	Avatar     string  `json:"avatar,omitempty"`
	Email      string  `json:"email,omitempty"`
	Position   string  `json:"position,omitempty"`
	Mobile     string  `json:"mobile,omitempty"`
}

func (this CorpApi) Oauth2GetUserDetail(user_ticket string) (*OauthUserDetail, error) {
	var res OauthUserDetail
	if err := this.DoPostObject(oauth2_getuserdetail,
		&struct{UserTicket string `json:"user_ticket"`}{UserTicket:user_ticket}, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type AuthorizeAccessToken struct {
	utils.CommonError
	ExpiresIn  int64  `json:"expires_in,omitempty"`
	UserTicket string `json:"user_ticket,omitempty"`
	DeviceId   string `json:"DeviceId,omitempty"`
	UserId     string `json:"UserId,omitempty"`
}

func (this CorpApi) Oauth2GetUserInfo(code string) (*AuthorizeAccessToken, error) {
	var res AuthorizeAccessToken
	err := this.DoGet(oauth2_getuserinfo, &res, code)
	if err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}
