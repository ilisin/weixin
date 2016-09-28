package mp

import (
	"encoding/json"
	"time"
)

type WeiXinToken struct {
	AccessToken string    `json:"access_token"`
	ExpireTime  time.Time `json:"expires_in"`
}

type weiXinTokenJson struct {
	AccessToken string `json:"access_token"`
	ExpireIn    int    `json:"expires_in"`
}

//for json unmarshal
func (this *WeiXinToken) UnmarshalJSON(data []byte) error {
	var token = weiXinTokenJson{}
	err := json.Unmarshal(data, &token)
	if err != nil {
		return err
	}
	this.AccessToken = token.AccessToken
	this.ExpireTime = time.Now().Add(time.Second * time.Duration(token.ExpireIn))
	return err
}

func (token WeiXinToken) IsExpired() bool {
	return token.ExpireTime.Before(time.Now())
}

//take a token from weixin server
func (this *WxMp) getToken() (*WeiXinToken, error) {
	token := &WeiXinToken{}
	err := this.HttpClient.HttpGetJson("/cgi-bin/token", map[string]interface{}{
		"grant_type": "client_credential",
		"appid":      this.config.AppId,
		"secret":     this.config.Secret,
	}, token)
	if err != nil {
		return nil, err
	}
	this.token = token
	return this.token, err
}

//take token , if not exist query from weixin server
func (this *WxMp) Token() (*WeiXinToken, error) {
	exist := false
	if this.token != nil && this.token.IsExpired() == false {
		exist = true
	}
	if !exist {
		return this.getToken()
	}
	return this.token, nil
}
