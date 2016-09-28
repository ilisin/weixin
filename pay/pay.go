package pay

import (
	//	"encoding/json"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"

	"github.com/Sirupsen/logrus"
	"github.com/ilisin/weixin/util"
	"net/http"
	"reflect"
	"sort"
	"strings"
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

	HttpsClient *util.HttpsClient
	Logger *logrus.Logger

	IpAddress string
}

func NewWxPay(config *WxPayConfig) *WxPay {
	p := &WxPay{
		config:config,
	}
	p.Logger = logrus.New()
	l, err := logrus.ParseLevel(config.LoggerLevel)
	if err != nil {
		logrus.Fatal("unkown weixin mp logger level")
	}
	p.Logger.Level = l
	p.HttpsClient,err = util.NewHttpsClient(config.MchBashUrl,config.Cert.Ca,config.Cert.Cert,config.Cert.Key,l)
	if err != nil {
		logrus.WithField("error",err).Fatal("init weixin pay sdk error")
	}
	p.IpAddress = util.GetALocalIpAddress()
	return p
}

func (this *WxPay) mchSign(model interface{}) string {
	typ := reflect.TypeOf(model)
	val := reflect.ValueOf(model)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	strs := make([]string, 0)
	dic := make(map[string]interface{})
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		valueField := val.Field(i)
		xmlTag := typeField.Tag.Get("xml")
		if len(xmlTag) == 0 {
			continue
		}
		if xmlTag == "sign" {
			continue
		}
		if xmlTag == "xml" {
			continue
		}
		if typeField.Type.Kind() == reflect.Int || typeField.Type.Kind() == reflect.Int64 {
			if valueField.Interface() == 0 {
				continue
			}
		}
		tstr := fmt.Sprintf("%v", valueField.Interface())
		if len(tstr) > 0 {
			strs = append(strs, xmlTag)
			dic[xmlTag] = valueField.Interface()
		}
	}
	sort.Strings(strs)
	pas := make([]string, len(dic))
	for i := 0; i < len(dic); i++ {
		pas[i] = fmt.Sprintf("%v=%v", strs[i], dic[strs[i]])
	}
	tempStr := strings.Join(pas, "&")
	tempStr = fmt.Sprintf("%v&key=%v", tempStr, this.config.ApiKey)
	logrus.WithField("tempStr", tempStr).Debug("sign")
	tempStr = util.GetMd5String(tempStr)
	return strings.ToUpper(tempStr)
}

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

func (this *WxPay) UnifiedOrderNative(orderNo, subMchId, productId, url, body string, fee int) (*UnifiedOrderRequest, *UnifiedOrderResponse, error) {
	req := &UnifiedOrderRequest{
		MchId:       this.config.MctId,
		SubMchId:    this.config.SubMctId,
		WxAppId:     this.config.AppId,
		OutTradeNo:  orderNo,
		TotalFee:    fee,
		NonceString: util.ProductANonceString(),
		TradeType:   UO_TT_NATIVE,
		ProductId:   productId,
		NotifyUrl:   url,
		Body:        body,
		Ip:          this.IpAddress,
	}
	return this.unifiedOrder(req)
}

func ReadAll(r io.Reader) ([]byte, error) {
	return readAll(r, bytes.MinRead) //const MinRead = 512
}

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
