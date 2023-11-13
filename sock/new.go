package sock

import (
	"context"
	"github.com/lafrinte/nops/str"
	"github.com/zeromq/goczmq"
)

func New(opts ...Option) *Sock {
	soc := &Sock{
		id:                     str.ID().String(),
		RetryAttempts:          DefaultRetryAttempts,
		RetryInterval:          DefaultRetryInterval,
		MaxBufferSize:          DefaultMaxBufferSize,
		ExitWaitTimeout:        DefaultExitWaitTimeout,
		EnableTcpKeepAlive:     false,
		TcpKeepAliveCnt:        DefaultTcpKeepAliveCnt,
		TcpKeepAliveIdleSec:    DefaultTcpKeepAliveIdleSec,
		HeartbeatIvlSec:        DefaultHeartbeatIvlSec,
		HeartbeatTimoutSec:     DefaultHeartbeatTimeoutSec,
		HeartbeatTTLSec:        DefaultHeartbeatTTLSec,
		ReconnectIvlMillSec:    DefaultReconnectIvlMillSec,
		ReconnectIvlMaxMillSec: DefaultReconnectIvlMaxMillSec,
		Sndhwm:                 DefaultSndhwm,
		dropMsgCount:           0,
		sendMsgCount:           0,
		RecvTimeoutSec:         DefaultRecvTimeoutSec,
		SendTimeoutSec:         DefaultSendTimeoutSec,
	}

	for _, opt := range opts {
		opt(soc)
	}

	if soc.ctx == nil {
		soc.ctx = context.Background()
	}

	switch soc.Type {
	case goczmq.Pub, goczmq.Push:
		if soc.in == nil {
			soc.in = make(chan []byte, soc.MaxBufferSize)
		}
		soc.out = make(chan []byte, 0)
		soc.retryCh = make(chan *RetryMsg, soc.MaxBufferSize)
	case goczmq.Sub, goczmq.Pull:
		if soc.out == nil {
			soc.out = make(chan []byte, soc.MaxBufferSize)
		}

		soc.in = make(chan []byte, 0)
		soc.retryCh = make(chan *RetryMsg, 0)
	case goczmq.Req, goczmq.Rep:
		soc.in = make(chan []byte, 1)
		soc.out = make(chan []byte, 1)
		soc.retryCh = make(chan *RetryMsg, 0)

		if soc.SendTimeoutSec == 0 {
			soc.SendTimeoutSec = 3
		}

		if soc.RecvTimeoutSec == 0 {
			soc.RecvTimeoutSec = 3
		}
	}

	return soc
}
