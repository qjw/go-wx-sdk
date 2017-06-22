package corp

import "encoding/xml"


// CommonToken 消息中通用的结构
type MixCommonToken struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `xml:"ToUserName" json:"to_username"`
	FromUserName string   `xml:"FromUserName" json:"from_username"`
	CreateTime   int64    `xml:"CreateTime" json:"create_time"`
	MsgType      string   `xml:"MsgType" json:"msg_type"`
	AgentID      int64    `xml:"AgentID" json:"agent_id"`
}

type MixMessage struct {
	MixCommonToken

	MsgID int64 `xml:"MsgId" json:"msg_id"`

	Content      string `xml:"Content" json:"content,omitempty"`
	PicURL       string `xml:"PicUrl" json:"pic_url,omitempty"`
	MediaID      string `xml:"MediaId" json:"media_id,omitempty"`
	Format       string `xml:"Format" json:"format,omitempty"`
	ThumbMediaID string `xml:"ThumbMediaId" json:"thumb_media_id,omitempty"`

	// location
	LocationX float64 `xml:"Location_X" json:"location_x,omitempty"`
	LocationY float64 `xml:"Location_Y" json:"location_y,omitempty"`
	Scale     float64 `xml:"Scale" json:"scale,omitempty"`
	Label     string  `xml:"Label" json:"label,omitempty"`

	// link
	Title       string `xml:"Title" json:"title,omitempty"`
	Description string `xml:"Description" json:"description,omitempty"`
	URL         string `xml:"Url" json:"url,omitempty"`
}
