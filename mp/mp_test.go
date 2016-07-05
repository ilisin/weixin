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
	os.Setenv("weixin.mp.appid", "xxxx")
	os.Setenv("weixin.mp.secret", "xxxxxxxxx")
	os.Setenv("weixin.mp.token", "xxxxxxxxxx")
	os.Setenv("GLOBAL_CONF", "path:://")
	err := configuration.Var(config)
	if err != nil {
		logrus.Error(err)
	}
	logrus.SetLevel(logrus.DebugLevel)
}

func TestGetMenu(t *testing.T) {
	//return
	wx := NewWxMp(config)
	menu, err := wx.GetMenu()
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", menu)
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

func TestUpdateMenu(t *testing.T) {
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
	err = wx.UpdateMenu(&menub)
	if err != nil {
		t.Error(err)
	}
}

func TestUserCount(t *testing.T) {
	return
	wx := NewWxMp(config)
	count, err := wx.UserCount()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", count)
}

func TestGetUser(t *testing.T) {
	return
	wx := NewWxMp(config)
	user, err := wx.GetUser("o8vf9t8qguhcJfgSpFlgIrm63iXI")
	if err != nil {
		t.Error(err)
		return
	}
	logrus.WithField("user", user).Info("get user")
}

func TestReadUsers(t *testing.T) {
	return
	wx := NewWxMp(config)
	users, err := wx.ReadUsers([]string{"o8vf9t8qguhcJfgSpFlgIrm63iXI", "o8vf9t2VMALDHvSZk5P8UHXWzQP0"})
	if err != nil {
		t.Error(err)
		return
	}
	logrus.WithField("users", users).Info("read users")
}

func TestReadUserIds(t *testing.T) {
	return
	wx := NewWxMp(config)
	ids, err := wx.ReadUserOpenIds("")
	if err != nil {
		t.Error(err)
		return
	}
	logrus.WithField("userOpenIds", ids).Info("read user openIds")
}

func TestUserUpdateRemark(t *testing.T) {
	return
	wx := NewWxMp(config)
	err := wx.UserUpdateRemark("o8vf9t8qguhcJfgSpFlgIrm63iXI", "gaogt")
	if err != nil {
		t.Error(err)
		return
	}
	logrus.Info("update remark seccuess")
}

func TestSendText(t *testing.T) {
	return
	wx := NewWxMp(config)
	err := wx.SendTextToUser("o8vf9t8qguhcJfgSpFlgIrm63iXI", "你好啊")
	if err != nil {
		t.Error(err)
	}
}
