package weixin

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/xml"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	HTTP_REQUEST_TIMEOUT = 20 * time.Second
)

const (
	FILE_TYPE_JSON = `json`
	FILE_TYPE_XML  = `xml`
)

func HttpRequest(method, url string, obj interface{}, ty string) (data []byte, err error) {
	var reader io.Reader = nil
	var dat []byte
	if obj != nil {
		switch ty {
		case FILE_TYPE_JSON:
			dat, err = json.Marshal(obj)
			if err != nil {
				return nil, err
			}
		case FILE_TYPE_XML:
			dat, err = xml.Marshal(obj)
			if err != nil {
				return nil, err
			}
		}
		logrus.WithField("data", string(dat)).Info("post http")
		buffer := bytes.NewBuffer(dat)
		reader = buffer
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: HTTP_REQUEST_TIMEOUT,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func HttpGet(url string, ty string) (data []byte, err error) {
	data, err = HttpRequest("GET", url, nil, ty)
	return
}

func HttpPost(url string, obj interface{}, ty string) (data []byte, err error) {
	data, err = HttpRequest("POST", url, obj, ty)
	return
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
		Timeout:   HTTP_REQUEST_TIMEOUT}
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

func HttpTlsGet(url string, ca, cert, certKey string) (data []byte, err error) {
	data, err = HttpTlsRequest("GET", url, ca, cert, certKey, nil)
	return
}

func HttpTlsPost(url string, ca, cert, certKey string, obj interface{}) (data []byte, err error) {
	data, err = HttpTlsRequest("POST", url, ca, cert, certKey, obj)
	return
}
