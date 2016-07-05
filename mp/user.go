package mp

import (
	"encoding/json"
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
	Total      int                 `json:"total"`
	Count      int                 `json:"count"`
	Data       map[string][]string `json:"data"`
	NextOpenId string              `json:"next_openid"`
}

type UserList struct {
	UserInfoList []*User `json:"user_info_list"`
}

func (this *WxMp) UserCount() (count int, err error) {
	dat, err := this.HttpGet("/cgi-bin/user/get", true, nil)
	if err != nil {
		return 0, err
	}
	userList := UserOpenIdList{}
	err = json.Unmarshal(dat, &userList)
	if err != nil {
		return 0, err
	}
	return userList.Total, nil
}

func (this *WxMp) GetUser(openId string) (user *User, err error) {
	dat, err := this.HttpGet("/cgi-bin/user/info", true, map[string]interface{}{
		"openid": openId,
		"lang":   LanguageZhCN,
	})
	if err != nil {
		return nil, err
	}
	user = &User{}
	err = json.Unmarshal(dat, user)
	return user, err
}

//获取用户openId数组，最多10000个，nextOpenId传值，则从对应的用户查找
func (this *WxMp) ReadUserOpenIds(nextOpenId string) (openIds []string, err error) {
	var params map[string]interface{} = nil
	if len(nextOpenId) > 0 {
		params = map[string]interface{}{
			"next_openid": nextOpenId,
		}
	}
	dat, err := this.HttpGet("/cgi-bin/user/get", true, params)
	if err != nil {
		return nil, err
	}
	userList := UserOpenIdList{}
	err = json.Unmarshal(dat, &userList)
	if err != nil {
		return nil, err
	}
	if x, ok := userList.Data["openid"]; ok {
		return x, nil
	}
	return nil, ErrFetch
}

//上限未100个
func (this *WxMp) ReadUsers(openIds []string) (users []*User, err error) {
	array := make([]map[string]interface{}, len(openIds))
	for i, openId := range openIds {
		array[i] = map[string]interface{}{
			"openid": openId,
			"lang":   LanguageZhCN,
		}
	}
	dat, err := this.HttpPost("/cgi-bin/user/info/batchget", true, nil, map[string]interface{}{
		"user_list": array,
	})
	if err != nil {
		return nil, err
	}
	userList := UserList{}
	err = json.Unmarshal(dat, &userList)
	users = userList.UserInfoList
	return
}

//修改用户的备注
func (this *WxMp) UserUpdateRemark(openId, remark string) error {
	_, err := this.HttpPost("/cgi-bin/user/info/updateremark", true, nil, map[string]string{
		"openid": openId,
		"remark": remark,
	})
	return err
}
