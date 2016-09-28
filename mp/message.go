package mp

import (
	"encoding/xml"
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
	MediaId string `xml:"MediaId" json:"media_id"`
}

type MessageVoiceItem struct {
	MediaId string `xml:"MediaId" json:"media_id"`
}

type MessageVideoItem struct {
	MediaId     string `xml:"MediaId"`
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
}

type MessageMusicItem struct {
	Title        string `xml:"Title" json:"title"`
	Description  string `xml:"Description" json:"description"`
	MusicUrl     string `xml:"MusicUrl" json:"musicurl"`
	HQMusicUrl   string `xml:"HQMusicUrl" json:"hqmusicurl"`
	ThumbMediaId string `xml:"ThumbMediaId" json:"thumb_media_id"`
}

//图文消息
type MessageNewsItem struct {
	Title       string `xml:"Title,omitempty" json:"title"`
	Description string `xml:"Description,omitempty" json:"description"`
	PicUrl      string `xml:"PicUrl,omitempty" json:"picurl"`
	Url         string `xml:"Url,omitempty" json:"url"`
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

type MessageHandler func(*MixMessage) *MixMessage

// send text to user
func (this *WxMp) SendTextToUser(openId, text string) error {
	msg := map[string]interface{}{
		"touser":  openId,
		"msgtype": MessageText,
		"text": map[string]string{
			"content": text,
		},
	}
	return this.HttpPostWithCommonResponse("/cgi-bin/message/custom/send", nil, msg)
}

func (this *WxMp) SendImageToUser(openId string, image *MessageImageItem) error {
	msg := map[string]interface{}{
		"touser":  openId,
		"msgtype": MessageImage,
		"image":   image,
	}
	return this.HttpPostWithCommonResponse("/cgi-bin/message/custom/send", nil, msg)
}

func (this *WxMp) SendVoiceToUser(openId string, voice *MessageVoiceItem) error {
	msg := map[string]interface{}{
		"touser":  openId,
		"msgtype": MessageVoice,
		"voice":   voice,
	}
	return this.HttpPostWithCommonResponse("/cgi-bin/message/custom/send", nil, msg)
}

func (this *WxMp) SendVideoToUser(openId string, video *MessageVideoItem, thumbId string) error {
	msg := map[string]interface{}{
		"touser":  openId,
		"msgtype": MessageVideo,
		"video": map[string]string{
			"media_id":       video.MediaId,
			"thumb_media_id": thumbId,
			"title":          video.Title,
			"description":    video.Description,
		},
	}
	return this.HttpPostWithCommonResponse("/cgi-bin/message/custom/send", nil, msg)
}

func (this *WxMp) SendMusicToUser(openId string, music *MessageMusicItem) error {
	msg := map[string]interface{}{
		"touser":  openId,
		"msgtype": MessageMusic,
		"music":   music,
	}
	return this.HttpPostWithCommonResponse("/cgi-bin/message/custom/send", nil, msg)
}

func (this *WxMp) SendNewsToUser(openId string, news []MessageNewsItem) error {
	msg := map[string]interface{}{
		"touser":  openId,
		"msgtype": MessageNews,
		"news": map[string]interface{}{
			"articles": news,
		},
	}
	return this.HttpPostWithCommonResponse("/cgi-bin/message/custom/send", nil, msg)
}

func (this *WxMp) RegisterTextMessageHandle(handle MessageHandler) {
	this.TextMessageHandler = handle
}

func (this *WxMp) RegisterImageMessageHandle(handle MessageHandler) {
	this.ImageMessageHandler = handle
}

func (this *WxMp) RegisterVoiceMessageHandle(handle MessageHandler) {
	this.VoiceMessageHandler = handle
}

func (this *WxMp) RegisterVideoMessageHandle(handle MessageHandler) {
	this.VideoMessageHandler = handle
}

func (this *WxMp) RegisterShortVideoMessageHandle(handle MessageHandler) {
	this.ShortVideoMessageHandler = handle
}

func (this *WxMp) RegisterLocationMessageHandle(handle MessageHandler) {
	this.LocationMessageHandler = handle
}

func (this *WxMp) RegisterLinkMessageHandle(handle MessageHandler) {
	this.LinkMessageHandler = handle
}
