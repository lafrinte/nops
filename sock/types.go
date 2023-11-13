package sock

import (
	"context"
	"fmt"
	"github.com/lafrinte/nops/atom"
	"github.com/lafrinte/nops/timer"
	"github.com/zeromq/goczmq"
	"reflect"
	"runtime"
	"time"
)

type RetryMsg struct {
	retry    uint8
	Msg      []byte
	MaxRetry uint8
}

// Retry get the retry state. if the retry time reach the max retry times, Retry will return false.
func (r *RetryMsg) Retry() bool {
	return r.MaxRetry > r.retry
}

// IterRetryTimes used to iter the retry times
func (r *RetryMsg) IterRetryTimes() {
	r.retry++
}

// GetRetryTimes gets msg current retry times
func (r *RetryMsg) GetRetryTimes() uint8 {
	return r.retry
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

	//
	sendMsgCount uint64
	dropMsgCount uint64
	recvMsgCount uint64

	// safe destroy args
	ExitWaitTimeout time.Duration

	SendTimeoutSec uint16
	RecvTimeoutSec uint16
}

// bind binds socket on endpoint
func (s *Sock) bind() (*goczmq.Sock, error) {
	soc := s.setOptions()
	if _, err := soc.Bind(s.Endpoint); err != nil {
		return nil, err
	}

	return soc, nil
}

// connect makes connection to endpoint
func (s *Sock) connect() (*goczmq.Sock, error) {
	soc := s.setOptions()
	if err := soc.Connect(s.Endpoint); err != nil {
		return nil, err
	}

	return soc, nil
}

