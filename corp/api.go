package corp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"strconv"
)

func (this CorpApi) GetAccessToken() (*string, error) {
	accessToken, err := this.Context.GetAccessToken()
	if err == nil {
		return &accessToken, nil
	} else {
		return nil, err
	}
}

type SignJsRes struct {
	AppID     string `json:"appid"`
	Timestamp string `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

func (this CorpApi) SignJsTicket(nonceStr, timestamp, url string) (*SignJsRes, error) {
	jsTicket, err := this.Context.GetJsTicket()
	if err != nil {
		return nil, err
	}
	if timestamp == "" {
		bb := utils.GetCurrTs()
		timestamp = strconv.FormatInt(bb, 10)
	}
	sign := utils.WXConfigSign(jsTicket, nonceStr, timestamp, url)
	return &SignJsRes{
		Signature: sign,
		AppID:     this.Context.Config.CorpID,
		NonceStr:  nonceStr,
		Timestamp: timestamp,
	}, nil
}

const (
	getcallbackip = "https://qyapi.weixin.qq.com/cgi-bin/getcallbackip?access_token=%s"
)

type IpList struct {
	IpList []string `json:"ip_list"`
}

func (this CorpApi) GetIpList() (*IpList, error) {
	var res IpList
	if err := this.DoGet(getcallbackip, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) GetJsTicket() (*string, error) {
	jsTicket, err := this.Context.GetJsTicket()
	if err == nil {
		return &jsTicket, nil
	} else {
		return nil, err
	}
}

const (
	convert2Openid = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_openid?access_token=%s"
	Convert2Userid = "https://qyapi.weixin.qq.com/cgi-bin/user/convert_to_userid?access_token=%s"
)

type Convert2OpenIDObj struct {
	UserID  string `json:"userid"`
	AgentID int64  `json:"agentid"`
}

type Convert2OpenIDRes struct {
	utils.CommonError
	OpenID string `json:"openid"`
	AppID  string `json:"appid"`
}

func (this CorpApi) Convert2OpenID(param *Convert2OpenIDObj) (*Convert2OpenIDRes, error) {
	var res Convert2OpenIDRes
	if err := this.DoPostObject(convert2Openid, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}


type Convert2UserIDRes struct {
	utils.CommonError
	UserID string `json:"userid"`
}

func (this CorpApi) Convert2UserID(openid string) (*Convert2UserIDRes, error) {
	var res Convert2UserIDRes
	if err := this.DoPostObject(Convert2Userid, &struct{
		OpenID string `json:"openid"`
	}{OpenID:openid}, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}