package mp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"
)

//-----------------------------------客服消息-----------------------------------------------------------------------------

const (
	sendKfMsg       = "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s"
	getkflist       = "https://api.weixin.qq.com/cgi-bin/customservice/getkflist?access_token=%s"
	addkfaccount    = "https://api.weixin.qq.com/customservice/kfaccount/add?access_token=%s"
	delkfaccount    = "https://api.weixin.qq.com/customservice/kfaccount/del?access_token=%s&kf_account=%s"
	updatekfaccount = "https://api.weixin.qq.com/customservice/kfaccount/update?access_token=%s"
)
const (
	addKfAccountTemp = `{
	     "kf_account" : "%s",
	     "nickname" : "%s",
	     "password" : "%s",
	}`
	sendKfTextMsgTemp = `{
		"touser":"%s",
		"msgtype":"text",
		"text":
		{
			"content":"%s"
		}
	}`
	sendKfImageMsgTemp = `{
	    "touser":"%s",
	    "msgtype":"image",
	    "image":
	    {
	      "media_id":"%s"
	    }
	}`
	sendKfVideoMsgTemp = `{
	    "touser":"%s",
	    "msgtype":"video",
	    "video":
	    {
	      "media_id":"%s",
	      "thumb_media_id":"%s",
	      "title":"%s",
	      "description":"%s"
	    }
	}`
	sendKfMpnewsMsgTemp = `{
	    "touser":"%s",
	    "msgtype":"mpnews",
	    "mpnews":
	    {
		 "media_id":"%s"
	    }
	}`
	sendKfArticleMsgTemp = `{
	    "touser":"{{.ToUser}}",
	    "msgtype":"news",
	    "news":{
		"articles": [
			{{range $index, $element := .Articles}}
	   		{{if $index}},{{end}}
			{
				"title": "{{$element.Title}}",
				"description": "{{$element.Description}}",
				"url": "{{$element.Url}}",
				"picurl": "{{$element.Picurl}}"
			}
			{{end}}
		 ]
	    }
	}`
	sendKfCardMsgTemp = `{
	    "touser":"%s",
	    "msgtype":"wxcard",
	    "wxcard":
	    {
		 "card_id":"%s"
	    }
	}`
)

