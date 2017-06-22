package corp

import (
	"github.com/qjw/go-wx-sdk/utils"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
)

type MessageHandle func(*MixMessage) utils.Reply

type Server struct {
	Request     *http.Request
	Responce    http.ResponseWriter

	CorpContext *Context
	AgentConfig *AgentConfig

	// 收到消息的回调
	MessageHandler MessageHandle
}

type ServerRequest struct {
	MixedMsg        *MixMessage
	RequestHttpBody *utils.RequestEncryptedXMLMsg

	Random    []byte
	Nonce     string
	Timestamp int64

	// 回复的消息
	ResponseRawXMLMsg []byte
	ResponseMsg       utils.Reply
}

//NewServer init
func NewServer(request *http.Request, responce http.ResponseWriter,
	handle MessageHandle, mpwcontext *Context, agentConfig *AgentConfig) *Server {
	return &Server{
		Request:        request,
		Responce:       responce,
		CorpContext:    mpwcontext,
		AgentConfig:    agentConfig,
		MessageHandler: handle,
	}
}

//Serve 处理微信的请求消息
func (srv *Server) Ping() {
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

	_, echostrRes, _, err := utils.DecryptMsg(srv.CorpContext.Config.CorpID, echostr, srv.AgentConfig.EncodingAESKey)
	if err != nil {
		log.Print("invalid DecryptMsg")
		http.Error(srv.Responce, "", http.StatusForbidden)
	}

	http.Error(srv.Responce, string(echostrRes), http.StatusOK)
}

//Serve 处理微信的请求消息
func (srv *Server) Serve() error {
	var svrReq ServerRequest

	// 解析 RequestHttpBody
	var requestHttpBody utils.RequestEncryptedXMLMsg
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

	if err = srv.buildResponse(&svrReq);err != nil{
		return err
	}
	if err = srv.send(&svrReq);err != nil{
		return err
	}
	return nil
}

func (srv *Server) validate(svrReq *ServerRequest, content string) bool {
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
	msgSignature2 := utils.Signature(srv.AgentConfig.Token, timestamp, nonce, content)
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
func (srv *Server) handleRequest(svrReq *ServerRequest) (err error) {
	err = srv.getMessage(svrReq)
	if err != nil {
		return
	}

	if srv.MessageHandler != nil {
		svrReq.ResponseMsg = srv.MessageHandler(svrReq.MixedMsg)
	}
	return
}

//getMessage 解析微信返回的消息
func (srv *Server) getMessage(svrReq *ServerRequest) error {
	// 解密
	random, rawMsgXML, appID, err := utils.DecryptMsg(
		srv.CorpContext.Config.CorpID,
		svrReq.RequestHttpBody.EncryptedMsg,
		srv.AgentConfig.EncodingAESKey)
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
	var mixedMsg MixMessage
	if err = xml.Unmarshal(rawMsgXML, &mixedMsg); err != nil {
		log.Print(err.Error())
		return err
	}
	svrReq.MixedMsg = &mixedMsg

	js, _ := json.MarshalIndent(mixedMsg, "", " ")
	log.Print(string(js))

	// 安全考虑再次验证
	if svrReq.RequestHttpBody.ToUserName != mixedMsg.ToUserName {
		err := fmt.Errorf("the RequestHttpBody's ToUserName(==%s) mismatch the MixedMessage's SuiteId",
			svrReq.RequestHttpBody.ToUserName)
		return err
	}

	return nil
}

func (srv *Server) buildResponse(svrReq *ServerRequest) (err error) {
	reply := svrReq.ResponseMsg
	if reply == nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic error: %v\n%s", e, debug.Stack())
		}
	}()
	msgType := reply.GetMsgType()
	switch msgType {
	case utils.MsgTypeText:
	case utils.MsgTypeImage:
	case utils.MsgTypeVoice:
	case utils.MsgTypeVideo:
	case utils.MsgTypeMusic:
	case utils.MsgTypeNews:
	case utils.MsgTypeTransfer:
	default:
		err = utils.ErrUnsupportReply
		return
	}

	reply.SetToUserName(svrReq.MixedMsg.FromUserName)
	reply.SetFromUserName(svrReq.MixedMsg.ToUserName)
	reply.SetCreateTime(utils.GetCurrTs())
	svrReq.ResponseRawXMLMsg, err = xml.Marshal(svrReq.ResponseMsg)
	return
}

//Send 将自定义的消息发送
func (srv *Server) send(svrReq *ServerRequest) (err error) {
	if svrReq.ResponseMsg == nil || svrReq.ResponseRawXMLMsg == nil {
		return
	}
	//安全模式下对消息进行加密
	var encryptedMsg []byte
	encryptedMsg, err = utils.EncryptMsg(svrReq.Random, svrReq.ResponseRawXMLMsg,
		srv.CorpContext.Config.CorpID,
		srv.AgentConfig.EncodingAESKey)
	if err != nil {
		return
	}
	// 如果获取不到timestamp nonce 则自己生成
	timestamp := svrReq.Timestamp
	timestampStr := strconv.FormatInt(timestamp, 10)
	msgSignature := utils.Signature(srv.AgentConfig.Token, timestampStr, svrReq.Nonce, string(encryptedMsg))
	replyMsg := utils.ResponseEncryptedXMLMsg{
		EncryptedMsg: string(encryptedMsg),
		MsgSignature: msgSignature,
		Timestamp:    timestamp,
		Nonce:        svrReq.Nonce,
	}

	data, _ := xml.MarshalIndent(replyMsg, "", "\t")
	log.Print(string(data))

	srv.Responce.Header().Set("Content-Type", "application/xml; charset=utf-8")
	srv.Responce.WriteHeader(http.StatusOK)
	srv.Responce.Write(data)
	return
}
