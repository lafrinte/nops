package sock

import (
	"context"
	"fmt"
	A "github.com/stretchr/testify/assert"
	"github.com/zeromq/goczmq"
	"strconv"
	"testing"
	"time"
)

func TestNewAndOption(t *testing.T) {
	assert := A.New(t)

	soc := New(
		WithCtx(context.Background()),
		WithSndhwm(1000),
		WithType("PUB"),
		WithEndpoint("inproc://xsub-test"),
		WithExitWaitTimeout(time.Second*10),
		WithRetryAttempts(3),
		WithRetryInterval(time.Second),
		WithMaxBufferSize(1000),
		EnableTcpKeepAlive(),
		WithTcpKeepAliveIdleSec(15),
		WithTcpKeepAliveCnt(3),
		WithHeartbeatIvlSec(5),
		WithHeartbeatTimoutSec(10),
		WithHeartbeatTTLSec(20),
		WithServiceRegisterEndpoint("inproc://register-test"),
	)

	assert.Equal(soc.Sndhwm, 1000)
	assert.Equal(soc.Type, goczmq.Pub)
	assert.Equal(soc.Endpoint, "inproc://xsub-test")
	assert.Equal(soc.ServiceRegisterEndpoint, "inproc://register-test")
	assert.Equal(soc.TcpKeepAliveIdleSec, int16(15))
	assert.Equal(soc.TcpKeepAliveCnt, int8(3))
	assert.Equal(soc.EnableTcpKeepAlive, true)
	assert.Equal(cap(soc.GetInChannel()), 1000)
	assert.Equal(cap(soc.GetOutChannel()), 0)
	assert.Equal(cap(soc.retryCh), 1000)
}

func TestSimplePub(t *testing.T) {
	assert := A.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	endpoint := "inproc://broadcast"
	soc := New(
		WithCtx(ctx),
		WithType("PUB"),
		WithEndpoint(endpoint),
		WithMaxBufferSize(4000),
	)

	go soc.Publisher()

	go func() {
		in := soc.GetInChannel()
		for i := 0; i < 4000; i++ {
			in <- []byte(strconv.Itoa(i))
		}
	}()

	<-ctx.Done()

	assert.Equal(soc.GetInCount(), 0)
}

func TestPubWithTypePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			expectErr := "publisher only enables by 'type': Push/Pub"

			if fmt.Sprintf("%s", r) != expectErr {
				t.Errorf("no panic raised when calling Publisher with type 'sub', got %v", r)
			}
		}
	}()
	endpoint := "inproc://broadcast"
	pub := New(
		WithType("SUB"),
		WithEndpoint(endpoint),
	)

	pub.Publisher()
}

func TestSimpleSub(t *testing.T) {
	assert := A.New(t)

	endpoint := "inproc://xsub-test"
	pub := goczmq.NewSock(goczmq.Pub)
	if _, err := pub.Bind(endpoint); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	soc := New(
		WithCtx(ctx),
		WithSndhwm(10000),
		WithType("Sub"),
		WithEndpoint("inproc://xsub-test"),
		WithAttach(),
		WithMaxBufferSize(4000),
	)

	go soc.Consumer()

	time.Sleep(time.Second)

	for i := 0; i < 4000; i++ {
		time.Sleep(time.Nanosecond) // mock real world situation
		if err := pub.SendFrame([]byte(strconv.Itoa(i)), goczmq.FlagNone); err != nil {
			panic(err)
		}
	}

	<-ctx.Done()

	assert.Equal(soc.GetOutCount(), 4000)
}

func TestSubWithTypePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			expectErr := "consumer only enables by 'type': Pull/Sub"

			if fmt.Sprintf("%s", r) != expectErr {
				t.Errorf("no panic raised when calling Publisher with type 'sub', got %v", r)
			}
		}
	}()
	endpoint := "inproc://broadcast"
	soc := New(
		WithType("PUB"),
		WithEndpoint(endpoint),
	)

	soc.Consumer()
}

func TestSimplePush(t *testing.T) {
	assert := A.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	endpoint := "inproc://push"
	soc := New(
		WithCtx(ctx),
		WithType("Push"),
		WithEndpoint(endpoint),
		WithMaxBufferSize(1000),
	)

	in := soc.GetInChannel()
	for i := 0; i < 1000; i++ {
		in <- []byte(strconv.Itoa(i))
	}

	go soc.Publisher()

	<-ctx.Done()

	assert.Equal(soc.GetOutCount(), 0)
}

func TestSimplePull(t *testing.T) {
	assert := A.New(t)

	endpoint := "inproc://xpull-test"
	push := goczmq.NewSock(goczmq.Push)
	if _, err := push.Bind(endpoint); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	soc := New(
		WithCtx(ctx),
		WithSndhwm(10000),
		WithType("Pull"),
		WithEndpoint(endpoint),
		WithAttach(),
		WithMaxBufferSize(1000),
	)

	go soc.Consumer()

	time.Sleep(time.Second)

	for i := 0; i < 1000; i++ {
		time.Sleep(time.Nanosecond) // mock real world situation
		if err := push.SendFrame([]byte(strconv.Itoa(i)), goczmq.FlagNone); err != nil {
			panic(err)
		}
	}

	<-ctx.Done()

	assert.Equal(soc.GetOutCount(), 1000)
}

func TestReconnect(t *testing.T) {

	endpoint := "tcp://0.0.0.0:5555"
	pub := goczmq.NewSock(goczmq.Pub)
	if _, err := pub.Bind(endpoint); err != nil {
		t.Error(err)
	}

	pub.SetReconnectIvl(100)

}
