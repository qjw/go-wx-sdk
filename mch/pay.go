package mch

import (
	"github.com/qjw/go-wx-sdk/utils"
	"fmt"
)

const (
	payUnifiedorder   = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	closeUnifiedorder = "https://api.mch.weixin.qq.com/pay/closeorder"
)

type Unifiedorder struct {
	XMLName    struct{}       `xml:"xml" json:"-" sign:"-"`
	AppID      string         `xml:"appid" json:"appid" doc:"微信支付分配的公众账号ID（企业号corpid即为此appId）"`
	MchID      string         `xml:"mch_id" json:"mch_id" doc:"微信支付分配的商户号"`
	DeviceInfo utils.CharData `xml:"device_info,omitempty" json:"device_info,omitempty" doc:"自定义参数，可以为终端设备号(门店号或收银设备ID)，PC网页或公众号内支付可以传\"WEB\""`
	Nonce      string         `xml:"nonce_str" json:"nonce_str" doc:"随机字符串，长度要求在32位以内"`
	Sign       string         `xml:"sign" json:"sign" doc:"通过签名算法计算得出的签名值"`
	SignType   string         `xml:"sign_type,omitempty" json:"sign_type,omitempty" doc:"签名类型，默认为MD5，支持HMAC-SHA256和MD5。"`
	FeeType    string         `xml:"fee_type,omitempty" json:"fee_type,omitempty" doc:"标价币种,符合ISO 4217标准的三位字母代码，默认人民币：CNY"`
	TradeType  string         `xml:"trade_type" json:"trade_type" doc:"JSAPI 取值如下：JSAPI，NATIVE，APP等，说明详见参数规定"`

	UnifiedordeObj
}

type UnifiedordeObj struct {
	Body           utils.CharData `xml:"body" json:"body" doc:"商品简单描述,例如腾讯充值中心-QQ会员充值"`
	Detail         utils.CharData `xml:"detail,omitempty" json:"detail,omitempty" doc:"商品详情"`
	Attach         utils.CharData `xml:"attach,omitempty" json:"attach,omitempty" doc:"附加数据,例如深圳分店。附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。"`
	OutTradeNo     string         `xml:"out_trade_no" json:"out_trade_no" doc:"商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。详见商户订单号"`
	TotalFee       int64          `xml:"total_fee" json:"total_fee" doc:"订单总金额，单位为分"`
	SpbillCreateIp string         `xml:"spbill_create_ip,omitempty" json:"spbill_create_ip,omitempty" doc:"终端IP,APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。"`
	TimeStart      string         `xml:"time_start,omitempty" json:"time_start,omitempty" doc:"订单生成时间，格式为yyyyMMddHHmmss"`
	TimeExpire     string         `xml:"time_expire,omitempty" json:"time_expire,omitempty" doc:"交易结束时间，格式为yyyyMMddHHmmss"`
	GoodsTag       string         `xml:"goods_tag,omitempty" json:"goods_tag,omitempty" doc:"订单优惠标记，使用代金券或立减优惠功能时需要的参数，说明详见代金券或立减优惠"`
	NotifyUrl      utils.CharData `xml:"notify_url,omitempty" json:"notify_url,omitempty" doc:"异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数"`
	ProductID      string         `xml:"product_id,omitempty" json:"product_id,omitempty" doc:"trade_type=NATIVE时（即扫码支付），此参数必传。此参数为二维码中包含的商品ID，商户自行定义。"`
	LimitPay       string         `xml:"limit_pay,omitempty" json:"limit_pay,omitempty" doc:"上传此参数no_credit--可限制用户不能使用信用卡支付"`
	OpenID         string         `xml:"openid,omitempty" json:"openid,omitempty" doc:"trade_type=JSAPI时（即公众号支付），此参数必传，此参数为微信用户在商户对应appid下的唯一标识"`
}

