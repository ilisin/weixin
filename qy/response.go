package qy

type WeiXinResult int

const (
	WXResultSeccuse       WeiXinResult = 0
	WXResultTakenTokenErr WeiXinResult = 43003 //未设置菜单
)

type Response struct {
	ErrCode WeiXinResult `json:"errcode"`
	ErrMsg  string       `json:"errmsg"`
}
