package mp

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/juju/errgo/errors"
	"github.com/ilisin/weixin/util"
)

const MESSAGE_TIMEOUT = 4 * time.Second

type WxMp struct {
	config  *WxMpConfig
	token   *WeiXinToken
	jsToken *JsapiTicket
	Logger  *logrus.Logger

	HttpClient *util.HttpClient

	//message handler
	TextMessageHandler       MessageHandler
	ImageMessageHandler      MessageHandler
	VoiceMessageHandler      MessageHandler
	VideoMessageHandler      MessageHandler
	ShortVideoMessageHandler MessageHandler
	LocationMessageHandler   MessageHandler
	LinkMessageHandler       MessageHandler

	//event handler
	SubscribeEventHandler   MessageHandler
	UnSubscribeEventHandler MessageHandler
	ScanEventHandler        MessageHandler
	LocationEventHandler    MessageHandler
	ClickEventHandler       MessageHandler //menu click
	ViewEventHandler        MessageHandler
}

func NewWxMp(config *WxMpConfig) *WxMp {
	mp := &WxMp{
		config: config,
	}
	mp.Logger = logrus.New()
	l, err := logrus.ParseLevel(config.LoggerLevel)
	if err != nil {
		logrus.Fatal("unkown weixin mp logger level")
	}
	mp.Logger.Level = l
	mp.HttpClient = util.NewClient(config.ServiceURL, mp.Logger.Level)
	return mp
}

func (this *WxMp) HttpServe(request *http.Request, responseWriter http.ResponseWriter) error {
	if strings.ToUpper(request.Method) == "GET" {
		return this.messageSignature(request, responseWriter)
	} else {
		return this.messageFilter(request, responseWriter)
	}
}

func (this *WxMp) messageFilter(request *http.Request, responseWriter http.ResponseWriter) error {
	if request == nil {
		return fmt.Errorf("request paramer error")
	}
	defer request.Body.Close()
	message := &MixMessage{}
	err := xml.NewDecoder(request.Body).Decode(message)
	if err != nil {
		this.Logger.WithField("error", err).Errorf("unkown weixin message request")
	}
	this.Logger.WithField("message", message).Debug("messageFilter")
	rc := make(chan *MixMessage)
	c := time.NewTicker(MESSAGE_TIMEOUT)
	go func() {
		repMessage := this.messageRoute(message)
		rc <- repMessage
	}()
	// if time out ,sync response null
	select {
	case <-c.C:
		go this.messageAsync(rc)
	case msg := <-rc:
		if msg != nil {
			repDat, err := xml.Marshal(msg)
			this.Logger.WithField("message", string(repDat)).Debug("response message")
			buffer := bytes.NewBuffer(repDat)
			responseWriter.Header().Set("Content-Type", "application/xml;charset=utf-8")
			responseWriter.WriteHeader(http.StatusOK)
			_, err = io.Copy(responseWriter, buffer)
			return err
		}
	}
	responseWriter.WriteHeader(http.StatusOK)
	return nil
}

func (this *WxMp) messageAsync(rc chan *MixMessage) {
	this.Logger.Debug("timeout , to message send sync")
	select {
	case msg := <-rc:
		if msg != nil {
			var err error
			switch msg.MessageType {
			case MessageText:
				err = this.SendTextToUser(msg.ToUserName, msg.Content)
			case MessageImage:
				err = this.SendImageToUser(msg.ToUserName, msg.Image)
			case MessageVoice:
				err = this.SendVoiceToUser(msg.ToUserName, msg.Voice)
			case MessageVideo:
				err = this.SendVideoToUser(msg.ToUserName, msg.Video, "")
			case MessageMusic:
				err = this.SendMusicToUser(msg.ToUserName, msg.Music)
			case MessageNews:
				err = this.SendNewsToUser(msg.ToUserName, msg.Articles)
			default:
				err = errors.New(fmt.Sprintf("unsurpport message type for sync message send :%v", msg.MessageType))
			}
			if err != nil {
				this.Logger.WithField("error", err).Error("send message sync")
			}
		}
	}
}

func (this *WxMp) messageRoute(inMsg *MixMessage) *MixMessage {
	if inMsg == nil {
		return nil
	}
	var handle MessageHandler
	switch inMsg.MessageType {
	case MessageText:
		handle = this.TextMessageHandler
	case MessageImage:
		handle = this.ImageMessageHandler
	case MessageVoice:
		handle = this.VoiceMessageHandler
	case MessageVideo:
		handle = this.VideoMessageHandler
	case MessageShortVideo:
		handle = this.ShortVideoMessageHandler
	case MessageLocation:
		handle = this.LocationMessageHandler
	case MessageLink:
		handle = this.LinkMessageHandler
	case MessageEvent:
		switch inMsg.Event {
		case EventSubscribe:
			handle = this.SubscribeEventHandler
		case EventUnsubscribe:
			handle = this.UnSubscribeEventHandler
		case EventScan:
			handle = this.ScanEventHandler
		case EventLocation:
			handle = this.LocationEventHandler
		case EventClick:
			handle = this.ClickEventHandler
		case EventView:
			handle = this.ViewEventHandler
		default:
			handle = nil
		}
	default:
		handle = nil
	}
	if handle == nil {
		return nil
	}
	return handle(inMsg)
}

func (this *WxMp) messageSignature(request *http.Request, responseWriter http.ResponseWriter) error {
	values := request.URL.Query()
	signature := values.Get("signature")
	timestamp := values.Get("timestamp")
	nonce := values.Get("nonce")
	echoStr := values.Get("echostr")
	this.Logger.WithFields(logrus.Fields{
		"signature": signature,
		"timestamp": timestamp,
		"nonce":     nonce,
		"echostr":   echoStr,
	}).Debug("message signature")
	strs := []string{this.config.Token, timestamp, nonce}
	sort.Strings(strs)
	tempStr := strings.Join(strs, "")
	c := sha1.New()
	_, err := io.WriteString(c, tempStr)
	if err != nil {
		return err
	}
	dat := c.Sum(nil)
	afterSign := hex.EncodeToString(dat[:])
	this.Logger.WithField("aftersign", afterSign).Debug("sign")
	if afterSign == signature {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(echoStr))
	} else {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("sign error"))
	}
	return nil
}
