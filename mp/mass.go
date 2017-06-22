package mp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"bytes"
	"fmt"
	"html/template"
)

//-----------------------------------群发--------------------------------------------------------------------------------
const (
	mass_preview = "https://api.weixin.qq.com/cgi-bin/message/mass/preview?access_token=%s"
	mass_send    = "https://api.weixin.qq.com/cgi-bin/message/mass/send?access_token=%s"
	mass_get     = "https://api.weixin.qq.com/cgi-bin/message/mass/get?access_token=%s"
	mass_delete  = "https://api.weixin.qq.com/cgi-bin/message/mass/delete?access_token=%s"
	mass_sendall = "https://api.weixin.qq.com/cgi-bin/message/mass/sendall?access_token=%s"
)

const (
	mass_preview_temp = `{
	   "touser":"%s",
	   "%s":{
		    "media_id":"%s"
		     },
	   "msgtype":"%s"
	}`
	mass_preview_txt_temp = `{
	    "touser":"%s",
	    "text":{
		   "content":"%s"
		   },
	    "msgtype":"text"
	}`
	mass_send_mpnews_temp = `{
	   "touser":[
	   		{{range $index, $element := .OpenIDList}}
	   		{{if $index}},{{end}}
				"{{$element}}"
			{{end}}
	   ],
	   "mpnews":{
	      "media_id":"{{.MediaID}}"
	   },
	   "msgtype":"mpnews"，
	   "send_ignore_reprint":{{.SendIgnoreReprint}}
	}`
	mass_send_txt_temp = `{
	   "touser":[
	   		{{range $index, $element := .OpenIDList}}
	   		{{if $index}},{{end}}
				"{{$element}}"
			{{end}}
	   ],
	   "msgtype": "text",
	   "text": { "content": "{{.Content}}"}
	}`
	mass_send_card_temp = `{
	   "touser":[
	   		{{range $index, $element := .OpenIDList}}
	   		{{if $index}},{{end}}
				"{{$element}}"
			{{end}}
	   ],
	   "msgtype": "wxcard",
	   "wxcard": { "card_id": "{{.CardID}}"}
	}`
	mass_send_msg_temp = `{
	   "touser":[
	   		{{range $index, $element := .OpenIDList}}
	   		{{if $index}},{{end}}
				"{{$element}}"
			{{end}}
	   ],
	   "{{.Tp}}":{
	      "media_id":"{{.MediaID}}"
	   },
	   "msgtype":"{{.Tp}}"
	}`
	mass_get_temp = `{
	   "msg_id": %d
	}`

	mass_sendall_mpnews_temp = `{
	   "filter":{
	      "is_to_all":%t,
	      "tag_id":%d
	   },
	   "mpnews":{
	      "media_id":"%s"
	   },
	   "msgtype":"mpnews",
	   "send_ignore_reprint":%d
	}`
	mass_sendall_text = `{
	   "filter":{
	      "is_to_all":%t,
	      "tag_id":%d
	   },
	   "text":{
	      "content":"%s"
	   },
	   "msgtype":"text"
	}`
	mass_sendall_msg = `{
	   "filter":{
	      "is_to_all":%t,
	      "tag_id":%d
	   },
	   "%s":{
	      "media_id":"%s"
	   },
	   "msgtype":"%s"
	}`
	mass_sendall_card = `{
	   "filter":{
	      "is_to_all":%t,
	      "tag_id":%d
	   },
	   "wxcard":{
               "card_id":"%s"
           },
   	   "msgtype":"wxcard"

	}`
)

type MassPreviewRes struct {
	utils.CommonError
	MsgID string `json:"msg_id"`
}

