package mp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"encoding/xml"
)

//MixMessage 存放所有微信发送过来的消息和事件
type MixMessage struct {
	MixCommonToken

	//基本消息
	MsgID        int64   `xml:"MsgId,omitempty" json:"msg_id,omitempty"`
	Content      string  `xml:"Content,omitempty" json:"content,omitempty"`
	PicURL       string  `xml:"PicUrl,omitempty" json:"pic_url,omitempty"`
	MediaID      string  `xml:"MediaId,omitempty" json:"media_id,omitempty"`
	Format       string  `xml:"Format,omitempty" json:"format,omitempty"`
	ThumbMediaID string  `xml:"ThumbMediaId,omitempty" json:"thumb_media_id,omitempty"`
	LocationX    float64 `xml:"Location_X,omitempty" json:"location_x,omitempty"`
	LocationY    float64 `xml:"Location_Y,omitempty" json:"location_y,omitempty"`
	Scale        float64 `xml:"Scale,omitempty" json:"scale,omitempty"`
	Label        string  `xml:"Label,omitempty" json:"label,omitempty"`
	Title        string  `xml:"Title,omitempty" json:"title,omitempty"`
	Description  string  `xml:"Description,omitempty" json:"description,omitempty"`
	URL          string  `xml:"Url,omitempty" json:"url,omitempty"`

	//事件相关
	Event     utils.EventType `xml:"Event,omitempty" json:"event,omitempty"`
	EventKey  string          `xml:"EventKey,omitempty" json:"event_key,omitempty"`
	Ticket    string          `xml:"Ticket,omitempty" json:"ticket,omitempty"`
	Latitude  string          `xml:"Latitude,omitempty" json:"latitude,omitempty"`
	Longitude string          `xml:"Longitude,omitempty" json:"longitude,omitempty"`
	Precision string          `xml:"Precision,omitempty,omitempty" json:"precision,omitempty"`
	MenuID    string          `xml:"MenuId,omitempty,omitempty" json:"menu_id,omitempty"`

	ScanCodeInfo *struct {
		ScanType   string `xml:"ScanType,omitempty" json:"scan_type,omitempty"`
		ScanResult string `xml:"ScanResult,omitempty" json:"scan_result,omitempty"`
	} `xml:"ScanCodeInfo,omitempty" json:"scan_code_info,omitempty"`

	//SendPicsInfo struct {
	//	Count   int32      `xml:"Count,omitempty" json:"count,omitempty"`
	//	PicList []EventPic `xml:"PicList>item"`
	//} `xml:"SendPicsInfo"`

	//SendLocationInfo struct {
	//	LocationX float64 `xml:"Location_X,omitempty" json:"location_x,omitempty"`
	//	LocationY float64 `xml:"Location_Y,omitempty" json:"location_y,omitempty"`
	//	Scale     float64 `xml:"Scale,omitempty" json:"scale,omitempty"`
	//	Label     string  `xml:"Label,omitempty" json:"label,omitempty"`
	//	Poiname   string  `xml:"Poiname" json:"poiname,omitempty"`
	//}

	CardEvent
}

