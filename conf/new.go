package conf

import (
	"github.com/spf13/viper"
	"sync"
)

var (
	GlobalViper *Config
	once        sync.Once
)

func New(opts ...Option) *Config {
	once.Do(func() {
		GlobalViper = new(Config)
		GlobalViper.viper = viper.New()

		for _, opt := range opts {
			opt(GlobalViper)
		}
	})

	return GlobalViper
}
