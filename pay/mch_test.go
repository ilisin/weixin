package pay

import (
	"encoding/xml"
	"github.com/Sirupsen/logrus"
	"testing"
)

var config = &WxPayConfig{}

func TestSendPack(t *testing.T) {
	return
	wx := NewWxPay(config)

	//wx.HttpTlsPost("/mmpaymkttransfers/gethbinfo")
	_, _, err := wx.SendRedPack("xxx", "xxx", "oThpRwnQ2MM1sESDcmscE2HFXBOg", 100)
	if err != nil {
		t.Error(err)
	}
}

func TestQueryRedPack(t *testing.T) {
	return
	wx := NewWxPay(config)

	req, resp, err := wx.GetRedPackInfo("xxxxxxxxxx")

	tt := QueryRedPackResponse{}
	tt.HbList.Items = make([]GroupRedPackInfoItem, 3)
	dd, _ := xml.Marshal(tt)
	logrus.WithField("data", string(dd)).Info("marshal")

	if err != nil {
		t.Error(err)
		return
	}

	logrus.WithFields(logrus.Fields{
		"req":  *req,
		"resp": *resp,
	}).Info("get info")
}
