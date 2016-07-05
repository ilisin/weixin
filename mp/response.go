package mp

type WeiXinResult int

const (
	WXResultSeccuse WeiXinResult = 0
	WXResultNonMenu WeiXinResult = 46003 //未设置菜单
)

type Response struct {
	ErrCode WeiXinResult `json:"errcode"`
	ErrMsg  string       `json:"errmsg"`
}