type CardEvent struct {
	// 卡券 审核事件推送 @card_pass_check/card_not_pass_check
	CardId       string `xml:"CardId,omitempty" json:"card_id,omitempty"`
	RefuseReason string `xml:"RefuseReason,omitempty" json:"refuse_reason,omitempty"`

	// 卡券 领取事件推送 @user_get_card
	IsGiveByFriend      int    `xml:"IsGiveByFriend,omitempty" json:"is_give_by_friend,omitempty" doc:"是否为转赠领取，1代表是，0代表否。"`
	UserCardCode        string `xml:"UserCardCode,omitempty" json:"user_card_code,omitempty" doc:"code序列号"`
	FriendUserName      string `xml:"FriendUserName,omitempty" json:"friend_username,omitempty" doc:"当IsGiveByFriend为1时填入的字段，表示发起转赠用户的openid"`
	OuterId             int    `xml:"OuterId,omitempty" json:"outer_id,omitempty"`
	OldUserCardCode     string `xml:"OldUserCardCode,omitempty" json:"old_user_card_code,omitempty" doc:"为保证安全，微信会在转赠发生后变更该卡券的code号，该字段表示转赠前的code。"`
	OuterStr            string `xml:"OuterStr,omitempty" json:"outer_str,omitempty"`
	IsRestoreMemberCard int    `xml:"IsRestoreMemberCard,omitempty" json:"is_restore_member_card,omitempty"`
	IsRecommendByFriend int    `xml:"IsRecommendByFriend,omitempty" json:"is_recommend_by_friend,omitempty" doc:"用户删除会员卡后可重新找回，当用户本次操作为找回时，该值为1，否则为0"`

	// 卡券 转赠事件推送 @user_gifting_card
	IsReturnBack int `xml:"IsReturnBack,omitempty" json:"is_return_back,omitempty" doc:"是否转赠退回，0代表不是，1代表是。"`
	IsChatRoom   int `xml:"IsChatRoom,omitempty" json:"is_chatroom,omitempty" doc:"是否是群转赠"`

	// 卡券 删除事件推送 @user_del_card
	// 卡券 核销事件推送 @user_consume_card
	ConsumeSource string `xml:"ConsumeSource,omitempty" json:"consume_source,omitempty" doc:"核销来源。支持开发者统计API核销（FROM_API）、公众平台核销（FROM_MP）、卡券商户助手核销（FROM_MOBILE_HELPER）（核销员微信号）"`
	LocationName  string `xml:"LocationName,omitempty" json:"location_name,omitempty" doc:"门店名称，当前卡券核销的门店名称（只有通过自助核销和买单核销时才会出现该字段）"`
	StaffOpenId   string `xml:"StaffOpenId,omitempty" json:"staff_openId,omitempty" doc:"核销该卡券核销员的openid（只有通过卡券商户助手核销时才会出现）"`
	VerifyCode    string `xml:"VerifyCode,omitempty" json:"verify_code,omitempty" doc:"自助核销时，用户输入的验证码"`
	RemarkAmount  string `xml:"RemarkAmount,omitempty" json:"remark_amount,omitempty" doc:"自助核销时，用户输入的备注金额"`

	// 卡券 买单事件推送 @user_pay_from_pay_cell
	TransId     string `xml:"TransId,omitempty" json:"trans_id,omitempty" doc:"微信支付交易订单号（只有使用买单功能核销的卡券才会出现）"`
	LocationId  int    `xml:"LocationId,omitempty" json:"location_id,omitempty" doc:"门店ID，当前卡券核销的门店ID（只有通过卡券商户助手和买单核销时才会出现）"`
	Fee         int    `xml:"Fee" json:"fee,omitempty" doc:"实付金额，单位为分"`
	OriginalFee int    `xml:"OriginalFee,omitempty" json:"original_fee,omitempty" doc:"应付金额，单位为分"`

	// 卡券 进入会员卡事件推送 @user_view_card
	// 从卡券进入公众号会话事件推送 @user_enter_session_from_card
	// 会员卡内容更新事件 @update_member_card
	ModifyBonus   int `xml:"ModifyBonus,omitempty" json:"modify_bonus,omitempty" doc:"变动的积分值。"`
	ModifyBalance int `xml:"ModifyBalance,omitempty" json:"modify_balance,omitempty" doc:"变动的余额值。"`

	// 库存报警事件 @card_sku_remind
	Detail string `xml:"Detail,omitempty" json:"detail,omitempty" doc:"报警详细信息"`

	// 券点流水详情事件 @card_pay_order
	// 会员卡激活事件推送 @submit_membercard_user_info
}

//EventPic 发图事件推送
type EventPic struct {
	PicMd5Sum string `xml:"PicMd5Sum"`
}

// CommonToken 消息中通用的结构
type MixCommonToken struct {
	XMLName      xml.Name      `xml:"xml" json:"-"`
	ToUserName   string        `xml:"ToUserName" json:"to_username"`
	FromUserName string        `xml:"FromUserName" json:"from_username"`
	CreateTime   int64         `xml:"CreateTime" json:"create_time"`
	MsgType      utils.MsgType `xml:"MsgType" json:"msg_type"`
}
