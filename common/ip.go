package common

import (
	"net"
	"errors"
)

func GetIntranceIp() (string, error) {
	// 返回系统的单播接口地址列表
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查Ip地址判断是否回环地址
		ipnet, ok := address.(*net.IPNet)
		// IsLoopback()判断ip是否是回环地址
		if ok && !ipnet.IP.IsLoopback() {
			// To4()将IPv4地址ip转换为4字节表示。如果ip不是IPv4地址，To4返回nil。
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("读取地址异常")
}