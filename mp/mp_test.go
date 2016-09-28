package mp

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/ilisin/configuration"
	"os"
	"testing"
)

var config = &WxMpConfig{}

func init() {
	os.Setenv("weixin.mp.appid", "wx1dfa40a36b0fef69")
	os.Setenv("weixin.mp.secret", "c44a8d989c5d9459d55e5835946b26de")
	os.Setenv("weixin.mp.token", "ilisiwx")
	os.Setenv("GLOBAL_CONF", "env:://")
	err := configuration.Var(config)
	if err != nil {
		logrus.Error(err)
	}
	logrus.SetLevel(logrus.DebugLevel)
}

func TestWxMp_GetToken(t *testing.T) {
	return
	wx := NewWxMp(config)
	to, err := wx.Token()
	if err != nil {
		t.Error(err)
	}
	t.Log(to, to.AccessToken)
	logrus.Info(to, "xx", to.AccessToken)
}

func TestGetMenu(t *testing.T) {
	return
	wx := NewWxMp(config)
	menu, err := wx.GetMenu()
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", menu)
	}
}

func TestWxMp_DeleteMenu(t *testing.T) {
	return
	wx := NewWxMp(config)
	err := wx.DeleteMenu()
	if err != nil {
		t.Error(err)
	}
}

func TestSetMenu(t *testing.T) {
	return
	str := `{"button":[{"name":"点我","sub_button":[{"type":"click","name":"添加菜单","key":"vip"},{"type":"click","name":"电子会员卡","key":"vip"},{"type":"view","name":"会员特权","url":"http://weixintest.imooly.com/show/"},{"type":"view","name":"充值卡","url":"http://weixintest.imooly.com/account/card"}]},{"name":"随便点","sub_button":[{"type":"click","name":"热门店铺","key":"HOTBUSINESS"},{"type":"click","name":"附近店铺","key":"NEARBUSINESS"}]},{"name":"不要点我","sub_button":[{"type":"view","name":"常见问题","url":"http://weixintest.imooly.com/show/problem"},{"type":"view","name":"调戏客服","url":"http://weixintest.imooly.com/show/contact"}]}]}`
	menub := Menu{}
	err := json.Unmarshal([]byte(str), &menub)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", menub)
	wx := NewWxMp(config)
	err = wx.SetMenu(&menub)
	if err != nil {
		t.Error(err)
	}
}

func TestSendText(t *testing.T) {
	return
	wx := NewWxMp(config)
	err := wx.SendTextToUser("o8vf9t8qguhcJfgSpFlgIrm63iXI", "你好啊")
	if err != nil {
		t.Error(err)
	}
}
