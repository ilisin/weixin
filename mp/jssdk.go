package mp

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"net/url"
	"strings"
	"time"
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
	dat, err := this.HttpGet("/cgi-bin/ticket/getticket", true, map[string]interface{}{
		"type": "jsapi",
	})
	if err != nil {
		return nil, err
	}
	jsTickt := &JsapiTicket{}
	err = json.Unmarshal(dat, jsTickt)
	if err == nil {
		jsTickt.ExpireTime = time.Now().Add(time.Second * time.Duration(jsTickt.ExpiresIn))
		this.jsToken = jsTickt
		logrus.Debug("验签不存在， 重新获取:", *this.jsToken)
	}
	return jsTickt, err
}

func (this *WxMp) JsapiToken() (*JsapiTicket, error) {
	exist := false
	if this.jsToken != nil && this.jsToken.IsExpired() == false {
		logrus.Debug("验签存在:", *this.jsToken)
		exist = true
	}
	if !exist {
		this.Token()
		return this.GetJsapiTicket()
	}
	return this.jsToken, nil
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func productANonceString() string {
	str := fmt.Sprintf("%v", time.Now().Nanosecond())
	return GetMd5String(str)
}

func (this *WxMp) JsSignature(urll string, timestamp int64) (nonceStr, sign string, err error) {
	//#后不参与签名
	unUrl, _ := url.QueryUnescape(urll)
	if i := strings.Index(unUrl, "#"); i > 0 {
		unUrl = string(unUrl[:i])
	}
	nonceStr = productANonceString()
	if _, err = this.JsapiToken(); err != nil {
		return
	}
	logrus.Debugf("url:%v", unUrl)
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
