package mch

import (
	"github.com/qjw/go-wx-sdk/utils"
)

type OrderCommonError struct {
	ReturnCode string `xml:"return_code,omitempty" json:"return_code,omitempty"`
	ReturnMsg  string `xml:"return_msg,omitempty" json:"return_msg,omitempty"`
}

type Config struct {
	AppID  string `json:"appid" doc:"微信支付分配的公众账号ID（企业号corpid即为此appId）"`
	MchID  string `json:"mch_id" doc:"微信支付分配的商户号"`
	ApiKey string `json:"apiKey" doc:"微信支付key"`
}

// Context struct
type Context struct {
	// 配置
	Config *Config
}

func NewContext(config *Config) *Context {
	context := &Context{
		Config: config,
	}
	return context
}

type PayApi struct {
	utils.ApiBaseXml
	Context *Context
}

func NewPayApi(context *Context) *PayApi {
	api := &PayApi{
		Context: context,
	}
	return api
}