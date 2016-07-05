package pay

import (
	"encoding/xml"
	"fmt"
	"github.com/Sirupsen/logrus"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"time"
)

type MchTlsConfig struct {
	//cert
	CaCrtFile string `json:"cert_ca"`
	CertFile  string `json:"cert_file"`
	CertKey   string `json:"cert_key"`

	MchId    string `json:"merchant_id"`
	SubMchId string `json:"sub_merchant_id"`
	SubAppId string `json:"sub_app_id"`
	MchName  string `json:"merchant_name"`

	ApiKey string `json:"api_key"`
}

type SendRedPackRequest struct {
	XMLName     xml.Name `xml:"xml"`
	NonceString string   `xml:"nonce_str"` //随机字符串
	Sign        string   `xml:"sign"`
	MchBillNo   string   `xml:"mch_billno"`   //商户订单号
	MchId       string   `xml:"mch_id"`       //商户号
	SubMchId    string   `xml:"sub_mch_id"`   //子商户号
	WxAppId     string   `xml:"wxappid"`      //公众账号appid
	MsgAppId    string   `xml:"msgappid"`     //触达用户appid,特约商户appid
	SendName    string   `xml:"send_name"`    //商户名称
	OpenId      string   `xml:"re_openid"`    //用户OpenId
	TotalAmount int      `xml:"total_amount"` //付款金额
	TotalNum    int      `xml:"total_num"`    //红包发放总人数
	WiShing     string   `xml:"wishing"`      //红包祝福语
	ClientIp    string   `xml:"client_ip"`    //调用接口的机器ip地址
	ActName     string   `xml:"act_name"`     //活动名称
	Remark      string   `xml:"remark"`       //备注
}

