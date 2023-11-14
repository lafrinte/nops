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
		WithEndpoint("inproc://test"),
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
		WithSendTimeoutSec(3),
		WithRecvTimeoutSec(3),
		WithInChannel(make(chan []byte, 1000)),
		WithOutChannel(make(chan []byte, 100)),
	)

	assert.Equal(soc.Sndhwm, 1000)
	assert.Equal(soc.Type, goczmq.Pub)
	assert.Equal(soc.Endpoint, "inproc://sub")
	assert.Equal(soc.TcpKeepAliveIdleSec, int16(15))
	assert.Equal(soc.TcpKeepAliveCnt, int8(3))
	assert.Equal(soc.EnableTcpKeepAlive, true)
	assert.Equal(cap(soc.GetInChannel()), 1000)
	assert.Equal(cap(soc.GetOutChannel()), 0)
	assert.Equal(cap(soc.retryCh), 1000)

	_ = New(
		WithType("SUB"),
		WithEndpoint("inproc://test"),
	)

	_ = New(
		WithType("Req"),
		WithEndpoint("inproc://test"),
	)
}

func TestBroadcastBindOnPub(t *testing.T) {
	assert := A.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	endpoint := "inproc://pub"

	pub := New(
		WithCtx(ctx),
		WithType("Pub"),
		WithEndpoint(endpoint),
		WithMaxBufferSize(1000),
	)

	in := pub.GetInChannel()
	go pub.Publisher()

	sub1 := New(
		WithCtx(ctx),
		WithType("Sub"),
		WithEndpoint(endpoint),
		WithAttach(),
		WithMaxBufferSize(1000),
	)

	out1 := sub1.GetOutChannel()
	go sub1.Consumer()

	sub2 := New(
		WithCtx(ctx),
		WithType("Sub"),
		WithEndpoint(endpoint),
		WithAttach(),
		WithMaxBufferSize(1000),
	)

	out2 := sub2.GetOutChannel()
	go sub2.Consumer()

	time.Sleep(time.Millisecond * 200) // wait sub1 and sub2 connection
	for i := 0; i < 1000; i++ {
		in <- []byte(strconv.Itoa(i))
	}

	<-ctx.Done()

	assert.Equal(len(in), 0)
	assert.Equal(len(out1), 1000)
	assert.Equal(len(out2), 1000)
}

func TestBroadcastBindOnSub(t *testing.T) {
	assert := A.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	endpoint := "inproc://sub"

	sub := New(
		WithCtx(ctx),
		WithType("Sub"),
		WithEndpoint(endpoint),
		WithMaxBufferSize(2000),
	)

	out := sub.GetOutChannel()
	go sub.Consumer()

	pub1 := New(
		WithCtx(ctx),
		WithType("Pub"),
		WithEndpoint(endpoint),
		WithAttach(),
		WithMaxBufferSize(1000),
	)

	in1 := pub1.GetInChannel()
	go pub1.Publisher()

	pub2 := New(
		WithCtx(ctx),
		WithType("Pub"),
		WithEndpoint(endpoint),
		WithAttach(),
		WithMaxBufferSize(1000),
	)

	in2 := pub2.GetInChannel()
	go pub2.Publisher()

	time.Sleep(time.Millisecond * 200) // wait sub1 and sub2 connection // wait sub1 and sub2 connection
	for i := 0; i < 1000; i++ {
		in1 <- []byte(strconv.Itoa(i))
		in2 <- []byte(strconv.Itoa(i))
	}

	<-ctx.Done()

	assert.Equal(len(in1), 0)
	assert.Equal(len(in2), 0)
	assert.Equal(len(out), 2000)
}

