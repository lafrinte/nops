package sock

import (
	"context"
	"fmt"
	"github.com/lafrinte/nops/atom"
	"github.com/lafrinte/nops/iputil"
	"github.com/lafrinte/nops/timer"
	"github.com/zeromq/goczmq"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"runtime"
	"time"
)

type RetryMsg struct {
	retry    uint8
	Msg      []byte
	MaxRetry uint8
}

func (r *RetryMsg) Retry() bool {
	return r.MaxRetry > r.retry
}

func (r *RetryMsg) IterRetryTimes() {
	r.retry++
}

func NewRetryMsg(msg []byte, retry uint8) *RetryMsg {
	if retry == 0 {
		retry = DefaultRetryAttempts
	}

	return &RetryMsg{Msg: msg, retry: 0, MaxRetry: retry}
}

type Sock struct {
	id  string
	ctx context.Context
	soc *goczmq.Sock

	// channel for msg and signal
	in            chan []byte
	out           chan []byte
	MaxBufferSize int

	// socket connection args
	Type     int
	Endpoint string
	attach   bool
	Sndhwm   int

	// reconnect args
	DisableRestart         atom.AtomicBool
	ReconnectIvlMillSec    int
	ReconnectIvlMaxMillSec int

	// heartbeat args
	HeartbeatIvlSec    uint16
	HeartbeatTimoutSec uint16
	HeartbeatTTLSec    uint16

	// tcp args
	EnableTcpKeepAlive  bool
	TcpKeepAliveIdleSec int16
	TcpKeepAliveCnt     int8

	// msg retry args
	retryCh       chan *RetryMsg
	RetryInterval time.Duration
	RetryAttempts uint8

	// safe destroy args
	ExitWaitTimeout time.Duration

	// register endpoint
	ServiceRegisterEndpoint string
}

// bind binds socket on endpoint
func (s *Sock) bind() (*goczmq.Sock, error) {
	soc := goczmq.NewSock(s.Type)
	if _, err := soc.Bind(s.Endpoint); err != nil {
		return nil, err
	}

	soc.SetSndhwm(s.Sndhwm)

	return s.setOptions(soc), nil
}

// connect makes connection to endpoint
func (s *Sock) connect() (*goczmq.Sock, error) {
	soc := goczmq.NewSock(s.Type)
	if err := soc.Connect(s.Endpoint); err != nil {
		return nil, err
	}

	soc.SetSndhwm(s.Sndhwm)

	return s.setOptions(soc), nil
}

// setOptions sets socket tcp and heartbeat option
func (s *Sock) setOptions(soc *goczmq.Sock) *goczmq.Sock {
	if soc.GetType() == goczmq.Sub {
		soc.SetSubscribe("")
	}

	if s.EnableTcpKeepAlive {
		soc.SetTcpKeepalive(1)
		soc.SetTcpKeepaliveIdle(int(s.TcpKeepAliveIdleSec * 1000))
		soc.SetTcpKeepaliveCnt(int(s.TcpKeepAliveCnt))
	}

	if s.HeartbeatTimoutSec > 0 {
		soc.SetHeartbeatTimeout(int(s.HeartbeatTimoutSec) * 1000)
	}

	if s.HeartbeatIvlSec > 0 {
		soc.SetHeartbeatIvl(int(s.HeartbeatIvlSec) * 1000)
	}

	if s.HeartbeatTTLSec > 0 {
		soc.SetHeartbeatTtl(int(s.HeartbeatTTLSec) * 1000)
	}

	soc.SetReconnectIvl(s.ReconnectIvlMillSec)

	if s.ReconnectIvlMaxMillSec > 0 {
		soc.SetReconnectIvlMax(s.ReconnectIvlMaxMillSec)
	}

	s.soc = soc

	return soc
}

/*
Attach attaches a socket to zero or more endpoints.

	call bind when attach is equal to false
	call connect when attach is equal to true
*/
func (s *Sock) Attach() (*goczmq.Sock, error) {
	if s.attach {
		return s.connect()
	}

	return s.bind()
}

