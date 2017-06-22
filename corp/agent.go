package corp

import "github.com/qjw/go-wx-sdk/utils"

const (
	agentGet  = "https://qyapi.weixin.qq.com/cgi-bin/agent/get?access_token=%s&agentid=%d"
	agentSet  = "https://qyapi.weixin.qq.com/cgi-bin/agent/set?access_token=%s"
	agentList = "https://qyapi.weixin.qq.com/cgi-bin/agent/list?access_token=%s"
)

type AgentLiteObj struct {
	AgentID       int64  `json:"agentid" doc:"企业应用id"`
	Name          string `json:"name" doc:"企业应用名称"`
	SquareLogoUrl string `json:"square_logo_url" doc:"企业应用方形头像"`
	RoundLogoUrl  string `json:"round_logo_url" doc:"企业应用圆形头像"`
}

type AgentListRes struct {
	utils.CommonError
	AgentList []AgentLiteObj `json:"agentlist"`
}

func (this CorpApi) GetAgentList() (*AgentListRes, error) {
	var res AgentListRes
	err := this.DoGet(agentList, &res)
	if err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type AgentObj struct {
	AgentLiteObj
	Description string `json:"description" doc:"企业应用详情"`

	AllowUserinfos struct {
		Users []struct {
			Userid string `json:"userid"`
			Status string `json:"status"`
		} `json:"user"`
	} `json:"allow_userinfos,omitempty" doc:"企业应用可见范围（人员），其中包括userid和关注状态state"`

	AllowPartys struct {
		Partyid []int64 `json:"partyid"`
	} `json:"allow_partys,omitempty" doc:"企业应用可见范围（部门）"`
	AllowTags struct {
		Tagid []int64 `json:"tagid"`
	} `json:"allow_tags,omitempty" doc:"企业应用可见范围（标签）"`
	Close              int    `json:"close" doc:"企业应用是否被禁用"`
	RedirectDomain     string `json:"redirect_domain,omitempty" doc:"企业应用可信域名"`
	ReportLocationFlag int    `json:"report_location_flag" doc:"企业应用是否打开地理位置上报 0：不上报；1：进入会话上报；2：持续上报"`
	Isreportuser       int    `json:"isreportuser" doc:"是否接收用户变更通知。0：不接收；1：接收"`
	Isreportenter      int    `json:"isreportenter" doc:"是否上报用户进入应用事件。0：不接收；1：接收"`
	ChatExtensionUrl   string `json:"chat_extension_url,omitempty" doc:"关联会话url"`
	Type               int    `json:"type" doc:"应用类型。1：消息型；2：主页型"`
}

type AgentDetailRes struct {
	utils.CommonError
	AgentObj
}

func (this CorpApi) GetAgentDetail(agentid int64) (*AgentDetailRes, error) {
	var res AgentDetailRes
	err := this.DoGet(agentGet, &res, agentid)
	if err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type AgentUpdateObj struct {
	AgentID            int64   `json:"agentid" doc:"企业应用id"`
	ReportLocationFlag *int    `json:"report_location_flag,omitempty" doc:"企业应用是否打开地理位置上报 0：不上报；1：进入会话上报；2：持续上报"`
	RedirectDomain     *string `json:"redirect_domain,omitempty" doc:"企业应用可信域名"`
	Name               *string `json:"name,omitempty" doc:"企业应用名称"`
	LogoMediaid        *string `json:"logo_mediaid,omitempty" doc:"企业应用头像的mediaid，通过多媒体接口上传图片获得mediaid，上传后会自动裁剪成方形和圆形两个头像"`
	Description        *string `json:"description,omitempty" doc:"企业应用详情"`
	Isreportuser       *int    `json:"isreportuser,omitempty" doc:"是否接收用户变更通知。0：不接收；1：接收"`
	Isreportenter      *int    `json:"isreportenter,omitempty" doc:"是否上报用户进入应用事件。0：不接收；1：接收"`
	ChatExtensionUrl   *string `json:"chat_extension_url,omitempty" doc:"关联会话url"`
	HomeUrl            *string `json:"home_url,omitempty" doc:"主页型应用url。url必须以http或者https开头。消息型应用无需该参数"`
}

func (this CorpApi) UpdateAgent(agent *AgentUpdateObj) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(agentSet, agent, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}