func TestQueueBindOnPush(t *testing.T) {
	assert := A.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	endpoint := "inproc://sub"

	push := New(
		WithCtx(ctx),
		WithType("Push"),
		WithEndpoint(endpoint),
		WithMaxBufferSize(1000),
	)

	in := push.GetInChannel()
	go push.Publisher()

	pull1 := New(
		WithCtx(ctx),
		WithType("Pull"),
		WithEndpoint(endpoint),
		WithAttach(),
		WithMaxBufferSize(1000),
	)

	out1 := pull1.GetOutChannel()
	go pull1.Consumer()

	pull2 := New(
		WithCtx(ctx),
		WithType("Pull"),
		WithEndpoint(endpoint),
		WithAttach(),
		WithMaxBufferSize(1000),
	)

	out2 := pull2.GetOutChannel()
	go pull2.Consumer()

	time.Sleep(time.Millisecond * 200) // wait sub1 and sub2 connection
	for i := 0; i < 1000; i++ {
		in <- []byte(strconv.Itoa(i))
	}

	<-ctx.Done()

	assert.Equal(len(in), 0)
	assert.Equal(len(out1)+len(out2), 1000)
}

func TestQueueBindOnPull(t *testing.T) {
	assert := A.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	endpoint := "inproc://sub"

	pull := New(
		WithCtx(ctx),
		WithType("Pull"),
		WithEndpoint(endpoint),
		WithMaxBufferSize(2000),
	)

	out := pull.GetOutChannel()
	go pull.Consumer()

	push1 := New(
		WithCtx(ctx),
		WithType("Push"),
		WithEndpoint(endpoint),
		WithAttach(),
		WithMaxBufferSize(1000),
	)

	in1 := push1.GetInChannel()
	go push1.Publisher()

	push2 := New(
		WithCtx(ctx),
		WithType("Push"),
		WithEndpoint(endpoint),
		WithAttach(),
		WithMaxBufferSize(1000),
	)

	in2 := push2.GetInChannel()
	go push2.Publisher()

	time.Sleep(time.Millisecond * 200) // wait sub1 and sub2 connection
	for i := 0; i < 1000; i++ {
		in1 <- []byte(strconv.Itoa(i))
		in2 <- []byte(strconv.Itoa(i))
	}

	<-ctx.Done()

	assert.Equal(len(in1), 0)
	assert.Equal(len(in2), 0)
	assert.Equal(len(out), 2000)
}

func TestReqRep(t *testing.T) {
	assert := A.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2500)
	defer cancel()

	endpoint := "inproc://rep"

	rep := New(
		WithCtx(ctx),
		WithType("Rep"),
		WithEndpoint(endpoint),
	)

	go rep.Responser()

	req := New(
		WithCtx(ctx),
		WithType("Req"),
		WithAttach(),
		WithEndpoint(endpoint),
	)

	go req.Requester()

	qIn := req.GetInChannel()
	qOut := req.GetOutChannel()
	pIn := rep.GetInChannel()
	pOut := rep.GetOutChannel()

	requestMsg := []byte("hello")
	qIn <- requestMsg

	msg := <-pOut
	assert.Equal(msg, requestMsg)

	responseMsg := []byte("world")
	pIn <- responseMsg

	msg = <-qOut
	assert.Equal(msg, responseMsg)

	<-ctx.Done()
}

func TestPublisherWithTypePanic(t *testing.T) {
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

func TestConsumerWithTypePanic(t *testing.T) {
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

func TestRequesterWithTypePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			expectErr := "requester only enables by 'type': Req"

			if fmt.Sprintf("%s", r) != expectErr {
				t.Errorf("no panic raised when calling Requester with type 'pub', got %v", r)
			}
		}
	}()
	endpoint := "inproc://broadcast"
	soc := New(
		WithType("PUB"),
		WithEndpoint(endpoint),
	)

	soc.Requester()
}

func TestResponserWithTypePanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			expectErr := "responser only enables by 'type': Rep"

			if fmt.Sprintf("%s", r) != expectErr {
				t.Errorf("no panic raised when calling Responser with type 'pub', got %v", r)
			}
		}
	}()
	endpoint := "inproc://broadcast"
	soc := New(
		WithType("PUB"),
		WithEndpoint(endpoint),
	)

	soc.Responser()
}
