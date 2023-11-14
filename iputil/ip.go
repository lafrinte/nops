package iputil

import (
	"net"
	"strings"
	"sync"
)

var (
	staticIP string
	once     sync.Once
)

// GetIP return one public IP
func GetIP() string {
	once.Do(func() {
		conn, err := net.Dial("udp", "8.8.8.8:53")
		if err == nil {
			localAddr := conn.LocalAddr().(*net.UDPAddr)
			staticIP = strings.Split(localAddr.String(), ":")[0]
		}
	})

	return staticIP
}

func GetInnerIP() string {
	infs, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, inf := range infs {
		if isEthDown(inf.Flags) || isLoopback(inf.Flags) {
			continue
		}

		addrs, err := inf.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					return ipNet.IP.String()
				}
			}
		}
	}

	return ""
}

func isEthDown(f net.Flags) bool {
	return f&net.FlagUp != net.FlagUp
}

func isLoopback(f net.Flags) bool {
	return f&net.FlagLoopback == net.FlagLoopback
}
