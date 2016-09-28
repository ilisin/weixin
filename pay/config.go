package pay

import (
	"github.com/ilisin/configuration"
)

type WxPayConfig struct {
	MchBashUrl string `conf:"weixin.pay.service.url,default(https://api.mch.weixin.qq.com)"`
	AppId      string `conf:"weixin.pay.appid"`
	MctId      string `conf:"weixin.pay.merchant.id"`
	SubMctId   string `conf:"weixin.pay.merchant.subid"`
	MctName    string `conf:"weixin.pay.merchant.name"`
	ApiKey     string `conf:"weixin.pay.apikey"`
	Cert       struct {
		Ca   string `conf:"weixin.pay.cert.ca"`
		Cert string `conf:"weixin.pay.cert.cert"`
		Key  string `conf:"weixin.pay.cert.key"`
	}
	//logger level
	LoggerLevel string `conf:"weixin.pay.logger.level,default(INFO)"`
}

func NewWxPayConfig() (*WxPayConfig, error) {
	conf := &WxPayConfig{}
	err := configuration.Var(conf)
	return conf, err
}
