package mp

import "github.com/ilisin/configuration"

type WxMpConfig struct {
	//weixin service base url
	ServiceURL string `conf:"weixin.mp.service.url,default(https://api.weixin.qq.com)"`
	//weixin AppId,record in http://mp.weixin.qq.com
	AppId string `conf:"weixin.mp.appid"`
	//weixin Secret,pair with weixin appId
	Secret string `conf:"weixin.mp.secret"`
	//Token ,with client url
	Token string `conf:"weixin.mp.token"`
	//logger level
	LoggerLevel string `conf:"weixin.mp.logger.level,default(INFO)"`
}

func NewWxMpConfig() (*WxMpConfig, error) {
	conf := &WxMpConfig{}
	err := configuration.Var(conf)
	return conf, err
}
