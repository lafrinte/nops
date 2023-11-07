package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/lafrinte/nops/fs"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strings"
)

type Option func(c *Config)

func WithAllowEmptyEnv() Option {
	return func(c *Config) {
		c.viper.AllowEmptyEnv(true)
	}
}

func WithEnvKeyReplacer(r *strings.Replacer) Option {
	return func(c *Config) {
		c.viper.SetEnvKeyReplacer(r)
	}
}

func WithAutomaticEnv() Option {
	return func(c *Config) {
		c.viper.AutomaticEnv()
	}
}

func WithBindEnv(envs map[string]string) Option {
	return func(c *Config) {
		for n, v := range envs {
			_ = c.viper.BindEnv(n, v)
		}
	}
}

func WithOptionConfigPath(path []string) Option {
	return func(c *Config) {
		for _, p := range path {
			c.viper.AddConfigPath(p)
		}
	}
}

func WithConfigName(in string) Option {
	return func(c *Config) {
		c.viper.SetConfigName(in)
	}
}

func WithConfigType(in string) Option {
	return func(c *Config) {
		c.viper.SetConfigType(in)
	}
}

func WithFsNotify() Option {
	return func(c *Config) {
		c.viper.WatchConfig()
		c.viper.OnConfigChange(func(event fsnotify.Event) {
			log.Info().Str("action", "conf").Msg(fmt.Sprintf("conf changed: %s", event))
		})
	}
}

func WithTemplate(template string) Option {
	return func(c *Config) {
		c.Template = strings.Replace(template, "\t", "    ", -1)
	}
}

func WriteTemplateFromFile(path string) Option {
	return func(c *Config) {
		if s, err := fs.ReadFile(path); err == nil {
			c.Template = strings.Replace(s, "\t", "    ", -1)
		}
	}
}

func WithWriteTo(path string) Option {
	return func(c *Config) {
		c.WriteTo = path
		dir, file := filepath.Split(path)
		c.viper.AddConfigPath(dir)
		c.viper.SetConfigName(file)
	}
}

// WithDefaultVal is used to set the conf struct unmarshalling to
func WithDefaultVal(i interface{}) Option {
	return func(c *Config) {
		c.DefaultVal = i
	}
}
