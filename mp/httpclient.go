package mp

import (
	"errors"
)

var ErrToken = errors.New("cann't found the access token")

func (this *WxMp) HttpGet(path string, params map[string]interface{}, respObj interface{}) error {
	token, err := this.Token()
	if err != nil {
		return ErrToken
	}
	if params == nil {
		params = map[string]interface{}{
			"access_token": token.AccessToken,
		}
	} else {
		params["access_token"] = token.AccessToken
	}
	return this.HttpClient.HttpGetJson(path, params, respObj)
}

func (this *WxMp) HttpGetWithCommonResponse(path string, params map[string]interface{}) error {
	resp := &Response{}
	err := this.HttpGet(path, params, resp)
	if err != nil {
		return err
	}
	if resp.ErrCode != WXResultSuccess {
		return errors.New(resp.ErrMsg)
	}
	return nil
}

// params is the part of path params
// bodyObj is the request data,with json format
func (this *WxMp) HttpPost(path string, params map[string]interface{}, bodyObj interface{}, respObj interface{}) error {
	if token, err := this.Token(); err == nil {
		params = map[string]interface{}{
			"access_token": token.AccessToken,
		}
	} else {
		return ErrToken
	}
	return this.HttpClient.HttpPostJson(path, params, bodyObj, respObj)
}

func (this *WxMp) HttpPostWithCommonResponse(path string, params map[string]interface{}, bodyObj interface{}) error {
	if token, err := this.Token(); err == nil {
		params = map[string]interface{}{
			"access_token": token.AccessToken,
		}
	} else {
		return ErrToken
	}
	resp := Response{}
	err := this.HttpClient.HttpPostJson(path, params, bodyObj, &resp)
	if err != nil {
		return err
	}
	if resp.ErrCode != WXResultSuccess {
		return errors.New(resp.ErrMsg)
	}
	return nil
}
