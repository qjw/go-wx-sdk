package mp

import "github.com/qjw/go-wx-sdk/utils"

const (
	getPoiList    = "https://api.weixin.qq.com/cgi-bin/poi/getpoilist?access_token=%s"
	getPoi        = "http://api.weixin.qq.com/cgi-bin/poi/getpoi?access_token=%s"
	getWxCategory = "http://api.weixin.qq.com/cgi-bin/poi/getwxcategory?access_token=%s"
)

type PoiObj struct {
	Sid           string `json:"sid,omitempty" doc:"商户自己的id，用于后续审核通过收到poi_id 的通知时，做对应关系。请商户自己保证唯一识别性"`
	Business_name string `json:"business_name" doc:"门店名称（仅为商户名，如：国美、麦当劳，不应包含地区、地址、分店名等信息，错误示例：北京国美）,不能为空，15个汉字或30个英文字符内"`
	Branch_name   string `json:"branch_name,omitempty" doc:"分店名称（不应包含地区信息，不应与门店名有重复，错误示例：北京王府井店）20个字以内"`
	Address       string `json:"address" doc:"门店所在的详细街道地址（不要填写省市信息）"`
	Telephone     string `json:"telephone" doc:"门店的电话（纯数字，区号、分机号均由“-”隔开）"`
	City          string `json:"city" doc:"门店所在的城市,10个字以内"`
	Introduction  string `json:"introduction" doc:"商户简介，主要介绍商户信息等 300字以内"`
	Province      string `json:"province" doc:"门店所在的省份（直辖市填城市名,如：北京市）10个字以内"`
	District   string   `json:"district" doc:"门店所在地区,10个字以内"`

	Recommend  string   `json:"recommend" doc:"推荐品，餐厅可为推荐菜；酒店为推荐套房；景点为推荐游玩景点等，针对自己行业的推荐内容200字以内"`
	Special    string   `json:"special" doc:"特色服务，如免费wifi，免费停车，送货上门等商户能提供的特色功能或服务"`
	OpenTime   string   `json:"open_time" doc:"营业时间，24 小时制表示，用“-”连接，如 8:00-20:00"`
	PoiID      string   `json:"poi_id" doc:"微信内部的id"`
	Categories []string `json:"categories" doc:"门店的类型（不同级分类用“,”隔开，如：美食，川菜，火锅。详细分类参见附件：微信门店类目表）"`
	OffsetType int      `json:"offset_type" doc:"坐标类型：1 为火星坐标,2 为sogou经纬度,3 为百度经纬度,4 为mapbar经纬度,5 为GPS坐标,6 为sogou墨卡托坐标,"`
	Longitude  float64  `json:"longitude" doc:"门店所在地理位置的经度"`
	Latitude   float64  `json:"latitude" doc:"门店所在地理位置的纬度（经纬度均为火星坐标，最好选用腾讯地图标记的坐标）"`
	PhotoList  []struct {
		PhotoUrl string `json:"photo_url"`
	} `json:"photo_list" doc:"图片列表，url 形式，可以有多张图片，尺寸为640*340px。必须为上一接口生成的url。图片内容不允许与门店不相关，不允许为二维码、员工合照（或模特肖像）、营业执照、无门店正门的街景、地图截图、公交地铁站牌、菜单截图等"`
	AvgPrice       int `json:"avg_price" doc:"人均价格，大于0 的整数"`
	AvailableState int `json:"available_state" doc:"门店是否可用状态。1 表示系统错误、2 表示审核中、3 审核通过、4 审核驳回。当该字段为1、2、4 状态时，poi_id 为空"`
	UpdateStatus   int `json:"update_status" doc:"扩展字段是否正在更新中。1 表示扩展字段正在更新中，尚未生效，不允许再次更新； 0 表示扩展字段没有在更新中或更新已生效，可以再次更新"`
}

type GetPoiListParam struct {
	Begin int `json:"begin" doc:"开始位置，0 即为从第一条开始查询"`
	Limit int `json:"limit" doc:"返回数据条数，最大允许50，默认为20"`
}

type GetPoiListResp struct {
	utils.CommonError
	BusinessList []struct {
		BaseInfo PoiObj `json:"base_info"`
	} `json:"business_list"`
}

func (this WechatApi) GetPoiList(param *GetPoiListParam) (*GetPoiListResp, error) {
	var res GetPoiListResp
	if err := this.DoPostObject(getPoiList, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type GetPoiDetailParam struct {
	PoiID      string   `json:"poi_id" doc:"微信内部的id"`
}

type GetPoiDetailResp struct {
	utils.CommonError
	Business struct {
		BaseInfo PoiObj `json:"base_info"`
	} `json:"business"`
}

func (this WechatApi) GetPoiDetail(param *GetPoiDetailParam) (*GetPoiDetailResp, error) {
	var res GetPoiDetailResp
	if err := this.DoPostObject(getPoi, param, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}