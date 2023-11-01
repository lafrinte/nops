package conf

import (
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/spf13/viper"
	"nops/fs"
	"nops/str"
	"regexp"
	"time"
)

type Config struct {
	viper      *viper.Viper
	Template   string
	WriteTo    string
	DefaultVal interface{}
}

/*
Write is used in cli for generate conf file with conf template which may contains description message for
each section and option.
*/
func (c *Config) Write() error {
	tpl, err := pongo2.FromString(c.Template)
	if err != nil {
		return fmt.Errorf("read template failed: err -> %s", err)
	}

	val := pongo2.Context{}
	if c.DefaultVal != nil {
		for k, v := range c.DefaultVal.(map[string]interface{}) {
			val[k] = v
		}
	}

	s, err := tpl.Execute(val)
	if err != nil {
		return fmt.Errorf("parsing conf failed: err -> %s", err)
	}

	// replace string
	// 1. prevent indent error: turn tab to 4 blank space
	// use lower case for bool type object parsing by pongo2. default True/False -> true/false
	s = str.Replaces(
		s,
		str.ReplacePoint{Old: "\t", New: "    ", N: -1},
		str.ReplacePoint{Old: ": True", New: ": true", N: -1},
		str.ReplacePoint{Old: ": False", New: ": false", N: -1},
	)

	// '\n    \n     \n' -> \n
	s = regexp.MustCompile("(?:\\n\\s*){2,}\\n").ReplaceAllString(s, "\n")

	// '\n     \n' -> '\n\n'
	s = regexp.MustCompile("\\n\\s+\\n").ReplaceAllString(s, "\n\n")

	// '\n\n    -' -> '\n    -'
	s = regexp.MustCompile("\\n\\n(\\s+-)").ReplaceAllString(s, "\n${1}")

	if err := fs.WriteFile(c.WriteTo, s, 0666); err != nil {
		return fmt.Errorf("write into conf failed: err -> %s", err)
	}

	return nil
}

/*
Read is used to parsing conf. when conf does not exist, func will try to write into conf file
from template parsing with default value.
*/
func (c *Config) Read() error {
	if err := c.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("no conf file, and write default template conf failed: err -> %s", err)
		}

		// parsing error will return directly
		return fmt.Errorf("reading failed: err -> %s", err)
	}

	return nil
}

func (c *Config) Get(key string) interface{} {
	return c.viper.Get(key)
}

func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *Config) GetInt32(key string) int32 {
	return c.viper.GetInt32(key)
}

func (c *Config) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

func (c *Config) GetUint(key string) uint {
	return c.viper.GetUint(key)
}

func (c *Config) GetUint16(key string) uint16 {
	return c.viper.GetUint16(key)
}

func (c *Config) GetUint32(key string) uint32 {
	return c.viper.GetUint32(key)
}

func (c *Config) GetUint64(key string) uint64 {
	return c.viper.GetUint64(key)
}

func (c *Config) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

func (c *Config) GetTime(key string) time.Time {
	return c.viper.GetTime(key)
}

func (c *Config) GetDuration(key string) time.Duration {
	return c.viper.GetDuration(key)
}

func (c *Config) GetIntSlice(key string) []int {
	return c.viper.GetIntSlice(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

func (c *Config) GetStringMap(key string) map[string]interface{} {
	return c.viper.GetStringMap(key)
}

func (c *Config) GetStringMapString(key string) map[string]string {
	return c.viper.GetStringMapString(key)
}

func (c *Config) GetStringMapStringSlice(key string) map[string][]string {
	return c.viper.GetStringMapStringSlice(key)
}

func (c *Config) GetMapSlice(key string) []map[string]interface{} {
	var out []map[string]interface{}
	dt := c.viper.Get(key).([]interface{})
	for _, v := range dt {
		out = append(out, v.(map[string]interface{}))
	}

	return out
}
