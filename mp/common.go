package mp

import (
	"github.com/qjw/go-wx-sdk/utils"
)


type Config struct {
	// 微信公众号app id
	AppID          string
	// 微信公众号密钥
	AppSecret      string
	// 微信公众号token
	Token          string
	// 微信公众号消息加密密钥
	EncodingAESKey string
}


type WechatApi struct {
	utils.ApiTokenBase
	Context *Context
}

func NewWechatApi(context *Context) * WechatApi{
	api := &WechatApi{
		Context:context,
	}
	api.ContextToken = context
	return api
}
