package small

import (
	"encoding/xml"
	"github.com/qjw/go-wx-sdk/utils"
)

//MixMessage 存放所有微信发送过来的消息和事件
type MixMessage struct {
	MixCommonToken

	//基本消息
	MsgID   int64  `xml:"MsgId,omitempty" json:"msg_id,omitempty"`
	Content string `xml:"Content,omitempty" json:"content,omitempty"`
	PicURL  string `xml:"PicUrl,omitempty" json:"pic_url,omitempty"`
	MediaID string `xml:"MediaId,omitempty" json:"media_id,omitempty"`

	//事件相关
	Event       utils.EventType `xml:"Event,omitempty" json:"event,omitempty"`
	// 进入会话事件  @user_enter_tempsession
	SessionFrom string          `xml:"SessionFrom,omitempty" json:"session_from,omitempty"`
}

// CommonToken 消息中通用的结构
type MixCommonToken struct {
	XMLName      xml.Name      `xml:"xml" json:"-"`
	ToUserName   string        `xml:"ToUserName" json:"to_username"`
	FromUserName string        `xml:"FromUserName" json:"from_username"`
	CreateTime   int64         `xml:"CreateTime" json:"create_time"`
	MsgType      utils.MsgType `xml:"MsgType" json:"msg_type"`
}
