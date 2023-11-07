package sock

import (
	"github.com/lafrinte/nops/iputil"
	A "github.com/stretchr/testify/assert"
	"github.com/zeromq/goczmq"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	assert := A.New(t)

	now := time.Now()
	msg := &RegisterMsg{
		ID:   "a",
		Host: iputil.GetIP(),
		TaskCount: &TaskCount{
			In:    int32(1),
			Out:   int32(0),
			Retry: int32(0),
		},
		SocketType: uint64(goczmq.Pull),
		Timestamp:  timestamppb.New(now),
	}

	assert.Equal(msg.GetID(), "a")
	assert.Equal(msg.GetHost(), iputil.GetIP())
	assert.Equal(msg.GetSocketType(), uint64(goczmq.Pull))
	assert.Equal(msg.GetTaskCount().GetIn(), int32(1))
	assert.Equal(msg.GetTaskCount().GetOut(), int32(0))
	assert.Equal(msg.GetTaskCount().GetRetry(), int32(0))
	assert.Equal(msg.GetTimestamp().GetSeconds(), now.Unix())
	assert.Equal(msg.GetTimestamp().GetNanos(), int32(now.Nanosecond()))

	buf, err := msg.Marshal()
	assert.Nil(err)

	newMsg := &RegisterMsg{}
	err = newMsg.Unmarshal(buf)
	assert.Nil(err)

	assert.Equal(msg.GetID(), newMsg.GetID())
	assert.Equal(msg.GetHost(), newMsg.GetHost())
	assert.Equal(msg.GetSocketType(), newMsg.GetSocketType())
	assert.Equal(msg.GetTaskCount().GetIn(), newMsg.GetTaskCount().GetIn())
	assert.Equal(msg.GetTaskCount().GetOut(), newMsg.GetTaskCount().GetOut())
	assert.Equal(msg.GetTaskCount().GetRetry(), newMsg.GetTaskCount().GetRetry())
}
