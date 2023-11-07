package pool

import (
	"github.com/bytedance/gopkg/util/gopool"
	"math"
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

func New(name string) gopool.Pool {
	if name == "" {
		name = GlobalPoolName
	}

	pool := gopool.GetPool(name)
	if pool != nil {
		return pool
	}

	pool = gopool.NewPool(
		GlobalPoolName,
		math.MaxInt32,
		gopool.NewConfig(),
	)

	_ = gopool.RegisterPool(pool)

	return pool
}

func DefaultPool() gopool.Pool {
	return gopool.GetPool(GlobalPoolName)
}