func (this WechatApi) UpdateKfAccount(kf_account string, nickname string, password string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(addKfAccountTemp, kf_account, nickname, password)
	if err := this.DoPost(updatekfaccount, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) DelKfAccount(kf_account string) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoGet(delkfaccount, &res, kf_account); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) AddKfAccount(kf_account string, nickname string, password string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(addKfAccountTemp, kf_account, nickname, password)
	if err := this.DoPost(addkfaccount, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type KfList struct {
	KfList []struct {
		KfAccount    string `json:"kf_account"`
		KfNick       string `json:"kf_nick"`
		KfID         int    `json:"kf_id"`
		KfHeadimgurl string `json:"kf_headimgurl"`
	} `json:"kf_list"`
}

func (this WechatApi) GetKfList() (*KfList, error) {
	var res KfList
	if err := this.DoGet(getkflist, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url,omitempty"`
	Picurl      string `json:"picurl,omitempty"`
}

func (this WechatApi) SendKfArticleMessage(touser string, articles []Article) (*utils.CommonError, error) {
	var res utils.CommonError
	ttt := template.New("SendKfArticleMessage")
	ttt.Parse(sendKfArticleMsgTemp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		ToUser   string
		Articles []Article
	}{
		ToUser:   touser,
		Articles: articles,
	})
	if err := this.DoPost(sendKfMsg, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) SendKfVideoMessage(touser string,
	media_id string,
	thumb_media_id string,
	title string,
	description string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(sendKfVideoMsgTemp, touser, media_id, thumb_media_id, title, description)
	if err := this.DoPost(sendKfMsg, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) SendKfMessage(touser string, content string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(sendKfTextMsgTemp, touser, content)
	if err := this.DoPost(sendKfMsg, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) SendKfImageMessage(touser string, media_id string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(sendKfImageMsgTemp, touser, media_id)
	if err := this.DoPost(sendKfMsg, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) SendKfMpnewsMessage(touser string, media_id string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(sendKfMpnewsMsgTemp, touser, media_id)
	if err := this.DoPost(sendKfMsg, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) SendKfCardMessage(touser string, card_id string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(sendKfCardMsgTemp, touser, card_id)
	if err := this.DoPost(sendKfMsg, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

//-----------------------------------用户和分组---------------------------------------------------------------------------

const (
	userGet            = "https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s&next_openid=%s"
	userDetailGet      = "https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN"
	tagListGet         = "https://api.weixin.qq.com/cgi-bin/tags/get?access_token=%s"
	tagUpdate          = "https://api.weixin.qq.com/cgi-bin/tags/update?access_token=%s"
	tagCreate          = "https://api.weixin.qq.com/cgi-bin/tags/create?access_token=%s"
	tagDelete          = "https://api.weixin.qq.com/cgi-bin/tags/delete?access_token=%s"
	tagUsers           = "https://api.weixin.qq.com/cgi-bin/user/tag/get?access_token=%s"
	tagUserBatchAdd    = "https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token=%s"
	tagUserBatchRm     = "https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token=%s"
	userTags           = "https://api.weixin.qq.com/cgi-bin/tags/getidlist?access_token=%s"
	userRemark         = "https://api.weixin.qq.com/cgi-bin/user/info/updateremark?access_token=%s"
	userBatchDetailGet = "https://api.weixin.qq.com/cgi-bin/user/info/batchget?access_token=%s"
)

const (
	tagUPdateTemp = `{
		"tag" : {
			"id" : %d,
			"name" : "%s"
		}
	}`
	tagCreateTemp = `{
		"tag" : {
			"name" : "%s"
		}
	}`
	tagDeleteTemp = `{
		"tag":{
			"id" : %d
		}
	}`
	tagUsersTemp = `{
		"tagid" : %d,
		"next_openid":"%s"
	}`
	tagUserBatchAddTemp = `{
		"openid_list": [
			{{range $index, $element := .OpenIDList}}
	   		{{if $index}},{{end}}
				"{{$element}}"
			{{end}}
		],
		"tagid": {{.TagID}}
	}`
	userTagsTemp = `{
		"openid" : "%s"
	}`
	userRemarkTemp = `{
		"openid":"%s",
		"remark":"%s"
	}`
	userBatchDetailGetTemp = `{
		"user_list": [
			{{range $index, $element := .OpenIDList}}
	   		{{if $index}},{{end}}
	   		{
				"openid": "{{$element}}",
				"lang": "zh_CN"
			}
			{{end}}
		]
	}`
)

type Users struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Data  struct {
		OpenID []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

func (this WechatApi) GetUser(next_openid string) (*Users, error) {
	var res Users
	if err := this.DoGet(userGet, &res, next_openid); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type UserInfo struct {
	Subscribe      int    `json:"subscribe"`
	Openid         string `json:"openid"`
	Nickname       string `json:"nickname"`
	Sex            int    `json:"sex"`
	Language       string `json:"language"`
	City           string `json:"city"`
	Province       string `json:"province"`
	Country        string `json:"country"`
	Headimgurl     string `json:"headimgurl"`
	Subscribe_time int64  `json:"subscribe_time"`
	Unionid        string `json:"unionid"`
	Remark         string `json:"remark"`
	Groupid        int    `json:"groupid"`
	Tagid_list     []int  `json:"tagid_list"`
}

type UserInfoList struct {
	UserInfoList []*UserInfo `json:"user_info_list"`
}

func (this WechatApi) GetUserDetail(openid string) (*UserInfo, error) {
	var res UserInfo
	if err := this.DoGet(userDetailGet, &res, openid); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) BatchGetUserDetail(openid_list []string) (*UserInfoList, error) {
	var res UserInfoList

	ttt := template.New("BatchGetUserDetail")
	ttt.Parse(userBatchDetailGetTemp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		OpenIDList []string
	}{
		OpenIDList: openid_list,
	})

	if err := this.DoPost(userBatchDetailGet, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type UserTag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type TagList struct {
	Tags []UserTag `json:"tags"`
}

func (this WechatApi) GetTagList() (*TagList, error) {
	var res TagList
	if err := this.DoGet(tagListGet, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) UpdateTag(id int, name string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(tagUPdateTemp, id, name)
	if err := this.DoPost(tagUpdate, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type TagCreate struct {
	Tag struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tag"`
}

func (this WechatApi) CreateTag(name string) (*TagCreate, error) {
	var res TagCreate
	body := fmt.Sprintf(tagCreateTemp, name)
	if err := this.DoPost(tagCreate, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) DeleteTag(id int) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(tagDeleteTemp, id)
	if err := this.DoPost(tagDelete, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type TagUsers struct {
	Count int `json:"count"`
	Data  struct {
		OpenID []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

func (this WechatApi) GetTagUsers(tagid int, next_openid string) (*TagUsers, error) {
	var res TagUsers
	body := fmt.Sprintf(tagUsersTemp, tagid, next_openid)
	if err := this.DoPost(tagUsers, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) AddTagMembers(tagid int, openid_list []string) (*utils.CommonError, error) {
	var res utils.CommonError
	ttt := template.New("AddTagMembers")
	ttt.Parse(tagUserBatchAddTemp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		TagID      int
		OpenIDList []string
	}{
		TagID:      tagid,
		OpenIDList: openid_list,
	})

	if err := this.DoPost(tagUserBatchAdd, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) RemoveTagMembers(tagid int, openid_list []string) (*utils.CommonError, error) {
	var res utils.CommonError

	ttt := template.New("RemoveTagMembers")
	ttt.Parse(tagUserBatchAddTemp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		TagID      int
		OpenIDList []string
	}{
		TagID:      tagid,
		OpenIDList: openid_list,
	})
	if err := this.DoPost(tagUserBatchRm, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type UserTags struct {
	TagidList []int `json:"tagid_list"`
}

func (this WechatApi) GetUserTags(openid string) (*UserTags, error) {
	var res UserTags
	body := fmt.Sprintf(userTagsTemp, openid)
	if err := this.DoPost(userTags, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) SetUserRemark(openid string, remark string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(userRemarkTemp, openid, remark)
	if err := this.DoPost(userRemark, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

//-----------------------------------黑名单------------------------------------------------------------------------------

const (
	getblacklist     = "https://api.weixin.qq.com/cgi-bin/tags/members/getblacklist?access_token=%s"
	batchunblacklist = "https://api.weixin.qq.com/cgi-bin/tags/members/batchunblacklist?access_token=%s"
	batchblacklist   = "https://api.weixin.qq.com/cgi-bin/tags/members/batchblacklist?access_token=%s"
)
const (
	getBlacklistTemp = `{
		"begin_openid":"%s"
	}`
	batchUnblacklistTemp = `{
		"openid_list":[
			{{range $index, $element := .OpenIDList}}
	   		{{if $index}},{{end}}
				"{{$element}}"
			{{end}}
		]
	}`
)

type Blacklist struct {
	Total      int    `json:"total"`
	Count      int    `json:"count"`
	NextOpenid string `json:"next_openid"`
	Data       struct {
		Openid []string `json:"openid"`
	} `json:"data"`
}

func (this WechatApi) Batchblacklist(openid_list []string) (*utils.CommonError, error) {
	var res utils.CommonError
	ttt := template.New("Batchblacklist")
	ttt.Parse(batchUnblacklistTemp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		OpenIDList []string
	}{
		OpenIDList: openid_list,
	})
	if err := this.DoPost(batchblacklist, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) BatchUnblacklist(openid_list []string) (*utils.CommonError, error) {
	var res utils.CommonError
	ttt := template.New("BatchUnblacklist")
	ttt.Parse(batchUnblacklistTemp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		OpenIDList []string
	}{
		OpenIDList: openid_list,
	})
	if err := this.DoPost(batchunblacklist, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) GetBlacklist(next_openid string) (*Blacklist, error) {
	var res Blacklist
	body := fmt.Sprintf(getBlacklistTemp, next_openid)
	if err := this.DoPost(getblacklist, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

//-----------------------------------模板--------------------------------------------------------------------------------
const (
	get_all_private_template = "https://api.weixin.qq.com/cgi-bin/template/get_all_private_template?access_token=%s"
	get_industry             = "https://api.weixin.qq.com/cgi-bin/template/get_industry?access_token=%s"
	add_template             = "https://api.weixin.qq.com/cgi-bin/template/api_add_template?access_token=%s"
	del_template             = "https://api.weixin.qq.com/cgi-bin/template/del_private_template?access_token=%s"
	send_template            = "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s"
)

const (
	addTemplateTemp = `{"template_id_short":"%s"}`
	delTemplateTemp = `{"template_id":"%s"}`
)

type TemplateSend struct {
	Touser     string `json:"touser"`
	TemplateID string `json:"template_id"`
	Url        string `json:"url"`
	Data       map[string]struct {
		Value string `json:"value"`
		Color string `json:"color,omitempty"`
	} `json:"data"`
}

type TemplateSendRes struct {
	utils.CommonError
	MsgID string `json:"msgid"`
}

func (this WechatApi) SendTemplateMsg(msg *TemplateSend) (*TemplateSendRes, error) {
	var res TemplateSendRes
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	if err := this.DoPost(send_template, string(body), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type TemplateRes struct {
	utils.CommonError
	TemplateId string `json:"template_id"`
}

func (this WechatApi) AddTemplate(template_id_short string) (*TemplateRes, error) {
	var res TemplateRes
	body := fmt.Sprintf(addTemplateTemp, template_id_short)
	if err := this.DoPost(add_template, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) DeleteTemplate(template_id string) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(delTemplateTemp, template_id)
	if err := this.DoPost(del_template, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type TemplateList struct {
	TemplateList []struct {
		TemplateID      string `json:"template_id"`
		Title           string `json:"title"`
		PrimaryIndustry string `json:"primary_industry"`
		DeputyIndustry  string `json:"deputy_industry"`
		Content         string `json:"content"`
		Example         string `json:"example"`
	} `json:"template_list"`
}

func (this WechatApi) GetAllPrivateTemplate() (*TemplateList, error) {
	var res TemplateList
	if err := this.DoGet(get_all_private_template, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type Industry struct {
	PrimaryIndustry struct {
		FirstClass  string `json:"first_class"`
		SecondClass string `json:"second_class"`
	} `json:"primary_industry"`
	SecondaryIndustry struct {
		FirstClass  string `json:"first_class"`
		SecondClass string `json:"second_class"`
	} `json:"secondary_industry"`
}

func (this WechatApi) GetIndustry() (*Industry, error) {
	var res Industry
	if err := this.DoGet(get_industry, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

//-----------------------------------二维码------------------------------------------------------------------------------

const (
	create_qrcode = "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s"
	show_qrcode   = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s"
	shorturl      = "https://api.weixin.qq.com/cgi-bin/shorturl?access_token=%s"
)
const (
	shortUrlTemp = `{
	  "action": "long2short",
	  "long_url": "%s"
	}`
)

type ShortUrl struct {
	utils.CommonError
	ShortUrl string `json:"short_url"`
}

func (this WechatApi) ShortUrl(long_url string) (*ShortUrl, error) {
	var res ShortUrl
	body := fmt.Sprintf(shortUrlTemp, long_url)
	if err := this.DoPost(shorturl, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type qrcodeParam struct {
	ActionName    string `json:"action_name"`
	ExpireSeconds int    `json:"expire_seconds,omitempty"`
	ActionInfo    struct {
		Scene struct {
			SceneStr string `json:"scene_str,omitempty"`
			SceneID  int    `json:"scene_id,omitempty"`
		} `json:"scene"`
	} `json:"action_info"`
}

type QrcodeTicket struct {
	Ticket        string `json:"ticket"`
	Expireseconds int    `json:"expire_seconds"`
	Url           string `json:"url"`
}

func (this WechatApi) CreateLimitStrQrcode(scene_str string) (*QrcodeTicket, error) {
	var res QrcodeTicket
	qrcodeParam := &qrcodeParam{
		ActionName: "QR_LIMIT_STR_SCENE",
	}
	qrcodeParam.ActionInfo.Scene.SceneStr = scene_str
	tmp, _ := json.Marshal(qrcodeParam)
	body := string(tmp)
	if err := this.DoPost(create_qrcode, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) CreateLimitQrcode(scene_id int) (*QrcodeTicket, error) {
	var res QrcodeTicket
	qrcodeParam := &qrcodeParam{
		ActionName: "QR_LIMIT_SCENE",
	}
	qrcodeParam.ActionInfo.Scene.SceneID = scene_id
	tmp, _ := json.Marshal(qrcodeParam)
	body := string(tmp)
	if err := this.DoPost(create_qrcode, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) CreateQrcode(expire_seconds int, scene_id int) (*QrcodeTicket, error) {
	var res QrcodeTicket
	qrcodeParam := &qrcodeParam{
		ActionName:    "QR_SCENE",
		ExpireSeconds: expire_seconds,
	}
	qrcodeParam.ActionInfo.Scene.SceneID = scene_id
	tmp, _ := json.Marshal(qrcodeParam)
	body := string(tmp)
	if err := this.DoPost(create_qrcode, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) ShowQrcode(ticket string) string {
	return string(fmt.Sprintf(show_qrcode, ticket))
}

//-------------------------------------菜单------------------------------------------------------------------------------
const (
	menu_get    = "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s"
	menu_create = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s"
	menu_delete = "https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s"
)

type MenuEntry struct {
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Key       string       `json:"key,omitempty"`
	Url       string       `json:"url,omitempty"`
	AppID     string       `json:"appid,omitempty"`
	Pagepath  string       `json:"pagepath,omitempty"`
	MediaID   string       `json:"media_id,omitempty"`
	SubButton []*MenuEntry `json:"sub_button,omitempty"`
}

type Menu struct {
	Menu struct {
		Buttons []*MenuEntry `json:"button"`
	} `json:"menu"`
}

func (this WechatApi) GetMenu() (*Menu, error) {
	var res Menu
	if err := this.DoGet(menu_get, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) CreateMenu(menus []MenuEntry) (*utils.CommonError, error) {
	var res utils.CommonError
	data := struct {
		Button []MenuEntry `json:"button"`
	}{
		Button: menus,
	}
	body, _ := json.Marshal(&data)

	if err := this.DoPostRaw(menu_create, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) DeleteMenu() (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoGet(menu_delete, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

//-------------------------------------杂项------------------------------------------------------------------------------

const (
	getcallbackip = "https://api.weixin.qq.com/cgi-bin/getcallbackip?access_token=%s"
)

type IpList struct {
	IpList []string `json:"ip_list"`
}

func (this WechatApi) GetIpList() (*IpList, error) {
	var res IpList
	if err := this.DoGet(getcallbackip, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) GetAccessToken() (*string, error) {
	accessToken, err := this.Context.GetAccessToken()
	if err == nil {
		return &accessToken, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) GetJsTicket() (*string, error) {
	jsTicket, err := this.Context.GetJsTicket()
	if err == nil {
		return &jsTicket, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) GetCardTicket() (*string, error) {
	jsTicket, err := this.Context.GetCardTicket()
	if err == nil {
		return &jsTicket, nil
	} else {
		return nil, err
	}
}

type SignJsRes struct {
	AppID     string `json:"appid"`
	Timestamp string `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

func (this WechatApi) SignJsTicket(nonceStr, timestamp, url string) (*SignJsRes, error) {
	jsTicket, err := this.Context.GetJsTicket()
	if err != nil {
		return nil, err
	}
	if timestamp == "" {
		bb := utils.GetCurrTs()
		timestamp = strconv.FormatInt(bb, 10)
	}
	sign := utils.WXConfigSign(jsTicket, nonceStr, timestamp, url)
	return &SignJsRes{
		Signature: sign,
		AppID:     this.Context.Config.AppID,
		NonceStr:  nonceStr,
		Timestamp: timestamp,
	}, nil
}

type SignChooseCardParam struct {
	ApiTicket string `json:"api_ticket"`
	ApiID     string `json:"app_id"`
	ShopID    string `json:"location_id"`
	TimesTamp int64  `json:"times_tamp"`
	NonceStr  string `json:"nonce_str"`
	CardID    string `json:"card_id"`
	CardType  string `json:"card_type"`
}

type SignChooseCardResp struct {
	ShopID    string `json:"shop_id"`
	TimesTamp int64  `json:"timetamp"`
	NonceStr  string `json:"nonce_str"`
	CardID    string `json:"card_id"`
	CardType  string `json:"card_type"`
	SignType  string `json:"sign_type"`
	CardSign  string `json:"card_sign"`
}

func (this WechatApi) signCardImp(variable interface{}) (sign string, err error) {
	ss := &utils.SignStructValue{}
	sign, err = ss.Sign(variable, nil)
	return
}

func (this WechatApi) SignChooseCard(shopID, cardID, cardType string) (*SignChooseCardResp, error) {
	jsTicket, err := this.Context.GetCardTicket()
	if err != nil {
		return nil, err
	}
	param := &SignChooseCardParam{
		ApiTicket: jsTicket,
		ApiID:     this.Context.Config.AppID,
		TimesTamp: utils.GetCurrTs(),
		NonceStr:  utils.RandString(30),
		CardID:    cardID,
		CardType:  cardType,
		ShopID:    shopID,
	}

	sign, err := this.signCardImp(param)
	if err != nil {
		return nil, err
	}
	return &SignChooseCardResp{
		ShopID:    param.ShopID,
		TimesTamp: param.TimesTamp,
		NonceStr:  param.NonceStr,
		CardID:    param.CardID,
		CardType:  param.CardType,
		SignType:  "SHA1",
		CardSign:  sign,
	}, nil
}

type SignAddCardParamObj struct {
	CardID              string `json:"card_id"`
	Code                string `json:"code,omitempty" doc:"指定的卡券code码，只能被领一次。自定义code模式的卡券必须填写，非自定义code和预存code模式的卡券不必填写。"`
	OpenID              string `json:"openid,omitempty" doc:"指定领取者的openid，只有该用户能领取。bind_openid字段为true的卡券必须填写，bind_openid字段为false不必填写。"`
	FixedBegintimestamp int64  `json:"fixed_begintimestamp,omitempty" doc:"卡券在第三方系统的实际领取时间，为东八区时间戳（UTC+8,精确到秒）。当卡券的有效期类型为DATE_TYPE_FIX_TERM时专用，标识卡券的实际生效时间，用于解决商户系统内起始时间和领取时间不同步的问题。"`
	OuterStr            string `json:"outer_str"`
}

type SignAddCardParam struct {
	CardList []SignAddCardParamObj `json:"cardList"`
}

type SignAddCardRespObj struct {
	CardID string `json:"cardId"`
	CardExt string `json:"cardExt"`
}

type SignAddCardResp struct {
	CardList []*SignAddCardRespObj `json:"cardList"`
}

func (this WechatApi) SignAddCard(param *SignAddCardParam) (*SignAddCardResp, error) {
	jsTicket, err := this.Context.GetCardTicket()
	if err != nil {
		return nil, err
	}

	resp := &SignAddCardResp{}
	resp.CardList = make([]*SignAddCardRespObj, 0, len(param.CardList))
	for _, value := range param.CardList {
		tmpVar := &struct {
			CardID              string `json:"-"`
			Code                string `json:"code,omitempty"`
			OpenID              string `json:"openid,omitempty"`
			ApiTicket string `json:"-"`
			TimesTamp int64  `json:"timetamp"`
			NonceStr  string `json:"nonce_str"`
			Sign string `json:"signature" sign:"-"`
			FixedBegintimestamp int64  `json:"fixed_begintimestamp,omitempty" sign:"-"`
			OuterStr            string `json:"outer_str" sign:"-"`
		}{
			CardID: value.CardID,
			Code: value.Code,
			OpenID: value.OpenID,
			FixedBegintimestamp: value.FixedBegintimestamp,
			OuterStr: value.OuterStr,
			ApiTicket: jsTicket,
			TimesTamp: utils.GetCurrTs(),
			NonceStr:  utils.RandString(30),
		}


		sign, err := this.signCardImp(tmpVar)
		if err != nil {
			return nil, err
		}
		tmpVar.Sign = sign

		cardExt,_ := json.Marshal(tmpVar)
		resp.CardList = append(resp.CardList,&SignAddCardRespObj{
			CardID: value.CardID,
			CardExt: string(cardExt),
		})
	}

	return resp,nil
}