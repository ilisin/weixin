package mp

import (
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
	Response
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"` //用户刷新token
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"` //用户授权的作用域 使用逗号(,)隔开
	UnionId      string `json:"unionid"`
}

type WXAuthScope string

const (
	WXAuthScopeUserInfo WXAuthScope = "snsapi_userinfo"
	WxAuthScopeBase     WXAuthScope = "snsapi_base"
)

// get auth2.0 redirect url
func (this *WxMp) GetAuth2Url(redirectUrl string, scope WXAuthScope, state string) string {
	return fmt.Sprintf("%s?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect", AUTHORIZE_URL, this.config.AppId, redirectUrl, scope, state)
}

func (this *WxMp) Auth2Token(code string) (*AuthToken, error) {
	token := &AuthToken{}
	err := this.HttpClient.HttpGetJson("/sns/oauth2/access_token", map[string]interface{}{
		"appid":      this.config.AppId,
		"secret":     this.config.Secret,
		"code":       code,
		"grant_type": "authorization_code",
	}, token)
	if err != nil {
		return nil, err
	}
	if er := token.Error(); er != nil {
		return nil, er
	}
	return token, nil
}

func (this *WxMp) AuthRefreshToken(token *AuthToken) error {
	respToken := &AuthToken{}
	err := this.HttpClient.HttpGetJson("/sns/oauth2/refresh_token", map[string]interface{}{
		"appid":         this.config.AppId,
		"grant_type":    "refresh_token",
		"refresh_token": token.RefreshToken,
	}, respToken)
	if err != nil {
		return err
	}
	if er := respToken.Error(); er != nil {
		return er
	}
	token = respToken
	return nil
}
