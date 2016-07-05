package pay

import (
	//	"encoding/json"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ilisin/weixin"
	"net/http"
)

//交易类型
type UnifiedOrderType string

const (
	UO_TT_JSAPI  UnifiedOrderType = "JSAPI"  //公众号支付
	UO_TT_NATIVE UnifiedOrderType = "NATIVE" //原生扫码支付
)

const (
	PayReturnCodeSuccess = "SUCCESS"
	PayReturnCodeFail    = "FAIL"
)

var (
	ErrSignError   = errors.New("sign error")
	ErrMchId       = errors.New("manchert id error")
	ErrXmlMarshal  = errors.New("xml marshal error")
	ErrUniforOrder = errors.New("unified order error")
	ErrProductId   = errors.New("product id error")
)

type WxPay struct {
	config *WxPayConfig
}

func NewWxPay(config *WxPayConfig) *WxPay {
	return &WxPay{config}
}

type (
	//统一订单接口实体
	UnifiedOrderRequest struct {
		WxAppId    string `xml:"appid"`        //公众账号appid
		MchId      string `xml:"mch_id"`       //商户号
		SubMchId   string `xml:"sub_mch_id"`   //子商户号
		OutTradeNo string `xml:"out_trade_no"` //商户订单号
		TotalFee   int    `xml:"total_fee"`    //支付金额
		//		TimeStart   string           `xml:"time_start"`       //订单生成时间
		//		TimeExpire  string           `xml:"time_expire"`      //订单失效时间
		TradeType   UnifiedOrderType `xml:"trade_type"`       //交易类型
		ProductId   string           `xml:"product_id"`       //扫码支付商品ID 二维码中的商品ID
		NotifyUrl   string           `xml:"notify_url"`       //接收微信支付异步通知回调地址
		NonceString string           `xml:"nonce_str"`        //随机字符串
		Sign        string           `xml:"sign"`             //签名
		Body        string           `xml:"body"`             //商品或支付单简要描述
		Ip          string           `xml:"spbill_create_ip"` //APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP
		OpenId      string           `xml:"openid"`
	}

	UnifiedOrderResponse struct {
		ReturnCode string `xml:"return_code"` //SUCCESS/FAIL
		ReturnMsg  string `xml:"return_msg"`

		//以下字段是在return_code为SUCCESS的时候返回
		WxAppId string `xml:"appid"` //公众账号appid
		//SubAppId     string `xml:"sub_appid"`     //公众账号appid
		MchId       string `xml:"mch_id"`    //商户号
		NonceString string `xml:"nonce_str"` //随机字符串
		Sign        string `xml:"sign"`
		ResultCode  string `xml:"result_code"` //业务结果SUCCESS/FAIL
		ErrCode     string `xml:"err_code"`
		ErrCodeDes  string `xml:"err_code_des"`

		//以下字段在return_code和result_code都为SUCCESS的时候返回
		TradeType string `xml:"trade_type"` //交易类型
		PrepayId  string `xml:"prepay_id"`  //预支付交易会话标识
		CodeUrl   string `xml:"code_url"`   //trade_type为NATIVE是有返回，可将该参数值生成二维码展示出来进行扫码支付
	}

	JsApiRequest struct {
		WxAppId     string `json:"appId" xml:"appId"`         //公众账号appid
		TimeStamp   string `json:"timeStamp" xml:"timeStamp"` //时间戳
		NonceString string `json:"nonceStr" xml:"nonceStr"`   //随机字符串
		Package     string `json:"package" xml:"package"`     //统一下单接口返回的prepay_id参数值，提交格式如：prepay_id=***
		SignType    string `json:"signType" xml:"signType"`   //签名方式MD5
		PaySign     string `json:"paySign" xml:"paySign"`     //签名
	}

	PayCallbackResponse struct {
		ReturnCode string `xml:"return_code"` //SUCCESS/FAIL
		ReturnMsg  string `xml:"return_msg"`

		//以下字段是在return_code为SUCCESS的时候返回
		WxAppId       string `xml:"appid"`       //公众账号appid
		MchId         string `xml:"mch_id"`      //商户号
		SubMchId      string `xml:"sub_mch_id"`  //子商户ID
		DeviceInfo    string `xml:"device_info"` //设备号
		NonceString   string `xml:"nonce_str"`   //随机字符串
		Sign          string `xml:"sign"`
		ResultCode    string `xml:"result_code"`    //业务结果SUCCESS/FAIL
		ErrCode       string `xml:"err_code"`       //错误代码
		ErrCodeDes    string `xml:"err_code_des"`   //错误代码
		OpenId        string `xml:"openid"`         //用户标识
		IsSubscribe   string `xml:"is_subscribe"`   //是否关注公众账号
		TradeType     string `xml:"trade_type"`     //交易类型
		BankType      string `xml:"bank_type"`      //付款银行
		TotalFee      int    `xml:"total_fee"`      //订单总金额，单位为分
		FeeType       string `xml:"fee_type"`       //货币种类
		CashFee       int    `xml:"cash_fee"`       //现金支付金额(分)
		CashFeeType   string `xml:"cash_fee_type"`  //现金支付货币类型
		CouponFee     int    `xml:"coupon_fee"`     //代金券或立减优惠金额
		CouponCount   int    `xml:"coupon_count"`   //代金券或立减优惠使用数量
		CouponId      string `xml:"coupon_id_$n"`   //代金券或立减优惠ID
		CouponFeeN    int    `xml:"coupon_fee_$n"`  //代金券或立减优惠金额
		TransactionId string `xml:"transaction_id"` //微信支付订单号
		OutTradeNo    string `xml:"out_trade_no"`   //商户订单号
		Attach        string `xml:"attach"`         //商家数据包
		TimeEnd       string `xml:"time_end"`       //支付完成时间 yyyyMMddHHmmss
	}

	//支付回调后返回给微信的信息
	PayCallbackToWxResponse struct {
		ReturnCode string `xml:"return_code"` //SUCCESS/FAIL
		ReturnMsg  string `xml:"return_msg"`
	}
)

