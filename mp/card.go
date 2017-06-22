package mp

import (
	"github.com/qjw/go-wx-sdk/utils"
)

const (
	batchGetCard     = "https://api.weixin.qq.com/card/batchget?access_token=%s"
	getCard          = "https://api.weixin.qq.com/card/get?access_token=%s"
	createCardQrcode = "https://api.weixin.qq.com/card/qrcode/create?access_token=%s"
	createCard       = "https://api.weixin.qq.com/card/create?access_token=%s"
	setCardWhitelist = "https://api.weixin.qq.com/card/testwhitelist/set?access_token=%s"
)

type GetCardListParam struct {
	utils.Pagination
	StatusList []string `json:"status_list" doc:"CARD_STATUS_NOT_VERIFY,待审核；CARD_STATUS_VERIFY_FAIL,审核失败；CARD_STATUS_VERIFY_OK，通过审核；CARD_STATUS_DELETE，卡券被商户删除；CARD_STATUS_DISPATCH，在公众平台投放过的卡券；"`
}

type GetCardListResp struct {
	utils.CommonError
	CardIDList []string `json:"card_id_list"`
}

func (this WechatApi) GetCardList(param *GetCardListParam) (*GetCardListResp, error) {
	var res GetCardListResp
	if err := this.DoPostObject(batchGetCard, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type CardAdvancedInfoObj struct {
	UseCondition *struct {
		AcceptCategory          string `json:"accept_category,omitempty" doc:"指定可用的商品类目，仅用于代金券类型，填入后将在券面拼写适用于xxx"`
		RejectCategory          string `json:"reject_category,omitempty" doc:"指定不可用的商品类目，仅用于代金券类型，填入后将在券面拼写不适用于xxxx"`
		CanUseWithOtherDiscount bool   `json:"can_use_with_other_discount,omitempty" doc:"不可以与其他类型共享门槛，填写false时系统将在使用须知里拼写“不可与其他优惠共享”，填写true时系统将在使用须知里拼写“可与其他优惠共享”，默认为true"`
	} `json:"use_condition,omitempty"`

	Abstract *struct {
		Abstract    string   `json:"abstract,omitempty" doc:"封面摘要简介。"`
		IconUrlList []string `json:"icon_url_list,omitempty" doc:"封面图片列表，仅支持填入一个封面图片链接，上传图片接口上传获取图片获得链接，填写非CDN链接会报错，并在此填入。建议图片尺寸像素850*350"`
	} `json:"abstract,omitempty"`

	TextImageList []struct {
		ImageUrl string `json:"image_url,omitempty" doc:"图片链接，必须调用上传图片接口上传图片获得链接，并在此填入，否则报错"`
		Text     string `json:"text,omitempty" doc:"图文描述"`
	} `json:"text_image_list,omitempty" doc:"图文列表，显示在详情内页，优惠券券开发者须至少传入一组图文列表"`

	TimeLimit []struct {
		Type        string `json:"type,omitempty" doc:"限制类型枚举值：支持填入MONDAY 周一 TUESDAY 周二 WEDNESDAY 周三 THURSDAY 周四 FRIDAY 周五 SATURDAY 周六 SUNDAY 周日 此处只控制显示，不控制实际使用逻辑，不填默认不显示"`
		BeginHour   int    `json:"begin_hour,omitempty" doc:"当前type类型下的起始时间（小时），如当前结构体内填写了MONDAY，此处填写了10，则此处表示周一 10:00可用"`
		BeginMinute int    `json:"begin_minute,omitempty" doc:"当前type类型下的起始时间（分钟），如当前结构体内填写了MONDAY，begin_hour填写10，此处填写了59，则此处表示周一 10:59可用"`
		EndHour     int    `json:"end_hour,omitempty" doc:"当前type类型下的结束时间（小时），如当前结构体内填写了MONDAY，此处填写了20，则此处表示周一 10:00-20:00可用"`
		EndMinute   int    `json:"end_minute,omitempty" doc:"当前type类型下的结束时间（分钟），如当前结构体内填写了MONDAY，begin_hour填写10，此处填写了59，则此处表示周一 10:59-00:59可用"`
	} `json:"time_limit,omitempty" doc:"使用时段限制"`
	BusinessService []string `json:"business_service,omitempty" doc:"商家服务类型：BIZ_SERVICE_DELIVER 外卖服务；BIZ_SERVICE_FREE_PARK 停车位；BIZ_SERVICE_WITH_PET 可带宠物；BIZ_SERVICE_FREE_WIFI 免费wifi，可多选"`
}

type CardBaseInfoCommonObj struct {
	LogoUrl     string `json:"logo_url" doc:"卡券的商户logo"`
	BrandName   string `json:"brand_name" doc:"商户名字,字数上限为12个汉字"`
	CodeType    string `json:"code_type" doc:"CODE_TYPE_TEXT文本；CODE_TYPE_BARCODE一维码 CODE_TYPE_QRCODE二维码CODE_TYPE_ONLY_QRCODE,二维码无code显示；CODE_TYPE_ONLY_BARCODE,一维码无code显示；CODE_TYPE_NONE，不显示code和条形码类型"`
	Title       string `json:"title" doc:"卡券名，字数上限为9个汉字。(建议涵盖卡券属性、服务及金额)。"`
	Color       string `json:"color" doc:"券颜色。按色彩规范标注填写Color010-Color102"`
	Notice      string `json:"notice" doc:"卡券使用提醒，字数上限为16个汉字。"`
	Description string `json:"description" doc:"卡券使用说明，字数上限为1024个汉字。"`

	DateInfo struct {
		Type           string `json:"type" doc:"DATE_TYPE_FIX_TIME_RANGE 表示固定日期区间，DATE_TYPE_FIX_TERM表示固定时长（自领取后按天算），DATE_TYPE_PERMANENT 表示永久有效（会员卡类型专用）"`
		BeginTimestamp int64  `json:"begin_timestamp,omitempty" doc:"type为DATE_TYPE_FIX_TIME_RANGE时专用，表示起用时间。从1970年1月1日00:00:00至起用时间的秒数，最终需转换为字符串形态传入。（东八区时间，单位为秒）"`
		EndTimestamp   int64  `json:"end_timestamp,omitempty" doc:"type为DATE_TYPE_FIX_TIME_RANGE时专用，表示结束时间，建议设置为截止日期的23:59:59过期。（东八区时间，单位为秒）截止日期必须大于当前时间"`
		FixedTerm      string `json:"fixed_term,omitempty" doc:"type为DATE_TYPE_FIX_TERM时专用，表示自领取后多少天内有效，不支持填写0。"`
		FixedBeginTerm string `json:"fixed_begin_term,omitempty" doc:"type为DATE_TYPE_FIX_TERM时专用，表示自领取后多少天开始生效，领取后当天生效填写0。（单位为天）"`
	} `json:"date_info" doc:"使用日期，有效期的信息"`

	LocationIDList  []int64 `json:"location_id_list,omitempty" doc:"支持更新适用门店列表。"`
	UseAllLocations bool    `json:"use_all_locations,omitempty" doc:"设置本卡券支持全部门店，与location_id_list互斥"`

	UseCustomCode bool   `json:"use_custom_code,omitempty" doc:"是否自定义Code码。填写true或false，默认为false。通常自有优惠码系统的开发者选择自定义Code码，并在卡券投放时带入Code码"`
	BindOpenid    bool   `json:"bind_openid,omitempty" doc:"是否指定用户领取，填写true或false。默认为false。通常指定特殊用户群体投放卡券或防止刷券时选择指定用户领取。"`
	ServicePhone  string `json:"service_phone,omitempty" doc:"客服电话。"`
	CanShare      bool   `json:"can_share,omitempty" doc:"卡券领取页面是否可分享。"`
	Source        string `json:"source,omitempty" doc:"第三方来源名，例如同程旅游、大众点评。"`

	CustomUrlName     string `json:"custom_url_name,omitempty" doc:"自定义入口名称"`
	CustomUrl         string `json:"custom_url,omitempty" doc:"自定义入口URL"`
	CustomUrlSubTitle string `json:"custom_url_sub_title,omitempty" doc:"显示在入口右侧的提示语。"`

	PromotionUrlName     string `json:"promotion_url_name,omitempty" doc:"营销场景的自定义入口名称。"`
	PromotionUrl         string `json:"promotion_url,omitempty" doc:"入口跳转外链的地址链接。"`
	PromotionUrlSubTitle string `json:"promotion_url_sub_title,omitempty" doc:"显示在营销入口右侧的提示语。"`
}

type CardBaseInfoCreateObj struct {
	CardBaseInfoCommonObj

	Sku struct {
		Quantity int64 `json:"quantity" doc:"卡券库存的数量，上限为100000000。"`
	} `json:"sku"`

	CenterTitle    string `json:"center_title,omitempty" doc:"卡券顶部居中的按钮，仅在卡券状态正常(可以核销)时显示"`
	CenterSubTitle string `json:"center_sub_title,omitempty" doc:"显示在入口下方的提示语，仅在卡券状态正常(可以核销)时显示。"`
	CenterUrl      string `json:"center_url,omitempty" doc:"顶部居中的url，仅在卡券状态正常(可以核销)时显示。"`

	GetLimit      int64 `json:"get_limit,omitempty" doc:"每人可领券的数量限制,不填写默认为50"`
	UseLimit      int64 `json:"use_limit,omitempty" doc:"每人可核销的数量限制,不填写默认为50"`
	CanGiveFriend bool  `json:"can_give_friend,omitempty" doc:"卡券是否可转赠。"`

	GetCustomCodeMode string `json:"get_custom_code_mode,omitempty" doc:"填入GET_CUSTOM_CODE_MODE_DEPOSIT表示该卡券为预存code模式卡券，须导入超过库存数目的自定义code后方可投放，填入该字段后，quantity字段须为0,须导入code后再增加库存"`
}

type CardBaseInfoObj struct {
	CardBaseInfoCommonObj

	Sku struct {
		Quantity      int64 `json:"quantity" doc:"卡券现有库存的数量"`
		TotalQuantity int64 `json:"total_quantity" doc:"卡券全部库存的数量，上限为100000000"`
	} `json:"sku"`

	Status string `json:"status" doc:"CARD_STATUS_NOT_VERIFY,待审核；CARD_STATUS_VERIFY_FAIL,审核失败；CARD_STATUS_VERIFY_OK，通过审核；CARD_STATUS_DELETE，卡券被商户删除；CARD_STATUS_DISPATCH，在公众平台投放过的卡券；"`

	ID         string `json:"id"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}

type CreateCardParam struct {
	Card struct {
		CardType string `json:"card_type" doc:"卡券类型。团购券：GROUPON; 折扣券：DISCOUNT; 礼品券：GIFT;代金券：CASH; 通用券：GENERAL_COUPON;会员卡：MEMBER_CARD; 景点门票：SCENIC_TICKET；电影票：MOVIE_TICKET； 飞机票：BOARDING_PASS；会议门票：MEETING_TICKET； 汽车票：BUS_TICKET;"`
		Groupon  *struct {
			BaseInfo     CardBaseInfoCreateObj `json:"base_info"`
			AdvancedInfo CardAdvancedInfoObj   `json:"advanced_info"`
			DealDetail   string                `json:"deal_detail" doc:"团购券专用，团购详情。"`
		} `json:"groupon,omitempty" doc:"团购券"`
		Cash *struct {
			BaseInfo     CardBaseInfoCreateObj `json:"base_info"`
			AdvancedInfo CardAdvancedInfoObj   `json:"advanced_info"`
			LeastCost    int                   `json:"least_cost" doc:"代金券专用，表示起用金额（单位为分）,如果无起用门槛则填0。"`
			ReduceCost   int                   `json:"reduce_cost" doc:"代金券专用，表示减免金额。（单位为分）"`
		} `json:"cash,omitempty" doc:"代金券"`
		Discount *struct {
			BaseInfo     CardBaseInfoCreateObj `json:"base_info"`
			AdvancedInfo CardAdvancedInfoObj   `json:"advanced_info"`
			Discount     int                   `json:"discount" doc:"折扣券专用，表示打折额度（百分比）。填30就是七折。"`
		} `json:"discount,omitempty" doc:"折扣券"`
		Gift *struct {
			BaseInfo     CardBaseInfoCreateObj `json:"base_info"`
			AdvancedInfo CardAdvancedInfoObj   `json:"advanced_info"`
			Gift         string                `json:"gift" doc:"兑换券专用，填写兑换内容的名称。"`
		} `json:"gift,omitempty" doc:"兑换券"`
		GeneralCoupon *struct {
			BaseInfo      CardBaseInfoCreateObj `json:"base_info"`
			AdvancedInfo  CardAdvancedInfoObj   `json:"advanced_info"`
			DefaultDetail string                `json:"default_detail" doc:"优惠券专用，填写优惠详情。"`
		} `json:"general_coupon,omitempty" doc:"优惠券"`
	} `json:"card"`
}

type GetCardDetailParam struct {
	CardID string `json:"card_id"`
}

type GetCardDetailResp struct {
	utils.CommonError
	Card struct {
		CardType string `json:"card_type" doc:"卡券类型。团购券：GROUPON; 折扣券：DISCOUNT; 礼品券：GIFT;代金券：CASH; 通用券：GENERAL_COUPON;会员卡：MEMBER_CARD; 景点门票：SCENIC_TICKET；电影票：MOVIE_TICKET； 飞机票：BOARDING_PASS；会议门票：MEETING_TICKET； 汽车票：BUS_TICKET;"`
		Groupon  *struct {
			BaseInfo     CardBaseInfoObj     `json:"base_info"`
			AdvancedInfo CardAdvancedInfoObj `json:"advanced_info"`
			DealDetail   string              `json:"deal_detail" doc:"团购券专用，团购详情。"`
		} `json:"groupon,omitempty" doc:"团购券"`
		Cash *struct {
			BaseInfo     CardBaseInfoObj     `json:"base_info"`
			AdvancedInfo CardAdvancedInfoObj `json:"advanced_info"`
			LeastCost    int                 `json:"least_cost" doc:"代金券专用，表示起用金额（单位为分）,如果无起用门槛则填0。"`
			ReduceCost   int                 `json:"reduce_cost" doc:"代金券专用，表示减免金额。（单位为分）"`
		} `json:"cash,omitempty" doc:"代金券"`
		Discount *struct {
			BaseInfo     CardBaseInfoObj     `json:"base_info"`
			AdvancedInfo CardAdvancedInfoObj `json:"advanced_info"`
			Discount     int                 `json:"discount" doc:"折扣券专用，表示打折额度（百分比）。填30就是七折。"`
		} `json:"discount,omitempty" doc:"折扣券"`
		Gift *struct {
			BaseInfo     CardBaseInfoObj     `json:"base_info"`
			AdvancedInfo CardAdvancedInfoObj `json:"advanced_info"`
			Gift         string              `json:"gift" doc:"兑换券专用，填写兑换内容的名称。"`
		} `json:"gift,omitempty" doc:"兑换券"`
		GeneralCoupon *struct {
			BaseInfo      CardBaseInfoObj     `json:"base_info"`
			AdvancedInfo  CardAdvancedInfoObj `json:"advanced_info"`
			DefaultDetail string              `json:"default_detail" doc:"优惠券专用，填写优惠详情。"`
		} `json:"general_coupon,omitempty" doc:"优惠券"`
	} `json:"card"`
}

func (this WechatApi) GetCardDetail(param *GetCardDetailParam) (*GetCardDetailResp, error) {
	var res GetCardDetailResp
	if err := this.DoPostObject(getCard, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type CreateCardQrcodeParam struct {
	Action_name    string `json:"action_name" doc:"QR_CARD"`
	Expire_seconds int64  `json:"expire_seconds"`
	Action_info    struct {
		Card struct {
			CardID       string `json:"card_id"`
			Code         string `json:"code,omitempty"`
			Openid       string `json:"openid,omitempty"`
			IsUniqueCode string `json:"is_unique_code,omitempty"`
			OuterID      int    `json:"outer_id,omitempty"`
			OuterStr     string `json:"outer_str,omitempty"`
		} `json:"card"`
	} `json:"action_info"`
}

type CreateCardQrcodeResp struct {
	utils.CommonError
	Ticket        string `json:"ticket"`
	Url           string `json:"url"`
	ShowQrcodeUrl string `json:"show_qrcode_url"`
}

func (this WechatApi) CreateCardQrcode(param *CreateCardQrcodeParam) (*CreateCardQrcodeResp, error) {
	var res CreateCardQrcodeResp
	if err := this.DoPostObject(createCardQrcode, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type CreateCardResp struct {
	utils.CommonError
	CardID string `json:"card_id"`
}

func (this WechatApi) CreateCard(param *CreateCardParam) (*CreateCardResp, error) {
	var res CreateCardResp
	if err := this.DoPostObject(createCard, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type SetCardWhitelistParam struct {
	OpenIDs   []string `json:"openid"`
	Usernames []string `json:"username"`
}

func (this WechatApi) SetCardWhitelist(param *SetCardWhitelistParam) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(setCardWhitelist, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

const (
	decryptCardCode     = "https://api.weixin.qq.com/card/code/decrypt?access_token=%s"
	deleteCard          = "https://api.weixin.qq.com/card/delete?access_token=%s"
	checkCardCode       = "https://api.weixin.qq.com/card/code/get?access_token=%s"
	getCardUseList      = "https://api.weixin.qq.com/card/user/getcardlist?access_token=%s"
	unavailableCardCode = "https://api.weixin.qq.com/card/code/unavailable?access_token=%s"
	// getcardbizuininfo = "https://api.weixin.qq.com/datacube/getcardbizuininfo?access_token=%s"
)

type GetCardUseListParam struct {
	CardID string `json:"card_id,omitempty"`
	OpenID string `json:"openid" doc:"用户openid"`
}

type GetCardUseListResp struct {
	utils.CommonError
	CardList []struct {
		CardIDParam
		Code string `json:"code"`
	} `json:"card_list"`
	HasShareCard bool `json:"has_share_card" doc:"是否有可用的朋友的券"`
}

func (this WechatApi) GetCardUseList(param *GetCardUseListParam) (*GetCardUseListResp, error) {
	var res GetCardUseListResp
	if err := this.DoPostObject(getCardUseList, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type CheckCardCodeParam struct {
	CardID       string `json:"card_id,omitempty"`
	Code         string `json:"code"`
	CheckConsume bool   `json:"check_consume" doc:"是否校验code核销状态，填入true和false时的code异常状态返回数据不同。"`
}

type CheckCardCodeResp struct {
	utils.CommonError
	Card struct {
		CardID    string `json:"card_id"`
		BeginTime int64  `json:"begin_time"`
		EndTime   int64  `json:"end_time"`
	} `json:"card"`
	OpenID         string `json:"openid" doc:"用户openid"`
	CanConsume     bool   `json:"can_consume" doc:"是否可以核销，true为可以核销，false为不可核销"`
	OuterStr       string `json:"outer_str"`
	UserCardStatus string `json:"user_card_status" doc:"当前code对应卡券的状态NORMAL正常,CONSUMED已核销,EXPIRE已过期,GIFTING转赠中,GIFT_TIMEOUT转赠超时,DELETE已删除,UNAVAILABLE已失效"`
}

func (this WechatApi) CheckCardCode(param *CheckCardCodeParam) (*CheckCardCodeResp, error) {
	var res CheckCardCodeResp
	if err := this.DoPostObject(checkCardCode, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type CardOuterParam struct {
	EncryptCode string `json:"encrypt_code" doc:"加密的code"`
	CardCode    string `json:"code" doc:"实际的code"`
	CardID      string `json:"card_id"`
	OpenID      string `json:"openid"`
	OuterStr    string `json:"outer_str"`
	OuterID     int    `json:"outer_id"`
}


type DecryptCardCodeParam struct {
	EncryptCode string `json:"encrypt_code" doc:"加密的code"`
}
type DecryptCardCodeResp struct {
	utils.CommonError
	CardCode string `json:"code" doc:"实际的code"`
}

func (this WechatApi) DecryptCardCode(param *DecryptCardCodeParam) (*DecryptCardCodeResp, error) {
	var res DecryptCardCodeResp
	if err := this.DoPostObject(decryptCardCode, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type UnavailableCardCodeParam struct {
	CardID string `json:"card_id" doc:"自定义code卡券的请求 需要此字段"`
	Code   string `json:"code"`
	Reason string `json:"reason"`
}

func (this WechatApi) UnavailableCardCode(param *UnavailableCardCodeParam) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(unavailableCardCode, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type CardIDParam struct {
	CardID string `json:"card_id"`
}

/*
 注意：如用户在商家删除卡券前已领取一张或多张该卡券依旧有效。即删除卡券不能删除已被用户领取，保存在微信客户端中的卡券。
*/
func (this WechatApi) DeleteCard(param *CardIDParam) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(deleteCard, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

const (
	createCardLandingpage = "https://api.weixin.qq.com/card/landingpage/create?access_token=%s"
	consumeCardCode       = "https://api.weixin.qq.com/card/code/consume?access_token=%s"
)

type CreateCardLandingpageParam struct {
	Banner    string `json:"banner" doc:"页面的banner图片链接，须调用，建议尺寸为640*300。"`
	PageTitle string `json:"page_title" doc:"页面的title。"`
	CanShare  bool   `json:"can_share" doc:"页面是否可以分享,填入true/false"`
	Scene     string `json:"scene" doc:"投放页面的场景值；SCENE_NEAR_BY附近,SCENE_MENU自定义菜单,SCENE_QRCODE二维码,SCENE_ARTICLE公众号文章,SCENE_H5 h5页面,SCENE_IVR 自动回复,SCENE_CARD_CUSTOM_CELL 卡券自定义cell"`
	CardList  []struct {
		CardID   string `json:"card_id"`
		ThumbUrl string `json:"thumb_url" doc:"缩略图url"`
	} `json:"card_list"`
}

type CreateCardLandingpageResp struct {
	utils.CommonError
	Url    string `json:"url" doc:"货架链接。"`
	PageID int64  `json:"page_id" doc:"货架ID。货架的唯一标识。"`
}

func (this WechatApi) CreateCardLandingpage(param *CreateCardLandingpageParam) (*CreateCardLandingpageResp, error) {
	var res CreateCardLandingpageResp
	if err := this.DoPostObject(createCardLandingpage, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type ConsumeCardCodeParam struct {
	CardID string `json:"card_id,omitempty" doc:"自定义code卡券的请求 需要此字段"`
	Code   string `json:"code"`
}

type ConsumeCardCodeResp struct {
	utils.CommonError
	Openid string `json:"openid" doc:"用户在该公众号内的唯一身份标识。"`
	Card   struct {
		CardID string `json:"card_id" doc:"卡券ID。"`
	} `json:"card"`
}

func (this WechatApi) ConsumeCardCode(param *ConsumeCardCodeParam) (*ConsumeCardCodeResp, error) {
	var res ConsumeCardCodeResp
	if err := this.DoPostObject(consumeCardCode, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}
