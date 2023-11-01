package str

import (
	"bytes"
	A "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReplaces(t *testing.T) {
	assert := A.New(t)

	baseString := "hello world!"
	expectString := "e#lo world!"
	rs := []ReplacePoint{
		{
			Old: "h",
			New: "",
		},
		{
			Old: "l",
			New: "#",
			N:   1,
		},
	}

	assert.Equal(expectString, Replaces(baseString, rs...))
}

func TestSplit(t *testing.T) {
	assert := A.New(t)

	baseString := "a  b c d e"
	expect := []string{"a", "b", "c", "d", "e"}
	assert.Equal(expect, Split(baseString, " "))
}

func TestIsAlpha(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Text   string
		Expect bool
	}{
		{
			Name:   "only alpha",
			Text:   "djlfasdhalk",
			Expect: true,
		},
		{
			Name:   "alpha and numeric",
			Text:   "asdfjlk1iu238197",
			Expect: false,
		},
		{
			Name:   "blankspace",
			Text:   " ",
			Expect: false,
		},
		{
			Name:   "special chars",
			Text:   "*&^%$#",
			Expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert.Equal(IsAlpha(c.Text), c.Expect)
		})
	}
}

func TestIsAlphanumeric(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Text   string
		Expect bool
	}{
		{
			Name:   "only alpha",
			Text:   "hello",
			Expect: true,
		},
		{
			Name:   "alpha and numeric",
			Text:   "hello123",
			Expect: true,
		},
		{
			Name:   "backspace",
			Text:   " ",
			Expect: false,
		},
		{
			Name:   "special chars",
			Text:   "*&^%$#",
			Expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert.Equal(IsAlphanumeric(c.Text), c.Expect)
		})
	}
}

func TestIsNumeric(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Text   string
		Expect bool
	}{
		{
			Name:   "only numeric",
			Text:   "123123",
			Expect: true,
		},
		{
			Name:   "only alpha",
			Text:   "djlfasdhalk",
			Expect: false,
		},
		{
			Name:   "alpha and numeric",
			Text:   "asdfjlk1iu238197",
			Expect: false,
		},
		{
			Name:   "blankspace",
			Text:   " ",
			Expect: false,
		},
		{
			Name:   "special chars",
			Text:   "*&^%$#",
			Expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert.Equal(IsNumeric(c.Text), c.Expect)
		})
	}
}

func TestIsUpper(t *testing.T) {
	assert := A.New(t)

	assert.True(IsUpper("ABC"))
	assert.False(IsUpper("ABc"))
}

func TestIsLower(t *testing.T) {
	assert := A.New(t)

	assert.False(IsLower("ABc"))
	assert.True(IsLower("abc"))
}

func TestHasLower(t *testing.T) {
	assert := A.New(t)

	assert.True(HasLower("ABc"))
	assert.False(HasLower("ABC"))
}

func TestHasUpper(t *testing.T) {
	assert := A.New(t)

	assert.True(HasUpper("ABc"))
	assert.False(HasUpper("abc"))
}

func TestHasEmpty(t *testing.T) {
	assert := A.New(t)

	assert.True(HasEmpty("a b"))
	assert.False(HasEmpty("ab"))
}

func TestIsMatch(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name     string
		Text     string
		Funcs    []func(raw string, pattern string) bool
		Patterns []string
		And      bool
		Expect   bool
	}{
		{
			Name: "only contains",
			Text: "hello world!",
			Funcs: []func(raw string, pattern string) bool{
				func(raw string, pattern string) bool {
					return strings.Contains(raw, pattern)
				},
			},
			Patterns: []string{"world"},
			And:      false,
			Expect:   true,
		},
		{
			Name: "only prefix",
			Text: "hello world!",
			Funcs: []func(raw string, pattern string) bool{
				func(raw string, pattern string) bool {
					return strings.HasPrefix(raw, pattern)
				},
			},
			Patterns: []string{"hello"},
			And:      false,
			Expect:   true,
		},
		{
			Name: "only suffix",
			Text: "hello world!",
			Funcs: []func(raw string, pattern string) bool{
				func(raw string, pattern string) bool {
					return strings.HasSuffix(raw, pattern)
				},
			},
			Patterns: []string{"world!"},
			And:      false,
			Expect:   true,
		},
		{
			Name: "both prefix and suffix",
			Text: "hello world!",
			Funcs: []func(raw string, pattern string) bool{
				func(raw string, pattern string) bool {
					return strings.HasPrefix(raw, pattern)
				},
				func(raw string, pattern string) bool {
					return strings.HasSuffix(raw, pattern)
				},
			},
			Patterns: []string{"hello", "world!"},
			And:      true,
			Expect:   true,
		},
		{
			Name: "both prefix and contains",
			Text: "hello world!",
			Funcs: []func(raw string, pattern string) bool{
				func(raw string, pattern string) bool {
					return strings.HasPrefix(raw, pattern)
				},
				func(raw string, pattern string) bool {
					return strings.Contains(raw, pattern)
				},
			},
			Patterns: []string{"hello", "world"},
			And:      true,
			Expect:   true,
		},
		{
			Name: "prefix or contains",
			Text: "hello world!",
			Funcs: []func(raw string, pattern string) bool{
				func(raw string, pattern string) bool {
					return strings.HasPrefix(raw, pattern)
				},
				func(raw string, pattern string) bool {
					return strings.Contains(raw, pattern)
				},
			},
			Patterns: []string{"hello", "ppt"},
			And:      false,
			Expect:   true,
		},
		{
			Name: "both prefix and contains failed",
			Text: "hello world!",
			Funcs: []func(raw string, pattern string) bool{
				func(raw string, pattern string) bool {
					return strings.HasPrefix(raw, pattern)
				},
				func(raw string, pattern string) bool {
					return strings.Contains(raw, pattern)
				},
			},
			Patterns: []string{"hello", "ppt"},
			And:      true,
			Expect:   false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			opts := []MatchOption{
				WithRawString(c.Text),
			}

			assert.Equal(len(c.Funcs), len(c.Patterns))

			for i := 0; i < len(c.Funcs); i++ {
				opts = append(opts, WithFunc(c.Funcs[i], c.Patterns[i]))
			}

			assert.Equal(IsMatch(c.And, opts...), c.Expect)
		})
	}
}

func TestByteAppend(t *testing.T) {
	assert := A.New(t)

	data := [][]byte{
		[]byte("0"),
		[]byte("1"),
		[]byte("2"),
		[]byte("3"),
		[]byte("4"),
		[]byte("5"),
	}

	assert.Equal(bytes.Join(data, []byte("")), ByteAppend(data...))
}

func TestStringAppend(t *testing.T) {
	assert := A.New(t)

	data := []string{
		"1", "2", "3", "4", "5",
	}

	assert.Equal(strings.Join(data, ""), StringAppend(data...))
}

func TestLastLine(t *testing.T) {
	assert := A.New(t)

	data := `a
b
c
d
e
f`
	assert.Equal("f", LastLine(data))
}
