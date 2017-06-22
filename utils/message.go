package utils

import (
	"encoding/xml"
	"errors"
)


//ErrInvalidReply 无效的回复
var ErrInvalidReply = errors.New("无效的回复消息")

//ErrUnsupportReply 不支持的回复类型
var ErrUnsupportReply = errors.New("不支持的回复消息")

// MsgType 基本消息类型
type MsgType string

// EventType 事件类型
type EventType string

const (
	//MsgTypeText 表示文本消息
	MsgTypeText MsgType = "text"
	//MsgTypeImage 表示图片消息
	MsgTypeImage = "image"
	//MsgTypeVoice 表示语音消息
	MsgTypeVoice = "voice"
	//MsgTypeVideo 表示视频消息
	MsgTypeVideo = "video"
	//MsgTypeShortVideo 表示短视频消息[限接收]
	MsgTypeShortVideo = "shortvideo"
	//MsgTypeLocation 表示坐标消息[限接收]
	MsgTypeLocation = "location"
	//MsgTypeLink 表示链接消息[限接收]
	MsgTypeLink = "link"
	//MsgTypeMusic 表示音乐消息[限回复]
	MsgTypeMusic = "music"
	//MsgTypeNews 表示图文消息[限回复]
	MsgTypeNews = "news"
	//MsgTypeTransfer 表示消息消息转发到客服
	MsgTypeTransfer = "transfer_customer_service"
	//MsgTypeEvent 表示事件推送消息
	MsgTypeEvent = "event"
)

const (
	//EventSubscribe 订阅
	EventSubscribe EventType = "subscribe"
	//EventUnsubscribe 取消订阅
	EventUnsubscribe = "unsubscribe"
	//EventScan 用户已经关注公众号，则微信会将带场景值扫描事件推送给开发者
	EventScan = "SCAN"
	//EventLocation 上报地理位置事件
	EventLocation = "LOCATION"
	//EventClick 点击菜单拉取消息时的事件推送
	EventClick = "CLICK"
	//EventView 点击菜单跳转链接时的事件推送
	EventView = "VIEW"
	//EventScancodePush 扫码推事件的事件推送
	EventScancodePush = "scancode_push"
	//EventScancodeWaitmsg 扫码推事件且弹出“消息接收中”提示框的事件推送
	EventScancodeWaitmsg = "scancode_waitmsg"
	//EventPicSysphoto 弹出系统拍照发图的事件推送
	EventPicSysphoto = "pic_sysphoto"
	//EventPicPhotoOrAlbum 弹出拍照或者相册发图的事件推送
	EventPicPhotoOrAlbum = "pic_photo_or_album"
	//EventPicWeixin 弹出微信相册发图器的事件推送
	EventPicWeixin = "pic_weixin"
	//EventLocationSelect 弹出地理位置选择器的事件推送
	EventLocationSelect = "location_select"
)


// 微信服务器请求 http body
type RequestEncryptedXMLMsg struct {
	XMLName struct{} `xml:"xml" json:"-"`

	ToUserName   string `xml:"ToUserName"`
	EncryptedMsg string `xml:"Encrypt"`
}

//ResponseEncryptedXMLMsg 需要返回的消息体
type ResponseEncryptedXMLMsg struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	EncryptedMsg string   `xml:"Encrypt"      json:"Encrypt"`
	MsgSignature string   `xml:"MsgSignature" json:"MsgSignature"`
	Timestamp    int64    `xml:"TimeStamp"    json:"TimeStamp"`
	Nonce        string   `xml:"Nonce"        json:"Nonce"`
}

type Reply interface {
	SetToUserName(toUserName string)()
	SetFromUserName(fromUserName string)()
	SetCreateTime(createTime int64)()
	SetMsgType(msgType MsgType)()
	GetMsgType()(MsgType)
}

// CommonToken 消息中通用的结构
type CommonToken struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CharData   `xml:"ToUserName"`
	FromUserName CharData   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      CharData  `xml:"MsgType"`
}

//SetToUserName set ToUserName
func (msg *CommonToken) SetToUserName(toUserName string) {
	msg.ToUserName = CharData(toUserName)
}

//SetFromUserName set FromUserName
func (msg *CommonToken) SetFromUserName(fromUserName string) {
	msg.FromUserName = CharData(fromUserName)
}

//SetCreateTime set createTime
func (msg *CommonToken) SetCreateTime(createTime int64) {
	msg.CreateTime = createTime
}

//SetMsgType set MsgType
func (msg *CommonToken) SetMsgType(msgType MsgType) {
	msg.MsgType = CharData(msgType)
}

func (msg CommonToken) GetMsgType() MsgType {
	return MsgType(msg.MsgType)
}


//Text 文本消息
type TextMessage struct {
	CommonToken
	Content CharData `xml:"Content"`
}

//NewText 初始化文本消息
func NewText(content string) *TextMessage {
	msg := &TextMessage{
		Content:CharData(content),
	}
	msg.MsgType = "text"
	return msg
}

//Text 文本消息
type VideoMessage struct {
	CommonToken
	Video struct {
		MediaId     CharData `xml:"MediaId"`
		Title       CharData `xml:"-"`
		Description CharData `xml:"-"`
	} `xml:"Video"`
}

//NewText 初始化图片消息
func NewVideo(media_id string,title string,desc string) *VideoMessage {
	msg := &VideoMessage{
		Video:struct{
			MediaId     CharData `xml:"MediaId"`
			Title       CharData `xml:"-"`
			Description CharData `xml:"-"`
		}{
			MediaId:CharData(media_id),
			Title:CharData(title),
			Description:CharData(desc),
		},
	}
	msg.MsgType = "video"
	return msg
}


//Text 文本消息
type VoiceMessage struct {
	CommonToken
	Voice struct {
		MediaId CharData `xml:"MediaId"`
	} `xml:"Voice"`
}

//NewText 初始化图片消息
func NewVoice(media_id string) *VoiceMessage {
	msg := &VoiceMessage{
		Voice:struct{
			MediaId CharData `xml:"MediaId"`
		}{
			MediaId:CharData(media_id),
		},
	}
	msg.MsgType = "voice"
	return msg
}

//Text 文本消息
type ImageMessage struct {
	CommonToken
	Image struct {
		MediaId CharData `xml:"MediaId"`
	} `xml:"Image"`
}

//NewText 初始化图片消息
func NewImage(media_id string) *ImageMessage {
	msg := new(ImageMessage)
	msg.Image.MediaId = CharData(media_id)
	msg.MsgType = "image"
	return msg
}

//ResAccessToken struct
type ResAccessToken struct {
	CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type ResJsTicket struct {
	CommonError

	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}