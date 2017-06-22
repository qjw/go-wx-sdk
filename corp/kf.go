package corp

import (
	"github.com/qjw/go-wx-sdk/utils"
)

const (
	kfSend = "https://qyapi.weixin.qq.com/cgi-bin/kf/send?access_token=%s"
	kfList = "https://qyapi.weixin.qq.com/cgi-bin/kf/list?access_token=%s&type=%s"
)

type KfMsgUserObj struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type KfTextMsgObj struct {
	Content string `json:"content"`
}

type KfImageMsgObj struct {
	MediaID string `json:"media_id"`
}

type KfMsgObj struct {
	Sender   KfMsgUserObj   `json:"sender"`
	Receiver KfMsgUserObj   `json:"receiver"`
	Msgtype  string         `json:"msgtype"`
	Text     *KfTextMsgObj  `json:"text,omitempty"`
	Image    *KfImageMsgObj `json:"image,omitempty"`
	File     *KfImageMsgObj `json:"file,omitempty"`
	Voice    *KfImageMsgObj `json:"voice,omitempty"`
}

func (this CorpApi) sendKfImp(from, to *KfMsgUserObj, msg *KfMsgObj) (*utils.CommonError, error) {
	var res utils.CommonError
	if err := this.DoPostObject(kfSend, msg, &res); err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}

func (this CorpApi) SendKfText(from, to *KfMsgUserObj, content string) (*utils.CommonError, error) {
	return this.sendKfImp(from, to, &KfMsgObj{
		Sender:   *from,
		Receiver: *to,
		Msgtype:  "text",
		Text: &KfTextMsgObj{
			Content: content,
		},
	})
}

func (this CorpApi) SendKfImage(from, to *KfMsgUserObj, mediaid string) (*utils.CommonError, error) {
	return this.sendKfImp(from, to, &KfMsgObj{
		Sender:   *from,
		Receiver: *to,
		Msgtype:  "image",
		Image: &KfImageMsgObj{
			MediaID: mediaid,
		},
	})
}

func (this CorpApi) SendKfFile(from, to *KfMsgUserObj, mediaid string) (*utils.CommonError, error) {
	return this.sendKfImp(from, to, &KfMsgObj{
		Sender:   *from,
		Receiver: *to,
		Msgtype:  "file",
		File: &KfImageMsgObj{
			MediaID: mediaid,
		},
	})
}

func (this CorpApi) SendKfVoice(from, to *KfMsgUserObj, mediaid string) (*utils.CommonError, error) {
	return this.sendKfImp(from, to, &KfMsgObj{
		Sender:   *from,
		Receiver: *to,
		Msgtype:  "voice",
		Voice: &KfImageMsgObj{
			MediaID: mediaid,
		},
	})
}

type KfObj struct {
	Users   []string `json:"user"`
	Parties []int64 `json:"party"`
	Tags    []string `json:"tag"`
}

type KfListRes struct {
	utils.CommonError
	Internal *KfObj `json:"internal,omitempty"`
	External *KfObj `json:"external,omitempty"`
}

func (this CorpApi) GetKfList(tp string) (*KfListRes, error) {
	var res KfListRes
	err := this.DoGet(kfList, &res, tp)
	if err == nil {
		return &res, nil
	} else {
		return nil, err
	}
}