package sock

import (
	"context"
	"fmt"
	"github.com/zeromq/goczmq"
	"strings"
	"time"
)

const (
	DefaultRetryAttempts          = uint8(5)
	DefaultRetryInterval          = time.Second
	DefaultExitWaitTimeout        = time.Second * 30
	DefaultMaxBufferSize          = 10000
	DefaultTcpKeepAliveCnt        = -1
	DefaultTcpKeepAliveIdleSec    = -1
	DefaultHeartbeatTTLSec        = 30 // unit: millisecond
	DefaultHeartbeatTimeoutSec    = 5  // unit: millisecond
	DefaultHeartbeatIvlSec        = 15 // unit: millisecond
	DefaultReconnectIvlMillSec    = 100
	DefaultReconnectIvlMaxMillSec = 0
	DefaultSndhwm                 = 10000
	DefaultSendTimeoutSec         = 0
	DefaultRecvTimeoutSec         = 0
)

type Option func(s *Sock)

func WithSndhwm(val int) Option {
	return func(s *Sock) {
		s.Sndhwm = val
	}
}

func WithType(val string) Option {
	return func(s *Sock) {
		switch strings.ToUpper(val) {
		case "PUB":
			s.Type = goczmq.Pub
		case "SUB":
			s.Type = goczmq.Sub

		// enable in next pr
		//case "ROUTER":
		//	s.Type = goczmq.Router
		//case "DEALER":
		//	s.Type = goczmq.Dealer
		case "PUSH":
			s.Type = goczmq.Push
		case "PULL":
			s.Type = goczmq.Pull
		case "REQ":
			s.Type = goczmq.Req
		case "REP":
			s.Type = goczmq.Rep
		default:
			msg := fmt.Sprintf("wrone socket mode: %s", strings.ToUpper(val))
			panic(fmt.Errorf(msg))
		}
	}
}

func WithCtx(val context.Context) Option {
	return func(s *Sock) {
		s.ctx = val
	}
}

func WithEndpoint(val string) Option {
	return func(s *Sock) {
		s.Endpoint = val
	}
}

func WithAttach() Option {
	return func(s *Sock) {
		s.attach = true
	}
}

func DisableRestart() Option {
	return func(s *Sock) {
		s.DisableRestart.Set(true)
	}
}

func WithExitWaitTimeout(val time.Duration) Option {
	return func(s *Sock) {
		s.ExitWaitTimeout = val
	}
}

func WithRetryInterval(val time.Duration) Option {
	return func(s *Sock) {
		s.RetryInterval = val
	}
}

func WithRetryAttempts(val int) Option {
	return func(s *Sock) {
		s.RetryAttempts = uint8(val)
	}
}

func WithMaxBufferSize(val int) Option {
	return func(s *Sock) {
		s.MaxBufferSize = val
	}
}

func EnableTcpKeepAlive() Option {
	return func(s *Sock) {
		s.EnableTcpKeepAlive = true
	}
}

func WithTcpKeepAliveIdleSec(val int) Option {
	return func(s *Sock) {
		s.TcpKeepAliveIdleSec = int16(val)
	}
}

func WithTcpKeepAliveCnt(val int) Option {
	return func(s *Sock) {
		s.TcpKeepAliveCnt = int8(val)
	}
}

func WithHeartbeatIvlSec(val int) Option {
	return func(s *Sock) {
		s.HeartbeatIvlSec = uint16(val)
	}
}

func WithHeartbeatTimoutSec(val int) Option {
	return func(s *Sock) {
		s.HeartbeatTimoutSec = uint16(val)
	}
}

func WithHeartbeatTTLSec(val int) Option {
	return func(s *Sock) {
		s.HeartbeatTTLSec = uint16(val)
	}
}

func WithSendTimeoutSec(val int) Option {
	return func(s *Sock) {
		s.SendTimeoutSec = uint16(val)
	}
}

func WithRecvTimeoutSec(val int) Option {
	return func(s *Sock) {
		s.RecvTimeoutSec = uint16(val)
	}
}

func WithInChannel(val chan []byte) Option {
	return func(s *Sock) {
		s.in = val
	}
}

func WithOutChannel(val chan []byte) Option {
	return func(s *Sock) {
		s.out = val
	}
}
