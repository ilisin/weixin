package util

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
	"errors"

	"github.com/Sirupsen/logrus"
	"fmt"
	"strings"
)

var (
	ErrUnkownRequestEncoding = errors.New("unkown request encoding type")
	ErrUnkownResponseEncoding = errors.New("unkown response encoding type")
	ErrEncodingRequest = errors.New("encoding request param error")
	ErrEncodingResponse = errors.New("encoding response param error")
)

type EncodingType string

const (
	EncodingTypeJson EncodingType = "json"
	EncodingTypeXml EncodingType = "xml"
)

type HttpRequestMethod string

const (
	GET HttpRequestMethod = "GET"
	POST HttpRequestMethod = "POST"
)

type HttpClient struct {
	BaseUrl string

	http.Client
	logger  *logrus.Logger
}

func NewClient(baseUrl string, l logrus.Level) *HttpClient {
	c := &HttpClient{
		BaseUrl:baseUrl,
		Client:http.Client{
			Timeout:20 * time.Second,
		},
		logger:logrus.New(),
	}
	c.logger.Level = l
	return c
}

type HttpsClient struct {
	HttpClient

	Ca      string
	Cert    string
	CertKey string
}

func NewHttpsClient(baseUrl, ca, cert, certKey string, l logrus.Level) (*HttpsClient, error) {
	c := &HttpsClient{
		HttpClient:HttpClient{
			BaseUrl:baseUrl,
			logger:logrus.New(),
		},
		Ca:ca,
		Cert:cert,
		CertKey:certKey,
	}
	c.logger.Level = l

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
			RootCAs:      pool,
			Certificates: []tls.Certificate{clicrt},
		},
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	c.Client = http.Client{
		Transport: tr,
		Timeout: 20 * time.Second,
	}
	return c, nil
}

func (hc HttpClient)Request(method HttpRequestMethod, path string, param interface{}, requestType EncodingType, responseType EncodingType, respObj interface{}) error {
	var reader io.Reader = nil
	var dat []byte
	var err error
	if param != nil {
		if requestType == EncodingTypeJson {
			dat, err = json.Marshal(param)
			if err != nil {
				hc.logger.WithField("err", err).Error("http request json encoding param")
				return ErrEncodingRequest
			}
		} else if requestType == EncodingTypeXml {
			dat, err = xml.Marshal(param)
			if err != nil {
				hc.logger.WithField("err", err).Error("http request xml encoding param error")
				return ErrEncodingRequest
			}
		} else {
			return ErrUnkownRequestEncoding
		}
		buffer := bytes.NewBuffer(dat)
		reader = buffer
	}
	req, err := http.NewRequest(string(method), hc.BaseUrl + path, reader)
	if err != nil {
		return err
	}
	hc.logger.WithFields(logrus.Fields{
		"method":method,
		"url":req.URL,
	}).Debug("http request")
	resp, err := hc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if (respObj != nil) {
		if (responseType == EncodingTypeJson) {
			err = json.NewDecoder(resp.Body).Decode(respObj)
			if err != nil {
				hc.logger.WithField("err", err).Error("http response json encoding error")
			}
		} else if (responseType == EncodingTypeXml) {
			err = xml.NewDecoder(resp.Body).Decode(respObj)
			if err != nil {
				hc.logger.WithField("err", err).Error("http response xml encoding error")
			}
		} else {
			return ErrUnkownResponseEncoding
		}
	}
	return nil
}

func (hc HttpClient) HttpGetJson(path string, params map[string]interface{}, respObj interface{}) error {
	if params != nil {
		pas := make([]string, len(params))
		i := 0
		for k, v := range params {
			pas[i] = fmt.Sprintf("%v=%v", k, v)
			i++
		}
		path = fmt.Sprintf("%v?%v", path, strings.Join(pas, "&"))
	}
	return hc.Request(GET, path, nil, EncodingTypeJson, EncodingTypeJson, respObj)
}

func (hc HttpClient) HttpPostJson(path string, params map[string]interface{}, postData interface{},respObj interface{})error{
	if params != nil {
		pas := make([]string, len(params))
		i := 0
		for k, v := range params {
			pas[i] = fmt.Sprintf("%v=%v", k, v)
			i++
		}
		path = fmt.Sprintf("%v?%v", path, strings.Join(pas, "&"))
	}
	return hc.Request(POST,path,postData,EncodingTypeJson,EncodingTypeJson,respObj)
}

func (hc HttpClient) HttpPostXml(path string,params map[string]interface{},postData interface{},respObj interface{})error{
	if params != nil {
		pas := make([]string, len(params))
		i := 0
		for k, v := range params {
			pas[i] = fmt.Sprintf("%v=%v", k, v)
			i++
		}
		path = fmt.Sprintf("%v?%v", path, strings.Join(pas, "&"))
	}
	return hc.Request(POST,path,postData,EncodingTypeXml,EncodingTypeXml,respObj)
}

//func HttpGet(url string, ty string) (data []byte, err error) {
//	data, err = HttpRequest("GET", url, nil, ty)
//	return
//}
//
//func HttpPost(url string, obj interface{}, ty string) (data []byte, err error) {
//	data, err = HttpRequest("POST", url, obj, ty)
//	return
//}

////tls request
//func HttpTlsRequest(method, url, ca, cert, certKey string, obj interface{}) (data []byte, err error) {
//	//logrus.WithFields(logrus.Fields{
//	//"ca":      ca,
//	//"cert":    cert,
//	//"certKey": certKey,
//	//}).Info("http tls request")
//	var reader io.Reader = nil
//	pool := x509.NewCertPool()
//	caCrt, err := ioutil.ReadFile(ca)
//	if err != nil {
//		return nil, err
//	}
//	pool.AppendCertsFromPEM(caCrt)
//
//	clicrt, err := tls.LoadX509KeyPair(cert, certKey)
//	if err != nil {
//		return nil, err
//	}
//
//	tr := &http.Transport{
//		TLSClientConfig: &tls.Config{
//			//MaxVersion:   tls.VersionSSL30,
//			RootCAs:      pool,
//			Certificates: []tls.Certificate{clicrt},
//			//InsecureSkipVerify: true,
//		},
//		Proxy: http.ProxyFromEnvironment,
//		Dial: (&net.Dialer{
//			Timeout:   30 * time.Second,
//			KeepAlive: 30 * time.Second,
//		}).Dial,
//		TLSHandshakeTimeout: 10 * time.Second,
//	}
//
//	client := &http.Client{
//		Transport: tr,
//		Timeout:   HTTP_REQUEST_TIMEOUT}
//	if obj != nil {
//		dat, err := xml.Marshal(obj)
//		if err != nil {
//			return nil, err
//		}
//		buffer := bytes.NewBuffer(dat)
//		reader = buffer
//	}
//	req, err := http.NewRequest(method, url, reader)
//	if err != nil {
//		return nil, err
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//	return ioutil.ReadAll(resp.Body)
//}

//func HttpTlsGet(url string, ca, cert, certKey string) (data []byte, err error) {
//	data, err = HttpTlsRequest("GET", url, ca, cert, certKey, nil)
//	return
//}
//
//func HttpTlsPost(url string, ca, cert, certKey string, obj interface{}) (data []byte, err error) {
//	data, err = HttpTlsRequest("POST", url, ca, cert, certKey, obj)
//	return
//}
