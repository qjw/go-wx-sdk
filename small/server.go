package small

import (
	"encoding/xml"
	"fmt"
	"log"
	"github.com/qjw/go-wx-sdk/utils"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
)

type MessageHandle func(*MixMessage) utils.Reply

type Server struct {
	Request  *http.Request
	Responce http.ResponseWriter
	SContext *Context

	// 公众号的OpenID
	// openID string

	// 收到消息的回调
	MessageHandler MessageHandle
}

type ServerRequest struct {
	MixedMsg *MixMessage
	// 加密模式下才有
	RequestHttpBody *utils.RequestEncryptedXMLMsg
	// 收到的原始数据
	RequestRawXMLMsg []byte

	// 安全（加密）模式
	IsSafeMode bool
	Random     []byte
	Nonce      string
	Timestamp  int64

	// 回复的消息
	ResponseMsg utils.Reply
}

//NewServer init
func NewServer(request *http.Request, responce http.ResponseWriter,
	handle MessageHandle, mpwcontext *Context) *Server {
	return &Server{
		Request:        request,
		Responce:       responce,
		MessageHandler: handle,
		SContext:       mpwcontext,
	}
}

//Serve 处理微信的请求消息
func (srv *Server) Ping() {
	if !srv.validate(nil) {
		http.Error(srv.Responce, "", http.StatusForbidden)
		return
	}

	echostr := srv.Request.URL.Query().Get("echostr")
	if echostr == "" {
		http.Error(srv.Responce, "", http.StatusForbidden)
		return
	}
	http.Error(srv.Responce, echostr, http.StatusOK)
}

//Serve 处理微信的请求消息
func (srv *Server) Serve() error {
	var svrReq ServerRequest
	if !srv.validate(&svrReq) {
		return fmt.Errorf("请求校验失败")
	}

	err := srv.handleRequest(&svrReq)
	if err != nil {
		return err
	}

	if err = srv.buildResponse(&svrReq); err != nil {
		return err
	}
	if err = srv.send(&svrReq); err != nil {
		return err
	}
	return nil
}

//Validate 校验请求是否合法
func (srv *Server) validate(svrReq *ServerRequest) bool {
	signature := srv.Request.URL.Query().Get("signature")
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

	if signature == utils.Signature(srv.SContext.Config.Token, timestamp, nonce) {
		if svrReq != nil {
			svrReq.Timestamp = timestampInt
			svrReq.Nonce = nonce
		}
		return true
	} else {
		return false
	}
}

//HandleRequest 处理微信的请求
func (srv *Server) handleRequest(svrReq *ServerRequest) (err error) {
	//set isSafeMode
	svrReq.IsSafeMode = false
	encryptType := srv.Request.URL.Query().Get("encrypt_type")
	if encryptType == "aes" {
		svrReq.IsSafeMode = true
	}

	//set openID
	// srv.openID = srv.Context.Query("openid")
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
	var rawXMLMsgBytes []byte
	var err error
	if svrReq.IsSafeMode {
		var encryptedXMLMsg utils.RequestEncryptedXMLMsg
		if err := xml.NewDecoder(srv.Request.Body).Decode(&encryptedXMLMsg); err != nil {
			return fmt.Errorf("从body中解析xml失败,err=%v", err)
		}
		svrReq.RequestHttpBody = &encryptedXMLMsg

		//解密
		svrReq.Random, rawXMLMsgBytes, _, err = utils.DecryptMsg(srv.SContext.Config.AppID,
			encryptedXMLMsg.EncryptedMsg,
			srv.SContext.Config.EncodingAESKey)
		if err != nil {
			return fmt.Errorf("消息解密失败, err=%v", err)
		}
	} else {
		rawXMLMsgBytes, err = ioutil.ReadAll(srv.Request.Body)
		if err != nil {
			return fmt.Errorf("从body中解析xml失败, err=%v", err)
		}
	}

	var msg MixMessage
	if err := xml.Unmarshal(rawXMLMsgBytes, &msg); err != nil {
		return err
	}

	svrReq.RequestRawXMLMsg = rawXMLMsgBytes
	svrReq.MixedMsg = &msg
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
	return
}

//Send 将自定义的消息发送
func (srv *Server) send(svrReq *ServerRequest) (err error) {
	if svrReq.ResponseMsg == nil {
		return
	}

	var replyMsg interface{} = svrReq.ResponseMsg
	if svrReq.IsSafeMode {
		responseRawXMLMsg, err := xml.Marshal(svrReq.ResponseMsg)
		if err != nil {
			return err
		}

		//安全模式下对消息进行加密
		var encryptedMsg []byte
		encryptedMsg, err = utils.EncryptMsg(svrReq.Random, responseRawXMLMsg,
			srv.SContext.Config.AppID,
			srv.SContext.Config.EncodingAESKey)
		if err != nil {
			return err
		}
		// 如果获取不到timestamp nonce 则自己生成
		timestamp := svrReq.Timestamp
		timestampStr := strconv.FormatInt(timestamp, 10)
		msgSignature := utils.Signature(srv.SContext.Config.Token, timestampStr,
			svrReq.Nonce, string(encryptedMsg))
		replyMsg = utils.ResponseEncryptedXMLMsg{
			EncryptedMsg: string(encryptedMsg),
			MsgSignature: msgSignature,
			Timestamp:    timestamp,
			Nonce:        svrReq.Nonce,
		}
	}
	if replyMsg != nil {
		data, _ := xml.MarshalIndent(replyMsg, "", "\t")

		srv.Responce.Header().Set("Content-Type", "application/xml; charset=utf-8")
		srv.Responce.WriteHeader(http.StatusOK)
		srv.Responce.Write(data)
	}
	return nil
}