func (s *Sock) EmptyBuffer() bool {
	return (len(s.retryCh) + len(s.in) + len(s.out)) == 0
}

// GetInCount gets msg count int channel 'in'
func (s *Sock) GetInCount() int {
	return len(s.in)
}

// GetOutCount gets msg count int channel 'out'
func (s *Sock) GetOutCount() int {
	return len(s.out)
}

// GetRetryCount gets msg count int channel 'retryCh'
func (s *Sock) GetRetryCount() int {
	return len(s.retryCh)
}

// GetBufferInfoMsg return msg count in every channel on sock in string
func (s *Sock) GetBufferInfoMsg() string {
	return fmt.Sprintf("in: [%d], out: [%d], retry: [%d]", s.GetInCount(), s.GetOutCount(), s.GetRetryCount())
}

// GetBufferInfoDict return msg count in every channel on sock in dict
func (s *Sock) GetBufferInfoDict() map[string]int {
	return map[string]int{
		"in":    s.GetInCount(),
		"out":   s.GetOutCount(),
		"retry": s.GetRetryCount(),
	}
}

// IsAutoRestart gets the value of DisableRestart
func (s *Sock) IsAutoRestart() bool {
	return s.DisableRestart.True()
}

// StopAutoRestart sets DisableRestart to true
func (s *Sock) StopAutoRestart() {
	s.DisableRestart.Set(true)
}

// GetInChannel returns in channel
func (s *Sock) GetInChannel() chan []byte {
	return s.in
}

// GetOutChannel returns out channel
func (s *Sock) GetOutChannel() chan []byte {
	return s.out
}

func (s *Sock) CloseCh() {
	close(s.in)
	close(s.out)
	close(s.retryCh)
}

// CloseSock calls Destroy and set soc to nil
func (s *Sock) CloseSock() {
	if s.soc != nil {
		s.soc.Destroy()
		s.soc = nil
	}
}

// sendFrame tries sending msg and puts msg into retry channel when any error occured
func (s *Sock) sendFrame(sock *goczmq.Sock, msg []byte, retry bool) {
	/* goczmq.Sock will panic in C while receiving or sending frame on a destroyed socket,
	   wrapped in thread will transfer thread-panic to main thread and be caught
	*/
	Pool.CtxGo(s.ctx, func() {
		if err := sock.SendFrame(msg, goczmq.FlagNone); err != nil && retry {
			s.retryCh <- NewRetryMsg(msg, s.RetryAttempts)
		}
	})
}

// recvFrame tries receiving msg and puts into out channel
func (s *Sock) recvFrame(sock *goczmq.Sock) {
	/* goczmq.Sock will panic in C while receiving or sending frame on a destroyed socket,
	   wrapped in thread will transfer thread-panic to main thread and be caught
	*/
	Pool.CtxGo(s.ctx, func() {
		for {
			if buf, _, err := sock.RecvFrame(); err == nil {
				s.out <- buf
			} else {
				// prevent call recv for send after socket be destroyed
				if err == goczmq.ErrRecvFrameAfterDestroy {
					panic(err)
				}
			}
		}
	})
}

// retry tries to publish msg in retry channel and re-puts into retry channel when any error occured
func (s *Sock) retry(sock *goczmq.Sock, msg *RetryMsg) {
	if msg.Retry() {

		Pool.CtxGo(s.ctx, func() {
			if err := sock.SendFrame(msg.Msg, goczmq.FlagNone); err != nil {
				msg.IterRetryTimes()

				t := timer.AcquireTimer(s.RetryInterval)
				if !t.Stop() {
					<-t.C
				}

				timer.ReleaseTimer(t)

				s.retryCh <- msg
			}
		})
	}
}