type SendRedPackResponse struct {
	XMLName xml.Name `xml:"xml"`

	ReturnCode string `xml:"return_code"` //SUCCESS or FAIL
	ReturnMsg  string `xml:"return_msg"`

	//以下字段在return_code为SUCCESS时返回
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`  //业务结果SUCCESS/FAIL
	ErrCode    string `xml:"err_code"`     //错误码信息
	ErrCodeDes string `xml:"err_code_des"` //错误代码描述

	//以下字段在return_code和result_code都为SUCCESS的时候有返回
	MchBillNo   string `xml:"mch_billno"`   //商户订单号
	MchId       string `xml:"mch_id"`       //商户号
	WxAppId     string `xml:"wxappid"`      //公众账号appid
	OpenId      string `xml:"re_openid"`    //用户openid
	TotalAmount int    `xml:"total_amount"` //付款金额
	SendTime    int64  `xml:"send_time"`    //红包发送事件
	SendListId  string `xml:"send_listid"`  //微信单号
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

func (this *WxPay) productAMchBillId() string {
	tm := time.Now()
	r := rand.New(rand.NewSource(tm.Unix()))
	return fmt.Sprintf("%v%v%010d", this.config.MctId, tm.Format("20060102"), r.Int31())
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
	tempStr = GetMd5String(tempStr)
	return strings.ToUpper(tempStr)
}

func (this *WxPay) mchPaySign(model interface{}) string {
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
	tempStr = GetMd5String(tempStr)
	return strings.ToUpper(tempStr)
}

//发送红包
//amount单位未分
//返回红包订单号 商户号+yyyymmdd+10位随机数
func (this *WxPay) SendRedPack(amount int, subMchId, subAppId, openId string) (*SendRedPackRequest, *SendRedPackResponse, error) {
	reqModel := &SendRedPackRequest{
		NonceString: productANonceString(),
		MchBillNo:   this.productAMchBillId(),
		MchId:       this.config.MctId,
		SubMchId:    subMchId,
		WxAppId:     this.config.AppId,
		MsgAppId:    subAppId,
		SendName:    this.config.MctName,
		OpenId:      openId,
		TotalAmount: amount,
		TotalNum:    1,
		WiShing:     "祝您开启魔幻之旅",
		ClientIp:    "192.168.11.11",
		ActName:     "魔力红包天天抢",
		Remark:      "越抢越快乐！",
	}
	reqModel.Sign = this.mchSign(reqModel)
	dat, err := this.HttpTlsPost("/mmpaymkttransfers/sendredpack", reqModel)
	if err != nil {
		return reqModel, nil, err
	}
	logrus.WithField("resp", string(dat)).Debug("get red pack")
	resp := SendRedPackResponse{}
	err = xml.Unmarshal(dat, &resp)
	if err != nil {
		return reqModel, nil, err
	}
	//chkSign := this.mchSign(resp)
	//if resp.ReturnCode == "SUCCESS" && chkSign != resp.Sign {
	//return reqModel, &resp, fmt.Errorf("sign验证失败[%v]", chkSign)
	//}
	if resp.ReturnCode == "SUCCESS" && resp.ResultCode == "SUCCESS" {
		return reqModel, &resp, nil
	}
	return reqModel, &resp, fmt.Errorf("[return:%v][result:%v]%v", resp.ReturnCode, resp.ResultCode, resp.ReturnMsg)
}

type QueryRedPackRequest struct {
	XMLName     xml.Name `xml:"xml"`
	NonceString string   `xml:"nonce_str"` //随机字符串
	Sign        string   `xml:"sign"`
	MchBillNo   string   `xml:"mch_billno"` //商户订单号
	MchId       string   `xml:"mch_id"`     //商户号
	WxAppId     string   `xml:"appid"`      //公众账号appid
	BillType    string   `xml:"bill_type"`  //MCHT
}

//红包状态
type RedPackStatus string

const (
	RPStatusSending  RedPackStatus = "SENDING"  //发送中
	RPStatusSend     RedPackStatus = "SENT"     //已发送待领取
	RPStatusFailed   RedPackStatus = "FAILED"   //发送失败
	RPStatusReceived RedPackStatus = "RECEIVED" //已领取
	RPStatusRefund   RedPackStatus = "REFUND"   //已退款
)

//发送方式
type RedPackSendType string

const (
	RPSendTypeApi      RedPackSendType = "API"      //通过api接口发送
	RPSendTypeUpload   RedPackSendType = "UPLOAD"   //通过上传文件方式发放
	RPSendTypeActivity RedPackSendType = "ACTIVITY" //通过活动方式发放
)

//红包类型
type RedPackType string

const (
	RPTypeNormal RedPackType = "NORMAL" //普通红包
	RPTypeGroup  RedPackType = "GROUP"  //裂变红包
)

type GroupRedPackInfoItem struct {
	XMLName xml.Name      `xml:"hbinfo"`
	OpenId  string        `xml:"openid"`
	Status  RedPackStatus `xml:"status"`
	Amount  int           `xml:"amount"`
	RcvTime string        `xml:"rcv_time"` //yyyy-MM-dd HH:mm:ss
}

type GroupRedPackInfo struct {
	XMLName xml.Name `xml:"hblist"`
	Items   []GroupRedPackInfoItem
}

type QueryRedPackResponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"` //SUCCESS/FAIL
	ReturnMsg  string   `xml:"return_msg"`

	//以下字段是在return_code为SUCCESS的时候返回
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"` //业务结果SUCCESS/FAIL
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`

	//以下字段在return_code和result_code都为SUCCESS的时候返回
	MchBillNo    string           `xml:"mch_billno"` //商户订单号
	MchId        string           `xml:"mch_id"`     //商户号
	AppId        string           `xml:"appid"`
	OpenId       string           `xml:"openid"`
	DetailId     string           `xml:"detail_id"`     //使用api发送红包时返回的红包单号
	Status       RedPackStatus    `xml:"status"`        //红包状态
	SendType     RedPackSendType  `xml:"send_type"`     //发送方式
	Type         RedPackType      `xml:"hb_type"`       //红包类型
	TotalNum     int              `xml:"total_num"`     //红包个数
	TotalAmount  int              `xml:"total_amount"`  //总额
	Reason       string           `xml:"reason"`        //发送失败原因
	SendTime     string           `xml:"send_time"`     //发送时间yyyy-MM-dd HH:mm:ss
	RefundTime   string           `xml:"refund_time"`   //退款时间yyyy-MM-dd HH:mm:ss
	RefundAmount int              `xml:"refund_amount"` //退款金额
	Wishing      string           `xml:"wishing"`       //祝福语
	Remark       string           `xml:"remark"`        //活动描述
	ActName      string           `xml:"act_name"`      //活动名称
	HbList       GroupRedPackInfo `xml:"hblist"`        //裂变红包领取列表
}

//查询红包
func (this *WxPay) GetRedPackInfo(redpackBillNo string) (*QueryRedPackRequest, *QueryRedPackResponse, error) {
	reqModel := &QueryRedPackRequest{
		NonceString: productANonceString(),
		MchBillNo:   redpackBillNo,
		MchId:       this.config.MctId,
		WxAppId:     this.config.AppId,
		BillType:    "MCHT",
	}
	reqModel.Sign = this.mchSign(reqModel)
	dat, err := this.HttpTlsPost("/mmpaymkttransfers/gethbinfo", reqModel)
	if err != nil {
		return reqModel, nil, err
	}
	resp := QueryRedPackResponse{}
	logrus.WithField("resp", string(dat)).Debug("get red pack")
	err = xml.Unmarshal(dat, &resp)
	if err != nil {
		return reqModel, nil, err
	}
	//chkSign := this.mchSign(resp)
	//if chkSign != resp.Sign {
	//return reqModel, &resp, fmt.Errorf("sign验证失败[%v]", chkSign)
	//}
	if resp.ReturnCode == "SUCCESS" && resp.ResultCode == "SUCCESS" {
		return reqModel, &resp, nil
	}
	return reqModel, &resp, fmt.Errorf("[return:%v][result:%v]%v", resp.ReturnCode, resp.ResultCode, resp.ReturnMsg)
}