func ErrorMessage(err error) string {
	switch err {
	case ErrSignError:
		return "签名错误"
	case ErrMchId:
		return "商户id错误"
	case ErrXmlMarshal:
		return "xml解析异常"
	case ErrUniforOrder:
		return "统一下单错误"
	case ErrProductId:
		return "product id异常"
	default:
		return "未知错误"
	}
}

func (this *PayCallbackResponse) IsSuccess() bool {
	return this.ReturnCode == PayReturnCodeSuccess
}
func (this *PayCallbackResponse) IsFail() bool {
	return this.ReturnCode == PayReturnCodeFail
}

func (this *WxPay) PayCallbackToWxResponseSuccess() *PayCallbackToWxResponse {
	return &PayCallbackToWxResponse{
		ReturnCode: PayReturnCodeSuccess,
		ReturnMsg:  "OK",
	}
}

func (this *WxPay) PayCallbackToWxResponseFail(msg ...string) *PayCallbackToWxResponse {
	resp := &PayCallbackToWxResponse{
		ReturnCode: PayReturnCodeFail,
		ReturnMsg:  "验证失败",
	}
	if len(msg) > 0 && len(msg[0]) > 0 {
		resp.ReturnMsg = msg[0]
	}
	return resp
}

/*
调用微信统一下单
orderNo 订单号
productId 扫码支付商品ID 二维码中的商品ID
url 接收微信支付异步通知回调地址
ty 支付类型
fee 支付金额（单位分）
*/
func (this *WxPay) unifiedOrder(req *UnifiedOrderRequest) (*UnifiedOrderRequest, *UnifiedOrderResponse, error) {
	req.Sign = this.mchSign(req)
	dat, err := weixin.HttpPost(this.config.MchBashUrl+"/pay/unifiedorder", req, weixin.FILE_TYPE_XML)
	if err != nil {
		return req, nil, err
	}
	resp := UnifiedOrderResponse{}
	logrus.WithField("resp", string(dat)).Debug("get red pack")
	err = xml.Unmarshal(dat, &resp)
	if err != nil {
		return req, nil, err
	}
	//chkSign := this.mchSign(resp)
	//if chkSign != resp.Sign {
	//	return nil, nil, fmt.Errorf("sign验证失败[%v]", chkSign)
	//}
	if resp.ReturnCode == "SUCCESS" && resp.ResultCode == "SUCCESS" {
		return req, &resp, nil
	}
	return req, &resp, fmt.Errorf("[return:%v][result:%v]%v", resp.ReturnCode, resp.ResultCode, resp.ReturnMsg)
}

