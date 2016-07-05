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
