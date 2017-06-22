package mp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"github.com/qjw/go-wx-sdk/cache"
)

const (
	//AccessTokenURL 获取access_token的接口
	accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"
	jsTicketUrl    = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
	cardTicketUrl  = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=wx_card"
)

const (
	jsTickctTemp = "js_ticket_%s"
	cardTickctTemp = "card_ticket_%s"
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

// SetJsAPITicketLock 设置jsAPITicket的lock
func (ctx *Context) setJsAPITicketLock(lock *sync.RWMutex) {
	ctx.jsAPITicketLock = lock
}

//SetAccessTokenLock 设置读写锁（一个appID一个读写锁）
func (ctx *Context) setAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

func NewContext(config *Config, cache cache.Cache) *Context {
	context := &Context{
		Config: config,
		Cache:  cache,
	}
	context.setAccessTokenLock(new(sync.RWMutex))
	context.setJsAPITicketLock(new(sync.RWMutex))
	return context
}

//GetAccessToken 获取access_token
func (ctx *Context) GetAccessToken() (accessToken string, err error) {
	ctx.accessTokenLock.Lock()
	defer ctx.accessTokenLock.Unlock()

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", ctx.Config.AppID)
	val := ctx.Cache.Get(accessTokenCacheKey)
	if val != nil {
		accessToken = val.(string)
		return
	}

	//从微信服务器获取
	var resAccessToken utils.ResAccessToken
	resAccessToken, err = ctx.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = resAccessToken.AccessToken
	return
}

//GetAccessTokenFromServer 强制从微信服务器获取token
func (ctx *Context) GetAccessTokenFromServer() (resAccessToken utils.ResAccessToken, err error) {
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", accessTokenURL,
		ctx.Config.AppID,
		ctx.Config.AppSecret)

	body, _, err := utils.HTTPGet(url)
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrMsg != "" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v",
			resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", ctx.Config.AppID)
	expires := resAccessToken.ExpiresIn - 1500
	err = ctx.Cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
	return
}

func (ctx *Context) GetJsTicket() (jsTicket string, err error) {
	return ctx.getTicket(jsTicketUrl,jsTickctTemp)
	/*
	ctx.jsAPITicketLock.Lock()
	defer ctx.jsAPITicketLock.Unlock()

	jsTicketCacheKey := fmt.Sprintf("js_ticket_%s", ctx.Config.AppID)
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
	*/
}

func (ctx *Context) GetCardTicket() (jsTicket string, err error) {
	return ctx.getTicket(cardTicketUrl, cardTickctTemp)
}

/*
func (ctx *Context) GetJsTicketFromServer() (resJsTicket *utils.ResJsTicket, err error) {
	var token string
	token, err = ctx.GetAccessToken()
	if err != nil {
		return
	}

	url := fmt.Sprintf(jsTicketUrl, token)
	var jsticket utils.ResJsTicket
	resJsTicket = &jsticket
	var body []byte
	body, err = utils.HTTPGet(url)
	err = json.Unmarshal(body, &jsticket)
	if err != nil {
		return
	}
	if resJsTicket.ErrCode != 0 || resJsTicket.ErrMsg != "ok" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v",
			resJsTicket.ErrCode, resJsTicket.ErrMsg)
		return
	}

	jsTicketCacheKey := fmt.Sprintf("js_ticket_%s", ctx.Config.AppID)
	expires := resJsTicket.ExpiresIn - 1500
	err = ctx.Cache.Set(jsTicketCacheKey, resJsTicket.Ticket, time.Duration(expires)*time.Second)
	return
}
*/


func (ctx *Context) getTicket(urlTemp, keyTemp string) (jsTicket string, err error) {
	ctx.jsAPITicketLock.Lock()
	defer ctx.jsAPITicketLock.Unlock()

	jsTicketCacheKey := fmt.Sprintf(keyTemp, ctx.Config.AppID)
	val := ctx.Cache.Get(jsTicketCacheKey)
	if val != nil {
		jsTicket = val.(string)
		return
	}

	//从微信服务器获取
	var resJsTicket *utils.ResJsTicket
	resJsTicket, err = ctx.getTicketFromServer(urlTemp,keyTemp)
	if err != nil {
		return
	}

	jsTicket = resJsTicket.Ticket
	return
}

func (ctx *Context) getTicketFromServer(urlTemp, keyTemp string) (resJsTicket *utils.ResJsTicket, err error,) {
	var token string
	token, err = ctx.GetAccessToken()
	if err != nil {
		return
	}

	url := fmt.Sprintf(urlTemp, token)
	var jsticket utils.ResJsTicket
	resJsTicket = &jsticket

	body, _, err := utils.HTTPGet(url)
	err = json.Unmarshal(body, &jsticket)
	if err != nil {
		return
	}
	if resJsTicket.ErrCode != 0 || resJsTicket.ErrMsg != "ok" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v",
			resJsTicket.ErrCode, resJsTicket.ErrMsg)
		return
	}

	jsTicketCacheKey := fmt.Sprintf(keyTemp, ctx.Config.AppID)
	expires := resJsTicket.ExpiresIn - 1500
	err = ctx.Cache.Set(jsTicketCacheKey, resJsTicket.Ticket, time.Duration(expires)*time.Second)
	return
}
