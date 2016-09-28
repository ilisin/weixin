package mp

import (
	"errors"

	"github.com/Sirupsen/logrus"
)

var (
	ErrFetch = errors.New("Cann't fetch response")

	_ = logrus.Debug
)

//性别
type Sex int

const (
	Unkown Sex = iota
	Male
	Female
)

type LanguageType string

const (
	LanguageZhCN LanguageType = "zh_CN"
	LanguageZhTW LanguageType = "zh_TW"
	LanguageEN   LanguageType = "en"
)

type User struct {
	Response
	Subscribe     int          `json:"subscribe"` //1为已关注，0为未关注，拉取不到其他信息
	OpenId        string       `json:"openid"`
	NickName      string       `json:"nickname"`
	Sex           Sex          `json:"sex"`
	Language      LanguageType `json:"language"` //ex : zh_CN
	City          string       `json:"city"`
	Province      string       `json:"province"`
	Country       string       `json:"country"`
	Avatar        string       `json:"headimgurl"`
	SubscribeTime int64        `json:"subscribe_time"`
	UnionId       string       `json:"unionid"`
	Remark        string       `json:"remark"` //对用户的备注
	GroupId       int          `json:"groupid"`
}

type UserOpenIdList struct {
	Response
	Total      int                 `json:"total"`
	Count      int                 `json:"count"`
	Data       map[string][]string `json:"data"`
	NextOpenId string              `json:"next_openid"`
}

type UserList struct {
	UserInfoList []*User `json:"user_info_list"`
}

func (this *WxMp) UserCount() (count int, err error) {
	resp := &UserOpenIdList{}
	err = this.HttpGet("/cgi-bin/user/get", nil, resp)
	if err != nil {
		return 0, err
	}
	return resp.Total, nil
}

func (this *WxMp) GetUser(openId string) (user *User, err error) {
	user = &User{}
	err = this.HttpGet("/cgi-bin/user/info", map[string]interface{}{
		"openid": openId,
		"lang":   LanguageZhCN,
	}, user)
	if err != nil {
		return nil, err
	}
	return user, err
}
