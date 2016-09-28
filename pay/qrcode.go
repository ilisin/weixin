// sence :https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=6_4
package pay

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
	"github.com/ilisin/weixin/util"
)

type (
	//统一订单接口实体
	QrNotifyRequest struct {
		XMLName     xml.Name `xml:"xml"`
		AppId       string   `xml:"appid"` //公众账号appid
		OpenId      string   `xml:"openid"`
		MchId       string   `xml:"mch_id"`       //商户号
		IsSubscribe string   `xml:"is_subscribe"` // Y or N
		NonceString string   `xml:"nonce_str"`    //随机字符串
		ProductId   string   `xml:"product_id"`   //扫码支付商品ID 二维码中的商品ID
		Sign        string   `xml:"sign"`         //签名
	}

	QrNotifyResponse struct {
		XMLName    xml.Name `xml:"xml"`
		ReturnCode string   `xml:"return_code"` //SUCCESS/FAIL
		ReturnMsg  string   `xml:"return_msg"`

		//以下字段是在return_code为SUCCESS的时候返回
		WxAppId     string `xml:"appid"`       //公众账号appid
		MchId       string `xml:"mch_id"`      //商户号
		NonceString string `xml:"nonce_str"`   //随机字符串
		PrepayId    string `xml:"prepay_id"`   //预支付交易会话标识
		ResultCode  string `xml:"result_code"` //业务结果SUCCESS/FAIL
		ErrCodeDes  string `xml:"err_code_des"`
		Sign        string `xml:"sign"`
	}
)

func (this *WxPay) MakeGoodsQrCode(pid string) string {
	model := struct {
		AppId       string `xml:"appid"`  //公众账号appid
		MchId       string `xml:"mch_id"` //商户号
		Timestamp   int64  `xml:"time_stamp"`
		NonceString string `xml:"nonce_str"` //随机字符串
		ProductId   string `xml:"product_id"`
		Sign        string `xml:"sign"`
	}{
		AppId:       this.config.AppId,
		MchId:       this.config.MctId,
		Timestamp:   time.Now().Unix(),
		NonceString: util.ProductANonceString(),
		ProductId:   pid,
	}
	return fmt.Sprintf("weixin://wxpay/bizpayurl?appid=%v&mch_id=%v&nonce_str=%v&product_id=%v&time_stamp=%v&sign=%v",
		model.AppId, model.MchId, model.NonceString, model.ProductId, model.Timestamp, this.mchSign(model))
}

// request,response 为注册的http回调请求及相应
// 回调参数传入openid和product_id,传出为订单号,商品描述，金额，及异常状态
func (this *WxPay) ParseQrNotifyRequest(request *http.Request, response http.ResponseWriter, notifyUrl string,
	callback func(string, string) (string, string, int, error)) error {
	req := &QrNotifyRequest{}
	resp := &QrNotifyResponse{}
	data, err := ReadAll(request.Body)
	if err == nil {
		err = xml.Unmarshal(data, req)
	}
	if err != nil {
		err = ErrXmlMarshal
	} else if req.MchId != this.config.MctId {
		err = ErrMchId
	}
	if this.mchSign(req) != req.Sign {
		// logrus.WithFields(logrus.Fields{
		// 	"reqest_sign": req.Sign,
		// 	"self_sign":   s,
		// }).Error("sign not equerl")
		err = ErrSignError
	}
	if err == nil {
		orderNo, orderDesc, fee, err := callback(req.OpenId, req.ProductId)
		if err != nil {
			err = ErrProductId
		} else {
			_, uniresp, err := this.UnifiedOrderNative(orderNo, "", req.ProductId, notifyUrl, orderDesc, fee)
			if err != nil {
				err = ErrUniforOrder
			} else {
				resp.PrepayId = uniresp.PrepayId
				resp.ResultCode = "SUCCESS"
			}
		}
	}
	if err != nil {
		resp.ReturnCode = "FAIL"
		resp.ReturnMsg = ErrorMessage(err)
	} else {
		resp.ReturnCode = "SUCCESS"
	}
	resp.WxAppId = this.config.AppId
	resp.MchId = this.config.MctId
	resp.NonceString = req.NonceString
	resp.Sign = this.mchSign(resp)
	byts, er := xml.Marshal(resp)
	if er != nil {
		return er
	}
	// logrus.WithField("byts", string(byts)).Debug("notify response")
	_, er = response.Write(byts)
	if er != nil {
		return er
	}
	response.WriteHeader(http.StatusOK)
	return err
}
