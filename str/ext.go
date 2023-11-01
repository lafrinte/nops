package str

import (
	"bytes"
	"regexp"
	"runtime"
	"strings"
	"unicode"
)

var Enter string

func init() {
	system := runtime.GOOS

	switch system {
	case "windows":
		Enter = "\r\n"
	default:
		Enter = "\n"
	}
}

type ReplacePoint struct {
	Old string
	New string
	N   int
}

func Replaces(s string, rs ...ReplacePoint) string {
	if len(rs) > 0 {
		for _, r := range rs {
			if r.N == 0 {
				r.N = -1
			}

			s = strings.Replace(s, r.Old, r.New, r.N)
		}
	}

	return s
}

func Split(s string, sep string) []string {
	var seg []string
	for _, sb := range strings.Split(s, sep) {
		if sb != "" {
			seg = append(seg, strings.TrimSpace(sb))
		}
	}

	return seg
}

func UpperCaseToUnderScore(s string) string {
	re := regexp.MustCompile(`(?U)([A-Z][a-z])`)
	snake := re.ReplaceAllString(s, "${1}_")

	snake = strings.ToLower(snake)
	snake = strings.TrimPrefix(snake, "_")
	snake = strings.TrimSuffix(snake, "_")

	return snake
}

func IsAlpha(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if !unicode.IsLetter(v) {
			return false
		}
	}
	return true
}

// IsAlphanumeric checks if the string contains only Unicode letters or digits.
func IsAlphanumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if !isAlphanumeric(v) {
			return false
		}
	}
	return true
}

// IsNumeric Checks if the string contains only digits. A decimal point is not a digit and returns false.
func IsNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}

func isAlphanumeric(v rune) bool {
	return unicode.IsDigit(v) || unicode.IsLetter(v)
}

func IsLower(s string) bool {
	if s == "" {
		return true
	}

	for _, v := range s {
		if !unicode.IsLower(v) {
			return false
		}
	}

	return true
}

func IsUpper(s string) bool {
	if s == "" {
		return true
	}

	for _, v := range s {
		if !unicode.IsUpper(v) {
			return false
		}
	}

	return true
}

func HasLower(s string) bool {
	if s == "" {
		return false
	}

	for _, v := range s {
		if unicode.IsLower(v) {
			return true
		}
	}

	return false
}

func HasUpper(s string) bool {
	if s == "" {
		return false
	}

	for _, v := range s {
		if unicode.IsUpper(v) {
			return true
		}
	}

	return false
}

func HasEmpty(args ...string) bool {
	for _, arg := range args {
		if len(arg) == 0 {
			return true
		}
	}

	return false
}

type MatchProxy struct {
	RawString string
	Func      []func(raw string, pattern string) bool
	Pattern   []string
}

type MatchOption func(m *MatchProxy)

func WithRawString(s string) MatchOption {
	return func(m *MatchProxy) {
		m.RawString = s
	}
}

func WithFunc(f func(raw string, pattern string) bool, pattern string) MatchOption {
	return func(m *MatchProxy) {
		m.Func = append(m.Func, f)
		m.Pattern = append(m.Pattern, pattern)
	}
}

func IsMatch(and bool, opts ...MatchOption) bool {
	m := &MatchProxy{}

	for _, opt := range opts {
		opt(m)
	}

	switch and {
	case true:
		for i, f := range m.Func {
			if !f(m.RawString, m.Pattern[i]) {
				return false
			}
		}
		return true
	case false:
		for i, f := range m.Func {
			if f(m.RawString, m.Pattern[i]) {
				return true
			}
		}
		return false
	}

	return false
}

func ByteAppend(b ...[]byte) []byte {
	buffer := bytes.Buffer{}
	for _, buf := range b {
		buffer.Write(buf)
	}

	return buffer.Bytes()
}

func StringAppend(ss ...string) string {
	builder := strings.Builder{}

	for _, s := range ss {
		builder.WriteString(s)
	}

	return builder.String()
}

func LastLineByte(buf []byte) []byte {
	index := bytes.LastIndex(buf, []byte(Enter))
	if index == -1 {
		return buf
	} else {
		return buf[index+1:]
	}
}

func LastLine(s string) string {
	return string(LastLineByte([]byte(s)))
}
