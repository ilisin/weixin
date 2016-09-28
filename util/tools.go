package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
	"net"
)

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func ProductANonceString() string {
	str := fmt.Sprintf("%v", time.Now().Nanosecond())
	return GetMd5String(str)
}


func GetALocalIpAddress() string{
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}