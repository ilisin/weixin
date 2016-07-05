package mp

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ilisin/weixin"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

//params is the part of url
func (this *WxMp) HttpGet(path string, withToken bool, params map[string]interface{}) (data []byte, err error) {
	url := fmt.Sprintf("%v%v", this.config.ServiceURL, path)
	if withToken {
		if params == nil {
			params = make(map[string]interface{})
		}
		if token, err := this.Token(); err == nil {
			params["access_token"] = token.AccessToken
		} else {
			return nil, err
		}
	}
	if params != nil {
		pas := make([]string, len(params))
		i := 0
		for k, v := range params {
			pas[i] = fmt.Sprintf("%v=%v", k, v)
			i++
		}
		url = fmt.Sprintf("%v?%v", url, strings.Join(pas, "&"))
	}
	logrus.WithField("url", url).Debug("http get")
	data, err = weixin.HttpGet(url, weixin.FILE_TYPE_JSON)
	resp := Response{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		logrus.Error(`HTTPGet 返回错误`, string(data), err)
		return nil, err
	}
	logrus.WithField("resp", resp).Debug("http get")
	if resp.ErrCode != WXResultSeccuse {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"resp": resp,
		}).Info("HTTPGet")
		if resp.ErrCode == WXResultNonMenu {
			return nil, nil
		}
		return nil, fmt.Errorf("[%v]%v", resp.ErrCode, resp.ErrMsg)
	}
	return data, nil
}

//params is the part of path params
//bodyObj is the request data,with json format
func (this *WxMp) HttpPost(path string, withToken bool, params map[string]interface{}, bodyObj interface{}) ([]byte, error) {
	url := fmt.Sprintf("%v%v", this.config.ServiceURL, path)
	if withToken {
		if params == nil {
			params = make(map[string]interface{})
		}
		if token, err := this.Token(); err == nil {
			params["access_token"] = token.AccessToken
		} else {
			return nil, err
		}
	}
	if params != nil {
		pas := make([]string, len(params))
		i := 0
		for k, v := range params {
			pas[i] = fmt.Sprintf("%v=%v", k, v)
			i++
		}
		url = fmt.Sprintf("%v?%v", url, strings.Join(pas, "&"))
	}
	logrus.WithField("url", url).Debug("http post")
	data, err := weixin.HttpPost(url, bodyObj, weixin.FILE_TYPE_JSON)
	resp := Response{}
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != WXResultSeccuse {
		return nil, fmt.Errorf("[%v]%v", resp.ErrCode, resp.ErrMsg)
	}
	return data, nil
}

//tls request
func HttpTlsRequest(method, url, ca, cert, certKey string, obj interface{}) (data []byte, err error) {
	//logrus.WithFields(logrus.Fields{
	//"ca":      ca,
	//"cert":    cert,
	//"certKey": certKey,
	//}).Info("http tls request")
	var reader io.Reader = nil
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, err
	}
	pool.AppendCertsFromPEM(caCrt)

	clicrt, err := tls.LoadX509KeyPair(cert, certKey)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			//MaxVersion:   tls.VersionSSL30,
			RootCAs:      pool,
			Certificates: []tls.Certificate{clicrt},
			//InsecureSkipVerify: true,
		},
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   weixin.HTTP_REQUEST_TIMEOUT}
	if obj != nil {
		dat, err := xml.Marshal(obj)
		if err != nil {
			return nil, err
		}
		buffer := bytes.NewBuffer(dat)
		reader = buffer
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
