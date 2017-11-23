package corp

import (
	"log"
	"encoding/xml"
	"fmt"
	"github.com/qjw/go-wx-sdk/utils"
	"net/http"
	"strconv"
)

type KfMessageHandle func(*KfMixMessage)

type KfServer struct {
	Request     *http.Request
	Responce    http.ResponseWriter
	CorpContext *Context
	KfConfig    *KfConfig

	// 收到消息的回调
	MessageHandler KfMessageHandle
}

type KfServerRequest struct {
	MixedMsg        *KfMixMessage
	RequestHttpBody *KfRequestEncryptedXMLMsg

	Random    []byte
	Nonce     string
	Timestamp int64

	// 回复的消息
	ResponseRawXMLMsg []byte
	ResponseMsg       utils.Reply
}

//NewServer init
func NewKfServer(request *http.Request, responce http.ResponseWriter,
	handle KfMessageHandle, mpwcontext *Context, kfConfig *KfConfig) *KfServer {
	return &KfServer{
		Request:        request,
		Responce:       responce,
		CorpContext:    mpwcontext,
		KfConfig:       kfConfig,
		MessageHandler: handle,
	}
}

//Serve 处理微信的请求消息
func (srv *KfServer) Ping() {
	echostr := srv.Request.URL.Query().Get("echostr")
	if echostr == "" {
		log.Print("invalid echostr")
		http.Error(srv.Responce, "", http.StatusForbidden)
		return
	}

	if !srv.validate(nil, echostr) {
		http.Error(srv.Responce, "", http.StatusForbidden)
		return
	}

	_, echostrRes, _, err := utils.DecryptMsg(srv.CorpContext.Config.CorpID, echostr, srv.KfConfig.EncodingAESKey)
	if err != nil {
		log.Print("invalid DecryptMsg")
		http.Error(srv.Responce, "", http.StatusForbidden)
	}

	http.Error(srv.Responce, string(echostrRes), http.StatusOK)
}

//Serve 处理微信的请求消息
func (srv *KfServer) Serve() error {
	var svrReq KfServerRequest

	// 解析 RequestHttpBody
	var requestHttpBody KfRequestEncryptedXMLMsg
	if err := xml.NewDecoder(srv.Request.Body).Decode(&requestHttpBody); err != nil {
		log.Print(err.Error())
		return err
	}

	if !srv.validate(&svrReq, requestHttpBody.EncryptedMsg) {
		return fmt.Errorf("请求校验失败")
	}

	svrReq.RequestHttpBody = &requestHttpBody
	err := srv.handleRequest(&svrReq)
	if err != nil {
		return err
	}

	// 企业在收到数据包时，需回复XML里的PackageId节点值，表示成功接收，否则企业号侧认为回调失败。
	http.Error(srv.Responce, strconv.FormatInt(svrReq.MixedMsg.PackageId, 10), http.StatusOK)
	return nil
}

func (srv *KfServer) validate(svrReq *KfServerRequest, content string) bool {
	signature := srv.Request.URL.Query().Get("msg_signature")
	if signature == "" {
		log.Print("invalid msg_signature")
		return false
	}
	timestamp := srv.Request.URL.Query().Get("timestamp")
	if timestamp == "" {
		log.Print("invalid timestamp")
		return false
	}

	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Print(err.Error())
		return false
	}

	nonce := srv.Request.URL.Query().Get("nonce")
	if nonce == "" {
		log.Print("invalid nonce")
		return false
	}

	// 验证签名
	msgSignature2 := utils.Signature(srv.KfConfig.Token, timestamp, nonce, content)
	if signature != msgSignature2 {
		log.Print("invalid signature")
		return false
	}

	if svrReq != nil {
		svrReq.Timestamp = timestampInt
		svrReq.Nonce = nonce
	}
	return true
}

//HandleRequest 处理微信的请求
func (srv *KfServer) handleRequest(svrReq *KfServerRequest) (err error) {
	err = srv.getMessage(svrReq)
	if err != nil {
		return
	}

	if srv.MessageHandler != nil {
		srv.MessageHandler(svrReq.MixedMsg)
	}
	return
}

//getMessage 解析微信返回的消息
func (srv *KfServer) getMessage(svrReq *KfServerRequest) error {
	// 解密
	random, rawMsgXML, appID, err := utils.DecryptMsg(
		srv.CorpContext.Config.CorpID,
		svrReq.RequestHttpBody.EncryptedMsg,
		srv.KfConfig.EncodingAESKey)
	if err != nil {
		log.Print("invalid DecryptMsg")
		return err
	}
	svrReq.Random = random

	if svrReq.RequestHttpBody.ToUserName != appID {
		err := fmt.Errorf("the RequestHttpBody's ToUserName(==%s) mismatch the appID with aes encrypt(==%s)",
			svrReq.RequestHttpBody.ToUserName, appID)
		return err
	}

	// 解密成功, 解析 MixedMessage
	var mixedMsg KfMixMessage
	if err = xml.Unmarshal(rawMsgXML, &mixedMsg); err != nil {
		log.Print(err.Error())
		return err
	}
	svrReq.MixedMsg = &mixedMsg

	// 安全考虑再次验证
	if svrReq.RequestHttpBody.ToUserName != mixedMsg.ToUserName {
		err := fmt.Errorf("the RequestHttpBody's ToUserName(==%s) mismatch the MixedMessage's SuiteId",
			svrReq.RequestHttpBody.ToUserName)
		return err
	}

	return nil
}

// 微信服务器请求 http body
type KfRequestEncryptedXMLMsg struct {
	utils.RequestEncryptedXMLMsg
	AgentType string `xml:"AgentType"`
}

// CommonToken 消息中通用的结构
type KfMixMessage struct {
	XMLName    xml.Name           `xml:"xml" json:"-"`
	ToUserName string             `xml:"ToUserName" json:"to_username"`
	AgentType  string             `xml:"AgentType" json:"agent_type"`
	ItemCount  int                `xml:"ItemCount" json:"item_count"`
	PackageId  int64              `xml:"PackageId" json:"package_id"`
	Item       []KfMixMessageItem `xml:"Item" json:"items"`
}

type KfMixMessageItemHead struct {
	FromUserName string `xml:"FromUserName" json:"from_username"`
	CreateTime   int64  `xml:"CreateTime" json:"create_time"`
	MsgType      string `xml:"MsgType" json:"msg_type"`
}

type KfMixMessageItem struct {
	KfMixMessageItemHead
	MsgID int64 `xml:"MsgId" json:"msg_id"`

	Content string `xml:"Content" json:"content,omitempty"`
	PicURL  string `xml:"PicUrl" json:"pic_url,omitempty"`
	MediaID string `xml:"MediaId" json:"media_id,omitempty"`
	//Format       string `xml:"Format" json:"format,omitempty"`
	//ThumbMediaID string `xml:"ThumbMediaId" json:"thumb_media_id,omitempty"`

	// link
	Title       string `xml:"Title" json:"title,omitempty"`
	Description string `xml:"Description" json:"description,omitempty"`
	URL         string `xml:"Url" json:"url,omitempty"`
	// location
	LocationX float64 `xml:"Location_X" json:"location_x,omitempty"`
	LocationY float64 `xml:"Location_Y" json:"location_y,omitempty"`
	Scale     float64 `xml:"Scale" json:"scale,omitempty"`
	Label     string  `xml:"Label" json:"label,omitempty"`
}
