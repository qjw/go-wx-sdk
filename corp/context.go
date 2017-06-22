package corp

import (
	"sync"
	"fmt"
	"encoding/json"
	"time"
	"github.com/qjw/go-wx-sdk/utils"
	"github.com/qjw/go-wx-sdk/cache"
)

const (
	accessTokenKey = "corp_access_token_%s_%s"
	jsTicketKey = "js_ticket_%s_%s"
	accessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	jsTicketUrl    = "https://qyapi.weixin.qq.com/cgi-bin/get_jsapi_ticket?access_token=%s"
)

// Context struct
type Context struct {
	// 配置
	Config *Config

	// 缓存处理器
	Cache cache.Cache

	//accessTokenLock 读写锁 同一个AppID一个
	accessTokenLock *sync.RWMutex

	//jsAPITicket 读写锁 同一个AppID一个
	jsAPITicketLock *sync.RWMutex
}

func (ctx *Context) setJsAPITicketLock(lock *sync.RWMutex) {
	ctx.jsAPITicketLock = lock
}

func NewContext(config *Config,cache cache.Cache) *Context{
	context := &Context{
		Config:config,
		Cache:cache,
	}
	context.setAccessTokenLock(new(sync.RWMutex))
	context.setJsAPITicketLock(new(sync.RWMutex))
	return context
}

//SetAccessTokenLock 设置读写锁（一个appID一个读写锁）
func (ctx *Context) setAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

//GetAccessToken 获取access_token
func (ctx *Context) GetAccessToken() (accessToken string, err error) {
	ctx.accessTokenLock.Lock()
	defer ctx.accessTokenLock.Unlock()

	accessTokenCacheKey := fmt.Sprintf(accessTokenKey, ctx.Config.CorpID, ctx.Config.Tag)
	val := ctx.Cache.Get(accessTokenCacheKey)
	if val != nil {
		accessToken = val.(string)
		return
	}

	//从微信服务器获取
	var resAccessToken *utils.ResAccessToken
	resAccessToken, err = ctx.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resAccessToken.AccessToken
	return
}

//GetAccessTokenFromServer 强制从微信服务器获取token
func (ctx *Context) GetAccessTokenFromServer() (resAccessToken* utils.ResAccessToken, err error) {
	url := fmt.Sprintf(accessTokenURL,
		ctx.Config.CorpID,
		ctx.Config.CorpSecret)

	body, _, err := utils.HTTPGet(url)
	var accessToken utils.ResAccessToken
	resAccessToken = &accessToken
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}

	// 企业号=="" ，企业微信=="ok"
	if resAccessToken.ErrMsg != "ok" && resAccessToken.ErrMsg != "" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v",
			resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	accessTokenCacheKey := fmt.Sprintf(accessTokenKey, ctx.Config.CorpID, ctx.Config.Tag)
	expires := resAccessToken.ExpiresIn - 1500
	err = ctx.Cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
	return
}

func (ctx *Context) GetJsTicket() (jsTicket string, err error) {
	ctx.jsAPITicketLock.Lock()
	defer ctx.jsAPITicketLock.Unlock()

	jsTicketCacheKey := fmt.Sprintf(jsTicketKey, ctx.Config.CorpID, ctx.Config.Tag)
	val := ctx.Cache.Get(jsTicketCacheKey)
	if val != nil {
		jsTicket = val.(string)
		return
	}

	//从微信服务器获取
	var resJsTicket *utils.ResJsTicket
	resJsTicket, err = ctx.GetJsTicketFromServer()
	if err != nil {
		return
	}

	jsTicket = resJsTicket.Ticket
	return
}

func (ctx *Context) GetJsTicketFromServer() (resJsTicket *utils.ResJsTicket, err error) {
	var token string
	token,err = ctx.GetAccessToken()
	if err != nil{
		return
	}

	url := fmt.Sprintf(jsTicketUrl, token)
	var jsticket utils.ResJsTicket
	resJsTicket = &jsticket

	body, _, err := utils.HTTPGet(url)
	err = json.Unmarshal(body, &jsticket)
	if err != nil {
		return
	}
	if resJsTicket.ErrCode != 0 || resJsTicket.ErrMsg != "ok"{
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v",
			resJsTicket.ErrCode, resJsTicket.ErrMsg)
		return
	}

	jsTicketCacheKey := fmt.Sprintf(jsTicketKey, ctx.Config.CorpID, ctx.Config.Tag)
	expires := resJsTicket.ExpiresIn - 1500
	err = ctx.Cache.Set(jsTicketCacheKey, resJsTicket.Ticket, time.Duration(expires)*time.Second)
	return
}