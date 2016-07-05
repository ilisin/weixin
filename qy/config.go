package qy

import "github.com/ilisin/configuration"

type WxQyConfig struct {
	//weixin service base url
	ServiceURL string `conf:"weixin.qy.service.url,default(https://qyapi.weixin.qq.com)"`
	//weixin corpid,record in http://qydev.weixin.qq.com
	CorpID string `conf:"weixin.qy.corpid"`
	//weixin Secret
	Secret string `conf:"weixin.mp.secret"`
}

func NewWxQyConfig() (*WxQyConfig, error) {
	conf := &WxQyConfig{}
	err := configuration.Var(conf)
	return conf, err
}