// setOptions sets socket tcp and heartbeat option
func (s *Sock) setOptions() *goczmq.Sock {
	soc := goczmq.NewSock(s.Type)

	soc.SetSndhwm(s.Sndhwm)

	if soc.GetType() == goczmq.Sub {
		log.Debug().Msg("sub mode will set default subscribe to ''")
		soc.SetSubscribe("")
	}

	if s.SendTimeoutSec > 0 {
		soc.SetSndtimeo(int(s.SendTimeoutSec * 1000))
	}

	if s.RecvTimeoutSec > 0 {
		soc.SetRcvtimeo(int(s.SendTimeoutSec * 1000))
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

func (s *Sock) release() {

	close(s.in)
	close(s.out)
	close(s.retryCh)

	if s.soc != nil {
		s.soc.Destroy()
		s.soc = nil
	}
}

// sendFrame tries sending msg and puts msg into retry channel when any error occurred
func (s *Sock) sendFrame(sock *goczmq.Sock, msg []byte, retry bool) error {
	/* goczmq.Sock will panic in C while receiving or sending frame on a destroyed socket,
	   wrapped in thread will transfer thread-panic to main thread and be caught
	*/
	if sock == nil {
		log.Error().Err(fmt.Errorf("sock pointer is nil")).Msg("sock may closed, exit now")
		return fmt.Errorf("sock pointer is nil")
	}

	err := sock.SendFrame(msg, goczmq.FlagNone)
	if err != nil {
		log.Error().Err(err).Bytes("data", msg).Msg("failed to send")

		if retry {
			log.Info().Bytes("data", msg).Msg("retry send")
			s.retryCh <- NewRetryMsg(msg, s.RetryAttempts)

			// not return error while msg can retry
			return nil
		}

		s.dropMsgCount++

		return fmt.Errorf("failed to send")
	}

	s.sendMsgCount++

	return nil
}

// recvFrame tries receiving msg and puts into out channel
func (s *Sock) recvFrame(sock *goczmq.Sock) error {
	/* goczmq.Sock will panic in C while receiving or sending frame on a destroyed socket,
	   wrapped in thread will transfer thread-panic to main thread and be caught
	*/
	for {
		if sock == nil {
			log.Error().Err(fmt.Errorf("sock pointer is nil")).Msg("sock may closed, exit now")
			return fmt.Errorf("sock pointer is nil")
		}

		buf, _, err := sock.RecvFrame()
		if err != nil {
			if err == goczmq.ErrRecvFrameAfterDestroy {
				log.Error().Err(err).Msg("call RecvFrame after sock been destroyed")
				panic(err)
			}

			log.Error().Err(err).Msg("RecvFrame failed")

			continue
		}

		s.out <- buf
		s.recvMsgCount++
	}
}

// retry tries to publish msg in retry channel and re-puts into retry channel when any error occured
func (s *Sock) retry(sock *goczmq.Sock, msg *RetryMsg) {
	if msg.Retry() {
		if err := sock.SendFrame(msg.Msg, goczmq.FlagNone); err != nil {
			msg.IterRetryTimes()

			log.Error().Err(err).Bytes("data", msg.Msg).Msgf("retry SendFrame failed the %d time", msg.GetRetryTimes())
			t := timer.AcquireTimer(s.RetryInterval)
			if !t.Stop() {
				<-t.C
			}

			timer.ReleaseTimer(t)

			s.retryCh <- msg
			return
		}
	}
}

// Release tries release socket after all buffers be triggered
func (s *Sock) Release() error {
	s.StopAutoRestart()

	if s.EmptyBuffer() {
		s.release()
		return nil
	}

	exitWaitTimeout := time.After(s.ExitWaitTimeout)
	for {
		if s.EmptyBuffer() {
			s.release()
			return nil
		}

		select {
		case <-exitWaitTimeout:
			s.release()
			msg := fmt.Sprintf("msg lost: in [%d] out [%d] retry [%d]", s.GetInCount(), s.GetOutCount(), s.GetRetryCount())
			return fmt.Errorf(msg)
		case buf := <-s.in:
			_ = s.sendFrame(s.soc, buf, true)
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
		log.Panic().Err(err).Msg("panic on socket Attach")
		panic(err)
	}

	for {
		select {
		case <-s.ctx.Done():
			if err := s.Release(); err != nil {
				log.Error().Str("func", "Publisher").Msgf("Release err: %s", err.Error())
			}
			return
		case b := <-s.in:
			_ = s.sendFrame(s.soc, b, true)
		case r := <-s.retryCh:
			s.retry(s.soc, r)
		}
	}
}

// Consumer receive msg from sock and charge into 'out' channel. optional socket type: SUB/PULL
func (s *Sock) Consumer() {
	switch s.Type {
	case goczmq.Pull, goczmq.Sub:
	default:
		panic(fmt.Errorf("consumer only enables by 'type': Pull/Sub"))
	}

	defer s.recovery(s.Consumer)

	_, err := s.Attach()
	if err != nil {
		log.Panic().Err(err).Msg("panic on socket Attach")
		panic(err)
	}

	_ = s.recvFrame(s.soc)

	<-s.ctx.Done()

	err = s.Release()
	if err != nil {
		log.Error().Str("func", "Consumer").Msgf("Release err: %s", err.Error())
	}
}

// Requester send request msg in 'in' channel and save reply msg in 'out' channel
func (s *Sock) Requester() {
	switch s.Type {
	case goczmq.Req:
	default:
		panic(fmt.Errorf("requester only enables by 'type': Req"))
	}

	defer s.recovery(s.Requester)

	_, err := s.Attach()
	if err != nil {
		log.Panic().Err(err).Msg("panic on socket Attach")
		panic(err)
	}

	for {
		select {
		case <-s.ctx.Done():
			if err := s.Release(); err != nil {
				log.Error().Str("func", "Requester").Msgf("Release err: %s", err.Error())
			}
			return
		case b := <-s.in:
			/* Sets the timeout for send operation on the socket.
			   If the value is 0, zmq_send(3) will return immediately, with a EAGAIN error if the message cannot be sent.
			   If the value is -1, it will block until the message is sent.
			   For all other values, it will try to send the message for that amount of time before returning with an EAGAIN error
			*/
			if err := s.sendFrame(s.soc, b, false); err == nil {
				/* Sets the timeout for receive operation on the socket. If the value is 0, zmq_recv(3)
				   will return immediately, with a EAGAIN error if there is no message to receive. If the value is -1,
				   it will block until a message is available. For all other values,
				   it will wait for a message for that amount of time before returning with an EAGAIN error.
				*/
				reply, _, err := s.soc.RecvFrame()
				if err != nil {
					log.Error().Err(err).Msgf("requester get no replay at %s", time.Now())

					// charge empty when an error occurred in RecvFrame
					s.out <- []byte("")
					return
				}

				// block when msg in 'out' has not been consumed.
				s.out <- reply
				s.recvMsgCount++
			}

			// charge empty when an error occurred in sendFrame
			s.out <- []byte("")
			s.sendMsgCount++
		}
	}
}

// Responser recharge request msg into 'out' channel and get its response msg from 'in' channel
func (s *Sock) Responser() {
	switch s.Type {
	case goczmq.Rep:
	default:
		panic(fmt.Errorf("requester only enables by 'type': Rep"))
	}

	defer s.recovery(s.Responser)

	_, err := s.Attach()
	if err != nil {
		log.Panic().Err(err).Msg("panic on socket Attach")
		panic(err)
	}

	for {
		select {
		case <-s.ctx.Done():
			if err := s.Release(); err != nil {
				log.Error().Str("func", "Responser").Msgf("Release err: %s", err.Error())
			}
			return
		default:
			/* Sets the timeout for receive operation on the socket. If the value is 0, zmq_recv(3)
			   will return immediately, with a EAGAIN error if there is no message to receive. If the value is -1,
			   it will block until a message is available. For all other values,
			   it will wait for a message for that amount of time before returning with an EAGAIN error.
			*/
			if request, _, err := s.soc.RecvFrame(); err == nil {
				/* Sets the timeout for send operation on the socket.
				   If the value is 0, zmq_send(3) will return immediately, with a EAGAIN error if the message cannot be sent.
				   If the value is -1, it will block until the message is sent.
				   For all other values, it will try to send the message for that amount of time before returning with an EAGAIN error
				*/
				s.recvMsgCount++

				s.out <- request
				response := <-s.in
				if err := s.sendFrame(s.soc, response, false); err == nil {
					s.sendMsgCount++
				}
			}
		}
	}
}
