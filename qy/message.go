package qy

import (
	"encoding/json"
	"errors"

	"github.com/Sirupsen/logrus"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
)

type Message struct {
	// milt with | slit
	ToUser  string      `json:"touser"`
	ToParty string      `json:"toparty"`
	ToTag   string      `json:"totag"`
	MsgType MessageType `json:"msgtype"`
	AgentId int         `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Safe int `json:"safe"`
}

type MessageResponse struct {
	ErrCode      WeiXinResult `json:"errcode"`
	ErrMsg       string       `json:"errmsg"`
	InvalidUser  string       `json:"invaliduser"`
	InvalidParty string       `json:"invalidparty"`
	InvalidTag   string       `json:"invalidtag"`
}

// send a message to a party
func (this *WxQy) SendTextMessageToParty(message, partId string) error {
	msg := &Message{}
	msg.ToParty = partId
	///msg.ToUser = "gaoguangting"
	msg.MsgType = MessageTypeText
	msg.Text.Content = message
	data, err := this.HttpPost("/cgi-bin/message/send", true, nil, msg)
	if err != nil {
		logrus.WithField("data", string(data)).Errorf("http post response error")
		return err
	}
	resp := &MessageResponse{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		logrus.WithField("err", err).Errorf("send message response json encoding error")
	}
	if resp.ErrCode != WXResultSeccuse {
		logrus.WithFields(logrus.Fields{
			"errcode":      resp.ErrCode,
			"errmsg":       resp.ErrMsg,
			"invaliduser":  resp.InvalidUser,
			"invalidparty": resp.InvalidParty,
			"invalidtag":   resp.InvalidTag,
		}).Errorf("send message error")
		return errors.New(resp.ErrMsg)
	}
	return nil
}