type UnifiedordeRes struct {
	OrderCommonError
	AppID      string `xml:"appid" json:"appid" doc:"调用接口提交的公众账号ID"`
	MchID      string `xml:"mch_id" json:"mch_id" doc:"调用接口提交的商户号"`
	DeviceInfo string `xml:"device_info,omitempty" json:"device_info,omitempty" doc:"自定义参数，可以为请求支付的终端设备号等"`
	Nonce      string `xml:"nonce_str" json:"nonce_str" doc:"微信返回的随机字符串"`
	Sign       string `xml:"sign" json:"sign" doc:"微信返回的签名值"`
	ResultCode string `xml:"result_code,omitempty" json:"result_code,omitempty" doc:"业务结果"`
	ErrCode    string `xml:"err_code" json:"err_code"`
	ErrCodeDes string `xml:"err_code_des" json:"err_code_des" doc:"错误信息描述"`
	TradeType  string `xml:"trade_type" json:"trade_type" doc:"交易类型，取值为：JSAPI，NATIVE，APP等"`
	PrepayID   string `xml:"prepay_id" json:"prepay_id" doc:"微信生成的预支付会话标识，用于后续接口调用中使用，该值有效期为2小时"`
	CodeUrl    string `xml:"code_url" json:"code_url" doc:"trade_type为NATIVE时有返回，用于生成二维码，展示给用户进行扫码支付"`
}

func (this PayApi) Sign(variable interface{}) (sign string, err error) {
	ss := &utils.SignStruct{
		ToLower: false,
		Tag:     "xml",
	}
	sign, err = ss.Sign(variable, nil, this.Context.Config.ApiKey)
	return
}

