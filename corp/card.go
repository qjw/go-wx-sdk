package corp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"io"
)

const (
	uploadLogo   = "https://qyapi.weixin.qq.com/cgi-bin/media/uploadimg?access_token=%s&type=card_logo"
	createCard   = "https://qyapi.weixin.qq.com/cgi-bin/card/create?access_token=%s"
	batchGetCard = "https://qyapi.weixin.qq.com/cgi-bin/card/batchget?access_token=%s"
)

type UploadCardLogoRes struct {
	utils.CommonError
	Url string `json:"url"`
}

func (this CorpApi) UploadCardLogo(reader io.Reader, filename string) (*UploadCardLogoRes, error) {
	var res UploadCardLogoRes
	if err := this.DoPostFile(reader, "media", filename, &res, uploadLogo); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type GeneralCouponObj struct {
	LogoUrl      string `json:"logo_url" doc:"卡券的商户logo"`
	BrandName    string `json:"brand_name" doc:"商户名字,字数上限为12个汉字"`
	CodeType     string `json:"code_type"`
	Title        string `json:"title" doc:"卡券名，字数上限为9个汉字。(建议涵盖卡券属性、服务及金额)。"`
	Color        string `json:"color" doc:"券颜色。按色彩规范标注填写Color010-Color102"`
	Notice       string `json:"notice" doc:"卡券使用提醒，字数上限为16个汉字。"`
	ServicePhone string `json:"service_phone,omitempty" doc:"客服电话。"`
	Description  string `json:"description" doc:"卡券使用说明，字数上限为1024个汉字。"`
	Sku          struct {
		Quantity int64 `json:"quantity" doc:"卡券库存的数量，上限为100000000。"`
	} `json:"sku"`
	GetLimit      int64 `json:"get_limit,omitempty" doc:"每人可领券的数量限制,不填写默认为50"`
	BindOpenid    bool  `json:"bind_openid,omitempty" doc:"是否指定用户领取，填写true或false。默认为false。通常指定特殊用户群体投放卡券或防止刷券时选择指定用户领取。"`
	CanShare      bool  `json:"can_share,omitempty" doc:"卡券领取页面是否可分享。"`
	CanGiveFriend bool  `json:"can_give_friend,omitempty" doc:"卡券是否可转赠。"`

	CustomUrlName     string `json:"custom_url_name" doc:"自定义入口名称"`
	CustomUrl         string `json:"custom_url" doc:"自定义入口URL"`
	CustomUrlSubTitle string `json:"custom_url_sub_title" doc:"显示在入口右侧的提示语。"`
	PromotionUrlName  string `json:"promotion_url_name" doc:"营销场景的自定义入口名称。"`
	PromotionUrl      string `json:"promotion_url" doc:"入口跳转外链的地址链接。"`
	Source            string `json:"source" doc:"第三方来源名，例如同程旅游、大众点评。"`

	DateInfo struct {
		Type           string `json:"type" doc:"使用时间的类型，DATE_TYPE_FIX_TIME_RANGE 表示固定日期区间，DATE_TYPE_FIX_TERM表示固定时长（自领取后按天算。）"`
		BeginTimestamp string `json:"begin_timestamp" doc:"type为DATE_TYPE_FIX_TIME_RANGE时专用，表示起用时间。从1970年1月1日00:00:00至起用时间的秒数，最终需转换为字符串形态传入。（东八区时间，单位为秒）"`
		EndTimestamp   string `json:"end_timestamp" doc:"type为DATE_TYPE_FIX_TIME_RANGE时专用，表示结束时间，建议设置为截止日期的23:59:59过期。（东八区时间，单位为秒）截止日期必须大于当前时间"`
		FixedTerm      string `json:"fixed_term" doc:"type为DATE_TYPE_FIX_TERM时专用，表示自领取后多少天内有效，不支持填写0。"`
		FixedBeginTerm string `json:"end_timestamp" doc:"type为DATE_TYPE_FIX_TERM时专用，表示自领取后多少天开始生效，领取后当天生效填写0。（单位为天）"`
	} `json:"date_info" doc:"使用日期，有效期的信息"`
}

type CouponCardCreateObj struct {
	Card struct {
		CardType      string `json:"card_type"`
		GeneralCoupon struct {
			BaseInfo      *GeneralCouponObj `json:"base_info"`
			DefaultDetail string            `json:"default_detail" doc:"优惠券专用，填写优惠详情。"`
		} `json:"general_coupon"`
	} `json:"card"`
}

type CouponCardCreateParam struct {
	BaseInfo      GeneralCouponObj `json:"base_info"`
	DefaultDetail string           `json:"default_detail" doc:"优惠券专用，填写优惠详情。"`
}

func newCouponObj(param *CouponCardCreateParam) *CouponCardCreateObj {
	var res CouponCardCreateObj
	res.Card.CardType = "GENERAL_COUPON"
	res.Card.GeneralCoupon.BaseInfo = &param.BaseInfo
	res.Card.GeneralCoupon.DefaultDetail = param.DefaultDetail
	return &res
}

type CouponCreateRes struct {
	utils.CommonError
	CardID string `json:"card_id"`
}

func (this CorpApi) CreateCouponCard(param *CouponCardCreateParam) (*CouponCreateRes, error) {
	var res CouponCreateRes
	if err := this.DoPostObject(createCard, newCouponObj(param), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type CardLiteObj struct {
	CardID        string `json:"card_id"`
	CardType      string `json:"card_type"`
	Title         string `json:"title"`
	Status        string `json:"status"`
	Quantity      int  `json:"quantity"`
	TotalQuantity int  `json:"total_quantity"`
	Createtime    int64  `json:"createtime"`
}

type CardListRes struct {
	utils.CommonError
	TotalNum       int64         `json:"total_num"`
	CardDigestList []CardLiteObj `json:"CardLiteObj"`
}

type CardListParam struct {
	Offset int `json:"offset"`
	Count int `json:"count"`
	Status string `json:"status,omitempty"`
}

func (this CorpApi) GetCards(param *CardListParam) (*CouponCreateRes, error) {
	var res CouponCreateRes
	if err := this.DoPostObject(batchGetCard, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}