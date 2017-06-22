package corp

import "github.com/qjw/go-wx-sdk/utils"

const (
	menuGet    = "https://qyapi.weixin.qq.com/cgi-bin/menu/get?access_token=%s&agentid=%d"
	menuCreate = "https://qyapi.weixin.qq.com/cgi-bin/menu/create?access_token=%s&agentid=%d"
	menuDelete = "https://qyapi.weixin.qq.com/cgi-bin/menu/delete?access_token=%s&agentid=%d"
)

type MenuEntryObj struct {
	Type      string          `json:"type"`
	Name      string          `json:"name"`
	Key       string          `json:"key,omitempty"`
	Url       string          `json:"url,omitempty"`
	AppID     string          `json:"appid,omitempty"`
	Pagepath  string          `json:"pagepath,omitempty"`
	MediaID   string          `json:"media_id,omitempty"`
	SubButton []*MenuEntryObj `json:"sub_button,omitempty"`
}

type MenuObj struct {
	Menu MenuCreateObj `json:"menu"`
}

type MenuCreateObj struct {
	Buttons []*MenuEntryObj `json:"button,omitempty"`
}

func (this CorpApi) GetMenu(agentid int64) (*MenuObj, error) {
	var res MenuObj
	if err := this.DoGet(menuGet, &res, agentid); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}


func (this CorpApi) CreateMenu(agentid int64,param *MenuCreateObj) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(menuCreate, param, &res, agentid); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) DeleteMenu(agentid int64) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoGet(menuDelete, &res, agentid); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}