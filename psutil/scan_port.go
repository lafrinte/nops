package psutil

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

type ScanPort struct {
	ip       net.IP
	port     []int
	Timeout  time.Duration
	scanType string
}

func NewScanPort(ipString string, ports []int, timeoutMillSecond int, protocol string) (*ScanPort, error) {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("%s is not valid ip address", ipString)
	}

	protocol = strings.ToLower(protocol)
	switch protocol {
	case TCP, UDP:
	default:
		return nil, fmt.Errorf("optional value: (%s|%s), current %s", TCP, UDP, protocol)
	}

	s := &ScanPort{
		ip:       ip,
		Timeout:  time.Millisecond * 500,
		scanType: protocol,
	}

	if timeoutMillSecond > 0 {
		s.Timeout = time.Duration(timeoutMillSecond) * time.Millisecond
	}

	s.port = s.portFilter(ports)

	return s, nil
}

func (s *ScanPort) portFilter(ps []int) []int {
	i := 0

	// do uniq for number in ps
	uniq := make(map[int]struct{})

	for _, p := range ps {
		if _, ok := uniq[p]; !ok {
			if p >= 1 && p < 65535 {
				uniq[p] = struct{}{}
				ps[i] = p
				i++
			}
		}
	}

	return ps[:i]
}

func (s *ScanPort) scanPort(port int) int {
	stateCode := 0
	conn, err := net.DialTimeout(s.scanType, fmt.Sprintf("%s:%d", s.ip, port), s.Timeout)
	if err != nil {
		stateCode = -1
		return stateCode
	}

	_ = conn.Close()

	return stateCode
}

func (s *ScanPort) Scan() map[int]int {
	var (
		syncMap sync.Map
		wg      sync.WaitGroup
		state   = make(map[int]int)
	)

	for _, port := range s.port {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			syncMap.Store(p, s.scanPort(p))
		}(port)
	}

	wg.Wait()

	syncMap.Range(func(key, value any) bool {
		state[key.(int)] = value.(int)
		return true
	})

	return state
}
