package corp

import (
	"github.com/qjw/go-wx-sdk/utils"
)


type Config struct {
	// 企业号corpid
	CorpID string
	// 企业号App密钥
	CorpSecret string
	// 因为同一个企业号会有多个Secret，这里用于区分
	Tag string
}

type AgentConfig struct {
	AgentID int64
	// 企业号token
	Token string
	// 企业号消息加密密钥
	EncodingAESKey string
}

type KfConfig struct {
	// 企业号token
	Token string
	// 企业号消息加密密钥
	EncodingAESKey string
}

type CorpApi struct {
	utils.ApiTokenBase
	Context *Context
}

func NewCorpApi(context *Context) *CorpApi{
	api := &CorpApi{
		Context:context,
	}
	api.ContextToken = context
	return api
}
