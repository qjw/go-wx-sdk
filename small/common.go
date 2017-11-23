package small

import "github.com/qjw/go-wx-sdk/utils"

type Config struct {
	// 小程序app id
	AppID          string
	// 小程序密钥
	AppSecret      string
	// 小程序token
	Token          string
	// 小程序消息加密密钥
	EncodingAESKey string
}

type SmallApi struct {
	utils.ApiTokenBase
	Context *Context
}

func NewSmallApi(context *Context) * SmallApi{
	api := &SmallApi{
		Context:context,
	}
	api.ContextToken = context
	return api
}