func (this PayApi) CreateUnifiedOrder(param *UnifiedordeObj) (*UnifiedordeRes, error) {
	var realParam Unifiedorder
	realParam.AppID = this.Context.Config.AppID
	realParam.MchID = this.Context.Config.MchID

	realParam.DeviceInfo = "WEB"
	realParam.FeeType = "CNY"
	realParam.TradeType = "JSAPI"
	realParam.Nonce = utils.RandString(30)
	realParam.SignType = "MD5"

	realParam.UnifiedordeObj = *param

	sign, err := this.Sign(&realParam)
	if err != nil {
		return nil, err
	}
	realParam.Sign = sign

	var res UnifiedordeRes
	if err := this.DoPostObject(payUnifiedorder, &realParam, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type OrderCloseObj struct {
	OrderCloseParam

	XMLName  struct{} `xml:"xml" json:"-" sign:"-"`
	AppID    string   `xml:"appid" json:"appid" doc:"微信支付分配的公众账号ID（企业号corpid即为此appId）"`
	MchID    string   `xml:"mch_id" json:"mch_id" doc:"微信支付分配的商户号"`
	Sign     *string  `xml:"sign" json:"sign" doc:"通过签名算法计算得出的签名值"`
	SignType *string  `xml:"sign_type,omitempty" json:"sign_type,omitempty" doc:"签名类型，默认为MD5，支持HMAC-SHA256和MD5。"`
	Nonce    string   `xml:"nonce_str" json:"nonce_str" doc:"随机字符串，长度要求在32位以内"`
}

type OrderCloseParam struct {
	OutTradeNo string `xml:"out_trade_no" json:"out_trade_no" doc:"商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。详见商户订单号"`
}

type UnifiedordeCloseRes struct {
	OrderCommonError
	AppID      string `xml:"appid" json:"appid" doc:"调用接口提交的公众账号ID"`
	MchID      string `xml:"mch_id" json:"mch_id" doc:"调用接口提交的商户号"`
	Nonce      string `xml:"nonce_str" json:"nonce_str" doc:"微信返回的随机字符串"`
	Sign       string `xml:"sign" json:"sign" doc:"微信返回的签名值"`
	ResultCode string `xml:"result_code,omitempty" json:"result_code,omitempty" doc:"业务结果"`
	ResultMsg  string `xml:"result_msg,omitempty" json:"result_msg,omitempty" doc:"业务结果"`
	ErrCode    string `xml:"err_code" json:"err_code"`
	ErrCodeDes string `xml:"err_code_des" json:"err_code_des" doc:"错误信息描述"`
}

func (this PayApi) CloseUnifiedOrder(param *OrderCloseParam) (*UnifiedordeCloseRes, error) {
	var realParam OrderCloseObj
	realParam.AppID = this.Context.Config.AppID
	realParam.MchID = this.Context.Config.MchID
	realParam.Nonce = utils.RandString(30)
	realParam.SignType = utils.String("MD5")
	realParam.OrderCloseParam = *param

	sign, err := this.Sign(&realParam)
	if err != nil {
		return nil, err
	}
	realParam.Sign = &sign

	var res UnifiedordeCloseRes
	if err := this.DoPostObject(closeUnifiedorder, &realParam, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

// xml的字段名不要修改
type JsUnifiedOrder struct {
	AppID     string  `xml:"appId" json:"appid" doc:"调用接口提交的公众账号ID"`
	Nonce     string  `xml:"nonceStr" json:"nonce_str" doc:"随机字符串，长度要求在32位以内"`
	Sign      *string `xml:"paySign" json:"sign,omitempty" doc:"通过签名算法计算得出的签名值"`
	SignType  string  `xml:"signType" json:"sign_type" doc:"签名类型，默认为MD5，支持HMAC-SHA256和MD5。"`
	Package   string  `xml:"package" json:"package"`
	TimeStamp int64   `xml:"timeStamp" json:"timestamp"`
}

type JsUnifiedOrderRes struct {
	OrderCommonError
	JsUnifiedOrder
}

func (this PayApi) CreateJsUnifiedOrder(param *UnifiedordeRes) (*JsUnifiedOrderRes, error) {
	var realRes JsUnifiedOrderRes
	realRes.AppID = this.Context.Config.AppID
	realRes.Nonce = utils.RandString(30)
	realRes.SignType = "MD5"
	realRes.TimeStamp = utils.GetCurrTs()
	realRes.Package = fmt.Sprintf("prepay_id=%s", param.PrepayID)

	sign, err := this.Sign(&realRes.JsUnifiedOrder)
	if err != nil {
		return nil, err
	}
	realRes.Sign = &sign
	return &realRes, nil
}

type UnifiedOrderNotifyObj struct {
	OrderCommonError
	AppID              string `xml:"appid" json:"appid" doc:"微信支付分配的公众账号ID（企业号corpid即为此appId）"`
	MchID              string `xml:"mch_id" json:"mch_id" doc:"微信支付分配的商户号"`
	DeviceInfo         string `xml:"device_info,omitempty" json:"device_info,omitempty" doc:"自定义参数，可以为终端设备号(门店号或收银设备ID)，PC网页或公众号内支付可以传\"WEB\""`
	Nonce              string `xml:"nonce_str" json:"nonce_str" doc:"随机字符串，长度要求在32位以内"`
	Sign               string `xml:"sign" json:"sign" doc:"通过签名算法计算得出的签名值"`
	SignType           string `xml:"sign_type,omitempty" json:"sign_type,omitempty" doc:"签名类型，默认为MD5，支持HMAC-SHA256和MD5。"`
	ResultCode         string `xml:"result_code,omitempty" json:"result_code,omitempty" doc:"业务结果"`
	ErrCode            string `xml:"err_code" json:"err_code"`
	ErrCodeDes         string `xml:"err_code_des" json:"err_code_des" doc:"错误信息描述"`
	OpenID             string `xml:"openid" json:"openid" doc:"用户在商户appid下的唯一标识"`
	IsSubscribe        string `xml:"is_subscribe" json:"is_subscribe" doc:"是否关注公众账号(Y/N)"`
	TradeType          string `xml:"trade_type" json:"trade_type" doc:"JSAPI 取值如下：JSAPI，NATIVE，APP等，说明详见参数规定"`
	BankType           string `xml:"bank_type" json:"bank_type" doc:"银行类型，采用字符串类型的银行标识"`
	TotalFee           int64  `xml:"total_fee" json:"total_fee" doc:"订单总金额，单位为分"`
	SettlementTotalFee int64  `xml:"settlement_total_fee" json:"settlement_total_fee" doc:"应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。"`
	CashFee            int64  `xml:"cash_fee" json:"cash_fee" doc:"现金支付金额订单现金支付金额"`
	CashFeeType        string `xml:"cash_fee_type" json:"cash_fee_type" doc:"现金支付货币类型	，符合ISO4217标准的三位字母代码，默认人民币：CNY"`
	CouponFee          int64  `xml:"coupon_fee" json:"coupon_fee" doc:"总代金券金额,代金券金额<=订单金额，订单金额-代金券金额=现金支付金额"`
	CouponCount        int64  `xml:"coupon_count" json:"coupon_count" doc:"代金券使用数量"`
	TransactionID      string `xml:"transaction_id" json:"transaction_id" doc:"微信支付订单号	"`
	TimeEnd            string `xml:"time_end,omitempty" json:"time_end,omitempty" doc:"支付完成时间，格式为yyyyMMddHHmmss"`
	Attach             string `xml:"attach,omitempty" json:"attach,omitempty" doc:"附加数据,例如深圳分店。附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。"`
	OutTradeNo         string `xml:"out_trade_no" json:"out_trade_no" doc:"商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。详见商户订单号"`
	FeeType            string `xml:"fee_type,omitempty" json:"fee_type,omitempty" doc:"标价币种,符合ISO 4217标准的三位字母代码，默认人民币：CNY"`
}