func (this WechatApi) MassPreviewMsg(touser, media_id, tp string) (*MassPreviewRes, error) {
	var res MassPreviewRes
	body := fmt.Sprintf(mass_preview_temp, touser, tp, media_id, tp)
	if err := this.DoPost(mass_preview, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) MassPreviewText(touser, content string) (*MassPreviewRes, error) {
	var res MassPreviewRes
	body := fmt.Sprintf(mass_preview_txt_temp, touser, content)
	if err := this.DoPost(mass_preview, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type MassSendRes struct {
	utils.CommonError
	MsgID     int64 `json:"msg_id"`
	MsgDataID int64 `json:"msg_data_id"`
}

func (this WechatApi) MassSendMpnews(tousers []string, media_id string, send_ignore_reprint int) (*MassSendRes, error) {
	var res MassSendRes
	ttt := template.New("MassSendMpnews")
	ttt.Parse(mass_send_mpnews_temp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		OpenIDList        []string
		MediaID           string
		SendIgnoreReprint int
	}{
		MediaID:           media_id,
		OpenIDList:        tousers,
		SendIgnoreReprint: send_ignore_reprint,
	})

	if err := this.DoPost(mass_send, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) MassSendText(tousers []string, content string) (*MassSendRes, error) {
	var res MassSendRes
	ttt := template.New("MassSendText")
	ttt.Parse(mass_send_txt_temp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		OpenIDList []string
		Content    string
	}{
		Content:    content,
		OpenIDList: tousers,
	})

	if err := this.DoPost(mass_send, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) MassSendCard(tousers []string, card_id string) (*MassSendRes, error) {
	var res MassSendRes
	ttt := template.New("MassSendText")
	ttt.Parse(mass_send_card_temp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		OpenIDList []string
		CardID     string
	}{
		CardID:     card_id,
		OpenIDList: tousers,
	})

	if err := this.DoPost(mass_send, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

// voice/image
func (this WechatApi) MassSendMsg(tousers []string, media_id string, tp string) (*MassSendRes, error) {
	var res MassSendRes
	ttt := template.New("MassSendMsg")
	ttt.Parse(mass_send_msg_temp)
	var buf bytes.Buffer
	ttt.Execute(&buf, struct {
		OpenIDList []string
		MediaID    string
		Tp         string
	}{
		MediaID:    media_id,
		OpenIDList: tousers,
		Tp:         tp,
	})

	if err := this.DoPost(mass_send, buf.String(), &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

type MassGetResult struct {
	utils.CommonError
	MsgID     int64  `json:"msg_id"`
	MsgStatus string `json:"msg_status"`
}

func (this WechatApi) MassGet(msg_id int64) (*MassGetResult, error) {
	var res MassGetResult
	body := fmt.Sprintf(mass_get_temp, msg_id)
	if err := this.DoPost(mass_get, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) MassDelete(msg_id int64) (*utils.CommonError, error) {
	var res utils.CommonError
	body := fmt.Sprintf(mass_get_temp, msg_id)
	if err := this.DoPost(mass_delete, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) MassSendAllMpnews(media_id string, tag_id int,
	is_to_all bool, send_ignore_reprint int) (*MassSendRes, error) {
	var res MassSendRes
	body := fmt.Sprintf(mass_sendall_mpnews_temp, is_to_all, tag_id, media_id, send_ignore_reprint)
	if err := this.DoPost(mass_sendall, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) MassSendAllText(content string, tag_id int, is_to_all bool) (*MassSendRes, error) {
	var res MassSendRes
	body := fmt.Sprintf(mass_sendall_text, is_to_all, tag_id, content)
	if err := this.DoPost(mass_sendall, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) MassSendAllCard(card_id string, tag_id int, is_to_all bool) (*MassSendRes, error) {
	var res MassSendRes
	body := fmt.Sprintf(mass_sendall_card, is_to_all, tag_id, card_id)
	if err := this.DoPost(mass_sendall, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this WechatApi) MassSendAllMsg(media_id string, tag_id int, is_to_all bool, tp string) (*MassSendRes, error) {
	var res MassSendRes
	body := fmt.Sprintf(mass_sendall_msg, is_to_all, tag_id, tp, media_id, tp)
	if err := this.DoPost(mass_sendall, body, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}
