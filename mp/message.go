package mp

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type MessageType string

const (
	MessageText       MessageType = "text"
	MessageEvent      MessageType = "event"
	MessageImage      MessageType = "image"
	MessageVoice      MessageType = "voice"
	MessageVideo      MessageType = "video"
	MessageShortVideo MessageType = "shortvideo"
	MessageLocation   MessageType = "location"
	MessageLink       MessageType = "link"
	MessageMusic      MessageType = "music"
	MessageNews       MessageType = "news"
)

type EventType string

const (
	EventSubscribe   EventType = "subscribe"   //订阅事件
	EventUnsubscribe EventType = "unsubscribe" //取消订阅
	EventScan        EventType = "SCAN"
	EventLocation    EventType = "LOCATION"
	EventClick       EventType = "CLICK"
	EventView        EventType = "VIEW" //菜单跳转链接事件
)

type MessageImageItem struct {
	MediaId string
}

type MessageVoiceItem struct {
	MediaId string
}

type MessageVideoItem struct {
	MediaId     string
	Title       string
	Description string
}

type MessageMusicItem struct {
	Title        string
	Description  string
	MusicUrl     string
	HQMusicUrl   string
	ThumbMediaId string
}

//图文消息
type MessageNewsItem struct {
	Title       string `xml:"Title,omitempty"`
	Description string `xml:"Description,omitempty"`
	PicUrl      string `xml:"PicUrl,omitempty"`
	Url         string `xml:"Url,omitempty"`
}

type MessageHeader struct {
	ToUserName   string      `xml:"ToUserName"`
	FromUserName string      `xml:"FromUserName"`
	CreateTime   int64       `xml:"CreateTime"`
	MessageType  MessageType `xml:"MsgType"`
}

type MixMessage struct {
	XMLName xml.Name `xml:"xml"`

	MessageHeader
	MessageId    int64   `xml:"MsgId"`
	Content      string  `xml:"Content"`
	MediaId      string  `xml:"MediaId"`
	PicUrl       string  `xml:"PicUrl"`
	Format       string  `xml:"Format"` //语音格式amr,speex
	Recognition  string  `xml:"Recognition"`
	ThumbMediaId string  `xml:"ThumbMediaId"` //视频消息缩略图的媒体id
	LocationX    float32 `xml:"Location_X"`   //纬度
	LocationY    float32 `xml:"Location_Y"`   //经度
	Scale        int     `xml:"Scale"`        //地图缩放大小
	Label        string  `xml:"Label"`        //地理未知信息
	Title        string  `xml:"Title"`        //消息标题
	Description  string  `xml:"Description"`  //消息描述
	Url          string  `xml:"Url"`          //消息链接

	Event EventType `xml:"Event"` //事件类型
	//事件KEY
	//Event为subscribe时,用户未关注时,qrscene_未前缀,后面未二维码参数值
	//Event未CLICK时，未菜单中KEY值
	EventKey  string  `xml:"EventKey"`
	Ticket    string  `xml:"Ticket"` //二维码的ticket值
	Latitude  float32 `xml:"Latitude"`
	Longitude float32 `xml:"Longitude"`
	Precision float32 `xml:"Precision"`

	Image        *MessageImageItem
	Voice        *MessageVoiceItem
	Video        *MessageVideoItem
	Music        *MessageMusicItem
	ArticleCount int
	Articles     []MessageNewsItem `xml:"Articles>item,omitempty"`
}

func (this *WxMp) MessageFilter(request *http.Request, responseWriter http.ResponseWriter, callback func(*MixMessage) *MixMessage) error {
	if request == nil {
		return fmt.Errorf("request paramer error")
	}
	defer request.Body.Close()
	dat, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err
	}
	logrus.WithField("dat", string(dat)).Info("messageFilter")
	message := MixMessage{}
	err = xml.Unmarshal(dat, &message)
	if err != nil {
		return err
	}
	logrus.WithField("message", message).Info("messageFilter")
	repMessage := callback(&message)
	if repMessage != nil {
		repDat, err := xml.Marshal(repMessage)
		logrus.WithField("message", string(repDat)).Info("response message")
		buffer := bytes.NewBuffer(repDat)
		responseWriter.Header().Set("Content-Type", "application/xml;charset=utf-8")
		responseWriter.WriteHeader(http.StatusOK)
		_, err = io.Copy(responseWriter, buffer)
		return err
	} else {
		responseWriter.WriteHeader(http.StatusOK)
	}
	return nil
}

func (this *WxMp) MessageSignature(request *http.Request, responseWriter http.ResponseWriter) error {
	values := request.URL.Query()
	signature := values.Get("signature")
	timestamp := values.Get("timestamp")
	nonce := values.Get("nonce")
	echoStr := values.Get("echostr")
	logrus.WithFields(logrus.Fields{
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
	logrus.WithField("aftersign", afterSign).Debug("sign")
	if afterSign == signature {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(echoStr))
	} else {
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("sign error"))
	}
	return nil
}

//给指定用户放信息
func (this *WxMp) SendTextToUser(openId, text string) error {
	msg := map[string]interface{}{
		"touser":  openId,
		"msgtype": "text",
		"text": map[string]string{
			"content": text,
		},
	}
	_, err := this.HttpPost("/cgi-bin/message/custom/send", true, nil, msg)
	return err
}
