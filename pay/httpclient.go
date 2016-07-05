package pay

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/ilisin/weixin"
)

func (this *WxPay) HttpTlsPost(path string, bodyObj interface{}) ([]byte, error) {
	if this.config == nil {
		return nil, fmt.Errorf("cann't found cert config")
	}
	url := fmt.Sprintf("%v%v", this.config.MchBashUrl, path)
	logrus.WithField("url", url).Debug("http tls post")
	return weixin.HttpTlsPost(url, this.config.Cert.Ca, this.config.Cert.Cert, this.config.Cert.Key, bodyObj)
}
