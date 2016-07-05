package pay

import (
	"github.com/Sirupsen/logrus"
	"testing"
)

func guid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

func TestUnifor(t *testing.T) {
	return
	wx := NewWxPay(config)
	order := guid()
	//req, resp, err := wx.UnifiedOrderNative(order, "1303441501", "74682837812332", "http://wx.xlh-tech.com/notify", "支付商品", 100)
	req, resp, err := wx.UnifiedOrderNative(order, "xxxxxxxx", "xxxxxxxx", "http://xxx/notify", "支付商品", 10)
	if err != nil {
		t.Fatal(err)
	}
	logrus.WithFields(logrus.Fields{
		"request":  req,
		"response": resp,
	}).Info("pay request ok")
	t.Logf("ok")
}
