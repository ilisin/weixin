package mp

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
	"github.com/ilisin/weixin/util"
)

type JsapiTicket struct {
	Ticket     string    `json:"ticket"`
	ExpiresIn  int       `json:"expires_in"`
	ExpireTime time.Time `json:"-"`
}

func (this *JsapiTicket) IsExpired() bool {
	return this.ExpireTime.Before(time.Now())
}

//避免直接调用，此时有限制
func (this *WxMp) GetJsapiTicket() (*JsapiTicket, error) {
	jsTickt := &JsapiTicket{}
	err := this.HttpGet("/cgi-bin/ticket/getticket", map[string]interface{}{
		"type": "jsapi",
	}, jsTickt)
	if err != nil {
		return nil, err
	}
	return jsTickt, err
}

func (this *WxMp) JsapiToken() (*JsapiTicket, error) {
	exist := false
	if this.jsToken != nil && this.jsToken.IsExpired() == false {
		exist = true
	}
	if !exist {
		return this.GetJsapiTicket()
	}
	return this.jsToken, nil
}

func (this *WxMp) JsSignature(urll string, timestamp int64) (nonceStr, sign string, err error) {
	//#后不参与签名
	unUrl, _ := url.QueryUnescape(urll)
	if i := strings.Index(unUrl, "#"); i > 0 {
		unUrl = string(unUrl[:i])
	}
	nonceStr = util.ProductANonceString()
	if _, err = this.JsapiToken(); err != nil {
		return
	}
	jsTickt := this.jsToken.Ticket
	strTemp := fmt.Sprintf("jsapi_ticket=%v&noncestr=%v&timestamp=%v&url=%v", jsTickt, nonceStr, timestamp, unUrl)
	c := sha1.New()
	_, err = io.WriteString(c, strTemp)
	if err != nil {
		return
	}
	dat := c.Sum(nil)
	sign = hex.EncodeToString(dat[:])
	return
}
