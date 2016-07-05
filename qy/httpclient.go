package qy

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ilisin/weixin"
	"strings"
)

//params is the part of url
func (this *WxQy) HttpGet(path string, withToken bool, params map[string]interface{}) (data []byte, err error) {
	url := fmt.Sprintf("%v%v", this.config.ServiceURL, path)
	if withToken {
		if params == nil {
			params = make(map[string]interface{})
		}
		if token, err := this.Token(); err == nil {
			params["access_token"] = token.AccessToken
		} else {
			return nil, err
		}
	}
	if params != nil {
		pas := make([]string, len(params))
		i := 0
		for k, v := range params {
			pas[i] = fmt.Sprintf("%v=%v", k, v)
			i++
		}
		url = fmt.Sprintf("%v?%v", url, strings.Join(pas, "&"))
	}
	logrus.WithField("url", url).Debug("http get")
	data, err = weixin.HttpGet(url, weixin.FILE_TYPE_JSON)
	resp := Response{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		logrus.Error(`HTTPGet 返回错误`, string(data), err)
		return nil, err
	}
	logrus.WithField("resp", resp).Debug("http get")
	if resp.ErrCode != WXResultSeccuse {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"resp": resp,
		}).Info("HTTPGet")
		return nil, fmt.Errorf("[%v]%v", resp.ErrCode, resp.ErrMsg)
	}
	return data, nil
}

//params is the part of path params
//bodyObj is the request data,with json format
func (this *WxQy) HttpPost(path string, withToken bool, params map[string]interface{}, bodyObj interface{}) ([]byte, error) {
	url := fmt.Sprintf("%v%v", this.config.ServiceURL, path)
	if withToken {
		if params == nil {
			params = make(map[string]interface{})
		}
		if token, err := this.Token(); err == nil {
			params["access_token"] = token.AccessToken
		} else {
			return nil, err
		}
	}
	if params != nil {
		pas := make([]string, len(params))
		i := 0
		for k, v := range params {
			pas[i] = fmt.Sprintf("%v=%v", k, v)
			i++
		}
		url = fmt.Sprintf("%v?%v", url, strings.Join(pas, "&"))
	}
	logrus.WithField("url", url).Debug("http post")
	data, err := weixin.HttpPost(url, bodyObj, weixin.FILE_TYPE_JSON)
	logrus.WithField("data", string(data)).Debugf("http post")
	resp := Response{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != WXResultSeccuse {
		return nil, fmt.Errorf("[%v]%v", resp.ErrCode, resp.ErrMsg)
	}
	return data, nil
}
