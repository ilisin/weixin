package mp

import "github.com/juju/errgo/errors"

type WeiXinResult int

const (
	WXResultSuccess WeiXinResult = 0
	WXResultNonMenu WeiXinResult = 46003 //未设置菜单
)

type Response struct {
	ErrCode WeiXinResult `json:"errcode"`
	ErrMsg  string       `json:"errmsg"`
}

func (r Response) Error() error {
	if r.ErrCode == WXResultSuccess {
		return nil
	}
	return errors.New(r.ErrMsg)
}
