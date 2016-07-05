package qy

import (
	"encoding/json"
)

type WxQy struct {
	// weixin qy config
	config *WxQyConfig
	// access token
	token *WXQyToken
}

//take a token from weixin server
func (this *WxQy) GetToken() (*WXQyToken, error) {
	dat, err := this.HttpGet("/cgi-bin/gettoken", false, map[string]interface{}{
		"corpid":     this.config.CorpID,
		"corpsecret": this.config.Secret,
	})
	if err != nil {
		return nil, err
	}
	this.token = &WXQyToken{}
	err = json.Unmarshal(dat, this.token)
	return this.token, err
}

//take token , if not exist query from weixin server
func (this *WxQy) Token() (*WXQyToken, error) {
	exist := false
	if this.token != nil && this.token.IsExpired() == false {
		exist = true
	}
	if !exist {
		return this.GetToken()
	}
	return this.token, nil
}

func NewWxQy(config *WxQyConfig) *WxQy {
	return &WxQy{
		config: config,
	}
}
