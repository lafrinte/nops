package pool

import (
	"github.com/bytedance/gopkg/util/gopool"
)

var (
	GlobalPoolName = "nops"
)

func init() {
	if err := gopool.RegisterPool(
		gopool.NewPool(
			GlobalPoolName,
			int32((1<<14)-1),
			gopool.NewConfig(),
		),
	); err != nil {
		panic(err)
	}
}

func New(name string, capacity int32) gopool.Pool {
	if name == "" {
		name = GlobalPoolName
	}

	if capacity == 0 {
		capacity = int32((1 << 14) - 1)
	}

	pool := gopool.GetPool(name)
	if pool != nil {
		return pool
	}

	pool = gopool.NewPool(
		GlobalPoolName,
		capacity,
		gopool.NewConfig(),
	)

	_ = gopool.RegisterPool(pool)

	return pool
}

func DefaultPool() gopool.Pool {
	return gopool.GetPool(GlobalPoolName)
}
