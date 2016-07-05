package qy

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

var config = &WxQyConfig{
	ServiceURL: "https://qyapi.weixin.qq.com",
	CorpID:     "xxxxxxxxxxxx",
	Secret:     "xxxxxxxxxxxxxxxxxx",
}

func TestMessage(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	qy := NewWxQy(config)
	err := qy.SendTextMessageToParty("发送测试", "1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("send success")
}
