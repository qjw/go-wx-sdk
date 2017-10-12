package corp

import (
	"errors"
	"github.com/qjw/go-wx-sdk/utils"
	"strconv"
	"strings"
)

const (
	msssageSend = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
)

type MsgSendRes struct {
	utils.CommonError
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

type textMsgObj struct {
	Content string `json:"content"`
}

type imageMsgObj struct {
	MediaID string `json:"media_id"`
}
type voiceMsgObj imageMsgObj
type fileMsgObj imageMsgObj

type videoMsgObj struct {
	MediaID     string `json:"media_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type NewsMsgObj struct {
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Picurl      string `json:"picurl"`
	} `json:"articles"`
}

type TextCardObj struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}

type MsgObj struct {
	ToUser   string       `json:"touser,omitempty" doc:"成员ID列表（消息接收者，多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为@all，则向关注该企业应用的全部成员发送"`
	ToParty  string       `json:"toparty,omitempty" doc:"部门ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数"`
	ToTag    string       `json:"totag,omitempty" doc:"标签ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数"`
	MsgType  string       `json:"msgtype"`
	AgentID  int64        `json:"agentid"`
	Safe     string       `json:"safe,omitempty" doc:"表示是否是保密消息，0表示否，1表示是，默认0"`
	Text     *textMsgObj  `json:"text,omitempty"`
	Image    *imageMsgObj `json:"image,omitempty"`
	Voice    *voiceMsgObj `json:"voice,omitempty"`
	File     *fileMsgObj  `json:"file,omitempty"`
	Video    *videoMsgObj `json:"video,omitempty"`
	TextCard *TextCardObj `json:"textcard,omitempty"`
	News     *NewsMsgObj  `json:"news,omitempty"`
}

// 用 '|' 连接 a 的各个元素的十进制字符串
func joinInt64(a []int64, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return strconv.FormatInt(a[0], 10)
	default:
		strs := make([]string, len(a))
		for i, n := range a {
			strs[i] = strconv.FormatInt(n, 10)
		}
		return strings.Join(strs, sep)
	}
}

func (this CorpApi) sendMsgImp(touser []string, toparty, totag []int64, msg *MsgObj) (*utils.CommonError, error) {
	if len(touser) == 0 && len(toparty) == 0 && len(totag) == 0 {
		return nil, errors.New("empty reciver")
	}
	if len(touser) > 0 {
		msg.ToUser = strings.Join(touser, "|")
	}
	if len(toparty) > 0 {
		msg.ToParty = joinInt64(toparty, "|")
	}
	if len(totag) > 0 {
		msg.ToTag = joinInt64(totag, "|")
	}

	var res utils.CommonError
	if err := this.DoPostObject(msssageSend, msg, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) SendTextCardMsg(touser []string, toparty, totag []int64, agentid int64,
	card *TextCardObj) (*utils.CommonError, error) {
	return this.sendMsgImp(touser, toparty, totag, &MsgObj{
		MsgType: "textcard",
		AgentID: agentid,
		TextCard: card,
	})
}

func (this CorpApi) SendTextMsg(touser []string, toparty, totag []int64, agentid int64,
	content string) (*utils.CommonError, error) {
	return this.sendMsgImp(touser, toparty, totag, &MsgObj{
		MsgType: "text",
		AgentID: agentid,
		Text: &textMsgObj{
			Content: content,
		},
	})
}

func (this CorpApi) SendImageMsg(touser []string, toparty, totag []int64, agentid int64,
	mediaid string) (*utils.CommonError, error) {
	return this.sendMsgImp(touser, toparty, totag, &MsgObj{
		MsgType: "image",
		AgentID: agentid,
		Image: &imageMsgObj{
			MediaID: mediaid,
		},
	})
}

func (this CorpApi) SendVoiceMsg(touser []string, toparty, totag []int64, agentid int64,
	mediaid string) (*utils.CommonError, error) {
	return this.sendMsgImp(touser, toparty, totag, &MsgObj{
		MsgType: "voice",
		AgentID: agentid,
		Voice: &voiceMsgObj{
			MediaID: mediaid,
		},
	})
}

func (this CorpApi) SendFileMsg(touser []string, toparty, totag []int64, agentid int64,
	mediaid string) (*utils.CommonError, error) {
	return this.sendMsgImp(touser, toparty, totag, &MsgObj{
		MsgType: "file",
		AgentID: agentid,
		File: &fileMsgObj{
			MediaID: mediaid,
		},
	})
}

func (this CorpApi) SendVideoMsg(touser []string, toparty, totag []int64, agentid int64,
	mediaid, title, desc string) (*utils.CommonError, error) {
	return this.sendMsgImp(touser, toparty, totag, &MsgObj{
		MsgType: "video",
		AgentID: agentid,
		Video: &videoMsgObj{
			Title:       title,
			Description: desc,
			MediaID:     mediaid,
		},
	})
}

func (this CorpApi) SendNewsMsg(touser []string, toparty, totag []int64, agentid int64,
	news *NewsMsgObj) (*utils.CommonError, error) {
	return this.sendMsgImp(touser, toparty, totag, &MsgObj{
		MsgType: "news",
		AgentID: agentid,
		News:    news,
	})
}