// Release tries release socket after all buffers be triggered
func (s *Sock) Release() error {
	s.StopAutoRestart()

	if s.EmptyBuffer() {
		s.CloseSock()
		return nil
	}

	exitWaitTimeout := time.After(s.ExitWaitTimeout)
	for {
		if s.EmptyBuffer() {
			s.CloseCh()
			s.CloseSock()
			return nil
		}

		select {
		case <-exitWaitTimeout:
			s.CloseCh()
			s.CloseSock()
			return fmt.Errorf(s.GetBufferInfoMsg())
		case buf := <-s.in:
			s.sendFrame(s.soc, buf, true)
		case r := <-s.retryCh:
			s.retry(s.soc, r)
		}
	}
}

func (s *Sock) recovery(f func()) {
	if r := recover(); r != nil {
		if err, ok := r.(runtime.Error); ok {
			log.Error().Msgf("recover: %s", err.Error())
		} else {
			log.Error().Msgf("recover: %v", r)
		}

		if s.IsAutoRestart() {
			t := timer.AcquireTimer(time.Duration(s.ReconnectIvlMillSec) * time.Millisecond)

			if !t.Stop() {
				<-t.C
			}

			// recovery the f
			fPtr := reflect.ValueOf(f).Pointer()
			log.Warn().Msgf("recover: func %s", runtime.FuncForPC(fPtr).Name())
			Pool.CtxGo(s.ctx, func() {
				f()
			})

			timer.ReleaseTimer(t)

			// exit current thread
			return
		}

		// no restart will exit current thread directlly
		return
	}
}

// Publisher sends msg in 'in' channel. optional sock type: PUB/PUSH
func (s *Sock) Publisher() {
	switch s.Type {
	case goczmq.Push, goczmq.Pub:
	default:
		panic(fmt.Errorf("publisher only enables by 'type': Push/Pub"))
	}

	defer s.recovery(s.Publisher)

	_, err := s.Attach()
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-s.ctx.Done():
			err := s.Release()
			if err != nil {
				log.Error().Str("func", "Publisher").Msgf("Release err: %s", err.Error())
			}
			return
		case b := <-s.in:
			s.sendFrame(s.soc, b, true)
		case r := <-s.retryCh:
			s.retry(s.soc, r)
		}
	}
}

func (s *Sock) Consumer() {
	switch s.Type {
	case goczmq.Pull, goczmq.Sub:
	default:
		panic(fmt.Errorf("consumer only enables by 'type': Pull/Sub"))
	}

	defer s.recovery(s.Consumer)

	_, err := s.Attach()
	if err != nil {
		panic(err)
	}

	s.recvFrame(s.soc)

	<-s.ctx.Done()

	err = s.Release()
	if err != nil {
		log.Error().Str("func", "Consumer").Msgf("Release err: %s", err.Error())
	}
}

func (s *Sock) ServiceRegister() {
	if s.ServiceRegisterEndpoint == "" {
		log.Warn().Str("func", "ServiceRegister").
			Msgf("sock has empty 'ServiceRegisterEndpoint' which will skip sending zmq register msg")
		return
	}

	defer s.recovery(s.Publisher)

	soc, err := goczmq.NewReq(s.ServiceRegisterEndpoint)
	if err != nil {
		panic(err)
	}

	t := time.NewTicker(time.Duration(s.ReconnectIvlMillSec) * time.Millisecond)
	dt := &RegisterMsg{
		ID:   s.id,
		Host: iputil.GetIP(),
		TaskCount: &TaskCount{
			In:    int32(s.GetInCount()),
			Out:   int32(s.GetOutCount()),
			Retry: int32(s.GetRetryCount()),
		},
		SocketType: uint64(s.Type),
	}

	for {
		select {
		case <-s.ctx.Done():
			t.Stop()
			soc.Destroy()
			return
		case <-t.C:
			Pool.CtxGo(s.ctx, func() {
				dt.Timestamp = timestamppb.New(time.Now())
				if b, err := dt.Marshal(); err == nil {
					if err := soc.SendFrame(b, goczmq.FlagNone); err != nil {
						t.Stop()
						panic(err)
					}

					if _, _, err := soc.RecvFrame(); err != nil {
						t.Stop()
						panic(err)
					}
				}
			})
		}
	}
}
