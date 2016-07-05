package mp

import (
	"encoding/json"
)

type WxMp struct {
	config  *WxMpConfig
	token   *WeiXinToken
	jsToken *JsapiTicket
}

//take a token from weixin server
func (this *WxMp) GetToken() (*WeiXinToken, error) {
	dat, err := this.HttpGet("/cgi-bin/token", false, map[string]interface{}{
		"grant_type": "client_credential",
		"appid":      this.config.AppId,
		"secret":     this.config.Secret,
	})
	if err != nil {
		return nil, err
	}
	this.token = &WeiXinToken{}
	err = json.Unmarshal(dat, this.token)
	return this.token, err
}

//take token , if not exist query from weixin server
func (this *WxMp) Token() (*WeiXinToken, error) {
	exist := false
	if this.token != nil && this.token.IsExpired() == false {
		exist = true
	}
	if !exist {
		return this.GetToken()
	}
	return this.token, nil
}

func NewWxMp(config *WxMpConfig) *WxMp {
	return &WxMp{
		config: config,
	}
}
