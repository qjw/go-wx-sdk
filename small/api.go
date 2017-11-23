package small

import (
	"fmt"
	"github.com/qjw/go-wx-sdk/utils"
)

//-----------------------------------客服消息-----------------------------------------------------------------------------

const (
	sendKfMsg = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s"
)
const (
	sendKfTextMsgTemp = `{
		"touser":"%s",
		"msgtype":"text",
		"text":
		{
			"content":"%s"
		}
	}`
	sendKfImageMsgTemp = `{
	    "touser":"%s",
	    "msgtype":"image",
	    "image":
	    {
	      "media_id":"%s"
	    }
	}`
)

func (this SmallApi) SendKfMessage(touser string, content string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(sendKfTextMsgTemp, touser, content)
	if err := this.DoPost(sendKfMsg, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this SmallApi) SendKfImageMessage(touser string, media_id string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(sendKfImageMsgTemp, touser, media_id)
	if err := this.DoPost(sendKfMsg, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

//-----------------------------------登录换取session-----------------------------------------------------------------------------
type UserSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}

const (
	jscode2SessionApi = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func (this SmallApi) Jscode2Session(code string) (*UserSession, error) {
	var res struct {
		utils.CommonError
		UserSession
	}

	if err := this.DoGetLite(jscode2SessionApi,
		&res,
		this.Context.Config.AppID,
		this.Context.Config.AppSecret,
		code); err == nil {
		if res.CommonError.IsOK() {
			return &(res.UserSession), nil
		} else {
			return nil, fmt.Errorf("%d - %s", res.ErrCode, res.ErrMsg)
		}
	} else {
		return nil, err
	}
}
