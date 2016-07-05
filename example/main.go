package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ilisin/weixin/mp"
	"github.com/ilisin/weixin/pay"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	"time"
)

var config = &mp.WxMpConfig{}
var wx *mp.WxMp

func Handler(c *echo.Context) error {
	return wx.MessageFilter(c.Request(), c.Response().Writer(), func(message *mp.MixMessage) *mp.MixMessage {
		replayMessage := &mp.MixMessage{
			MessageHeader: mp.MessageHeader{
				ToUserName:   message.FromUserName,
				FromUserName: message.ToUserName,
				CreateTime:   time.Now().Unix(),
				MessageType:  mp.MessageText,
			},
		}
		user, err := wx.GetUser(message.FromUserName)
		if err != nil {
			replayMessage.Content = "未查找到您的个人信息,请重新发送"
		} else {
			if message.MessageType == mp.MessageText {
				replayMessage.Content = fmt.Sprintf("您好：%v\n您刚对我说了:%v", user.NickName, message.Content)
			} else {
				replayMessage.Content = fmt.Sprintf("对不起，我不知道您说的是什么")
			}
		}
		return replayMessage
	})
}

func GetHandle(c *echo.Context) error {
	return wx.MessageSignature(c.Request(), c.Response().Writer())
}

func Exec(c *echo.Context) error {
	//	MchBashUrl string `conf:"weixin.pay.service.url,default(https://api.mch.weixin.qq.com)"`
	//	AppId      string `conf:"weixin.pay.appid"`
	//	MctId      string `conf:"weixin.pay.merchant.id"`
	//	MctName    string `conf:"weixin.pay.merchant.name"`
	//	ApiKey     string `conf:"weixin.pay.apikey"`
	//	Cert       struct {
	//Ca   string `conf:"weixin.pay.cert.ca"`
	//Cert string `conf:"weixin.pay.cert.cert"`
	//Key  string `conf:"weixin.pay.cert.key"`
	//}
	conf := pay.NewWxPayConfig()
	//var config = weixin.NewWeiXinConfig("wx6dfeebdd15859916", "4850b8b2bace833c98f968519d4b2c84", "moolyweixin")
	wxpay := pay.NewWxPay(conf)

	//wx.HttpTlsPost("/mmpaymkttransfers/gethbinfo")
	err := wxpay.SendRedPack("olPlQuL1MCgkTLH6mcdzvqThOsng", 10)
	if err != nil {
		logrus.Error(err)
		return c.String(200, fmt.Sprintf("%v", err))
	}
	return c.String(200, "ok")
}

func init() {
	os.Setenv("weixin.mp.appid", "xxxx")
	os.Setenv("weixin.mp.secret", "d4xxx")
	os.Setenv("weixin.mp.token", "xx")
	os.Setenv("GLOBAL_CONF", "path:://")
	err := configuration.Var(config)
	if err != nil {
		logrus.Error(err)
	}
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	wx = mp.NewWxMp(config)

	e := echo.New()
	e.Use(middleware.Logger())

	e.Get("/weixin", GetHandle)
	e.Post("/weixin", Handler)

	e.Get("/exec", Exec)

	e.Run(":8080")
}