func (this *WxPay) UnifiedOrderJSAPI(openId, orderNo, url, body string, fee int) (*UnifiedOrderRequest, *UnifiedOrderResponse, error) {
	req := &UnifiedOrderRequest{
		MchId:       this.config.MctId,
		SubMchId:    this.config.SubMctId,
		WxAppId:     this.config.AppId,
		OutTradeNo:  orderNo,
		TotalFee:    fee,
		NonceString: productANonceString(),
		TradeType:   UO_TT_JSAPI,
		ProductId:   "",
		NotifyUrl:   url,
		Body:        body,
		Ip:          "192.168.10.1",
		OpenId:      openId,
	}
	return this.unifiedOrder(req)
}

func (this *WxPay) UnifiedOrderJSAPIParameter(openId, orderNo, url, body string, fee int) (*JsApiRequest, error) {
	req, resp, err := this.UnifiedOrderJSAPI(openId, orderNo, url, body, fee)
	if err != nil {
		return nil, err
	}
	_, p := this.MakeJsApiParameter(req, resp)
	return p, nil
}

func (this *WxPay) MakeJsApiParameter(req *UnifiedOrderRequest, resp *UnifiedOrderResponse) (string, *JsApiRequest) {
	if len(resp.PrepayId) == 0 {
		return "", nil
	}

	p := &JsApiRequest{
		WxAppId:     req.WxAppId,
		TimeStamp:   strconv.FormatInt(time.Now().Unix(), 10),
		NonceString: productANonceString(),
		Package:     `prepay_id=` + resp.PrepayId,
		SignType:    `MD5`,
	}
	p.PaySign = this.mchSign(p)
	return resp.PrepayId, p
}

func (this *WxPay) UnifiedOrderNative(orderNo, subMchId, productId, url, body string, fee int) (*UnifiedOrderRequest, *UnifiedOrderResponse, error) {
	req := &UnifiedOrderRequest{
		MchId:       this.config.MctId,
		SubMchId:    this.config.SubMctId,
		WxAppId:     this.config.AppId,
		OutTradeNo:  orderNo,
		TotalFee:    fee,
		NonceString: productANonceString(),
		TradeType:   UO_TT_NATIVE,
		ProductId:   productId,
		NotifyUrl:   url,
		Body:        body,
		Ip:          "192.168.10.1",
	}
	return this.unifiedOrder(req)
}

func ReadAll(r io.Reader) ([]byte, error) {
	return readAll(r, bytes.MinRead) //const MinRead = 512
}

//
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, capacity))
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		//buf太大会返回相应错误
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(r) //关键就是这个家伙
	return buf.Bytes(), err
}
func (this *WxPay) MakeJsApiResult(body io.ReadCloser) (*PayCallbackResponse, error) {
	data, err := ReadAll(body)
	if err != nil {
		return nil, err
	}
	resp := &PayCallbackResponse{}
	logrus.WithField("resp", string(data)).Debug("MakeJsApiResult")
	err = xml.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}
	chkSign := this.mchSign(resp)
	if chkSign != resp.Sign {
		return resp, fmt.Errorf("sign验证失败[%v]", chkSign)
	}
	return resp, nil
}

// request,response 为注册的http回调请求及相应
// 回调参数传入*PayCallbackResponse,传出为异常状态
func (this *WxPay) ParseJsApiNotifyRequest(request *http.Request, response http.ResponseWriter,
	callback func(*PayCallbackResponse) error) error {
	rep := &PayCallbackResponse{}
	data, err := ReadAll(request.Body)
	if err == nil {
		logrus.WithField("resp", string(data)).Debug("ParseJsApiNotifyRequest")
		err = xml.Unmarshal(data, rep)
	}
	if err != nil {
		err = ErrXmlMarshal
	} else if rep.MchId != this.config.MctId {
		err = ErrMchId
	} else if this.mchSign(rep) != rep.Sign {
		err = ErrSignError
	}
	var resp *PayCallbackToWxResponse
	resp = this.PayCallbackToWxResponseSuccess()
	if err == nil {
		err = callback(rep)
		if err != nil {
			resp = this.PayCallbackToWxResponseFail(err.Error())
		}
	} else {
		resp = this.PayCallbackToWxResponseFail(err.Error())
	}
	bytes, er := xml.Marshal(resp)
	if er != nil {
		return er
	}
	_, er = response.Write(bytes)
	if er != nil {
		return er
	}
	response.WriteHeader(http.StatusOK)
	return err
}
