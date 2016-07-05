package mp

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
)

var (
	_ = logrus.Debug
)

const (
	AUTHORIZE_URL = "https://open.weixin.qq.com/connect/oauth2/authorize"
)

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"` //用户刷新token
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"` //用户授权的作用域 使用逗号(,)隔开
	UnionId      string `json:"unionid"`
}

//获取auth2.0授权地址
//redirectUrl:授权成功后的重定向地址
//apiBase:为true时Scope为snsapi_base，否则为snsapi_userinfo
//state：附加参数，回传给从定向地址
func (this *WxMp) GetAuth2Url(redirectUrl string, apiBase bool, state string) string {
	scope := "snsapi_userinfo"
	if apiBase {
		scope = "snsapi_base"
	}
	return fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect", AUTHORIZE_URL, this.config.AppId, redirectUrl, scope, state)
}

func (this *WxMp) Auth2Token(code string) (*AuthToken, error) {
	dat, err := this.HttpGet("/sns/oauth2/access_token", false, map[string]interface{}{
		"appid":      this.config.AppId,
		"secret":     this.config.Secret,
		"code":       code,
		"grant_type": "authorization_code",
	})
	if err != nil {
		return nil, err
	}
	authToken := &AuthToken{}
	err = json.Unmarshal(dat, authToken)
	if err != nil {
		return nil, err
	}
	return authToken, nil
}

func (this *WxMp) AuthRefreshToken(token *AuthToken) error {
	dat, err := this.HttpGet("/sns/oauth2/refresh_token", false, map[string]interface{}{
		"appid":         this.config.AppId,
		"grant_type":    "refresh_token",
		"refresh_token": token.RefreshToken,
	})
	if err != nil {
		return err
	}
	logrus.WithField("data", string(dat)).Debug("refresh token")
	err = json.Unmarshal(dat, token)
	if err != nil {
		return err
	}
	return nil
}
