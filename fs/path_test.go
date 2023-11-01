package fs

import (
	"errors"
	A "github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestDirName(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Path   string
		Expect string
	}{
		{
			Name:   "simple file",
			Path:   "/Users/apple/xcode",
			Expect: "/Users/apple",
		},
		{
			Name:   "simple directory",
			Path:   "/Users/apple/",
			Expect: "/Users/apple",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert.Equal(c.Expect, GetDirName(c.Path))
		})
	}
}

func TestGetBaseName(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Path   string
		Expect string
	}{
		{
			Name:   "simple file",
			Path:   "/Users/apple/xcode",
			Expect: "xcode",
		},
		{
			Name:   "simple directory",
			Path:   "/Users/apple/",
			Expect: "apple",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert.Equal(c.Expect, GetBaseName(c.Path))
		})
	}
}

func TestGetExecDir(t *testing.T) {
	GetExecDir()
}

func TestIsExist(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Path   string
		Expect bool
		Dir    string
		Temp   bool
	}{
		{
			Name:   "exist dir path",
			Path:   "/tmp",
			Expect: true,
			Temp:   false,
		},
		{
			Name:   "exist file path",
			Dir:    "/tmp",
			Expect: true,
			Temp:   true,
		},
		{
			Name:   "no path",
			Expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if c.Temp {
				p, err := TempFile(c.Dir, "unitest")

				_ = p.Close()
				defer func(name string) {
					_ = os.Remove(name)
				}(p.Name())

				assert.Nil(err)
				assert.Equal(c.Expect, IsExist(p.Name()))
			} else {
				assert.Equal(c.Expect, IsExist(c.Path))
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Path   string
		Expect bool
		Dir    string
		Temp   bool
	}{
		{
			Name:   "dir path",
			Path:   "/tmp",
			Expect: false,
			Temp:   false,
		},
		{
			Name:   "file path",
			Dir:    "/tmp",
			Expect: true,
			Temp:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if c.Temp {
				p, err := TempFile(c.Dir, "unites")

				_ = p.Close()
				defer func(name string) {
					_ = os.Remove(name)
				}(p.Name())

				assert.Nil(err)
				assert.Equal(c.Expect, IsFile(p.Name()))
			} else {
				assert.Equal(c.Expect, IsFile(c.Path))
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Path   string
		Expect bool
		Dir    string
		Temp   bool
	}{
		{
			Name:   "dir path",
			Path:   "/tmp",
			Expect: true,
			Temp:   false,
		},
		{
			Name:   "file path",
			Dir:    "/tmp",
			Expect: false,
			Temp:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if c.Temp {
				p, err := TempFile(c.Dir, "unites")

				_ = p.Close()
				defer func(name string) {
					_ = os.Remove(name)
				}(p.Name())

				assert.Nil(err)
				assert.Equal(c.Expect, IsDir(p.Name()))
			} else {
				assert.Equal(c.Expect, IsDir(c.Path))
			}
		})
	}
}

func TestIsAbs(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Path   string
		Expect bool
	}{
		{
			Name:   "abs-path not exist",
			Path:   "/test/file",
			Expect: false,
		},
		{
			Name:   "abs-path exist",
			Path:   "/tmp",
			Expect: true,
		},
		{
			Name:   "not abs-path",
			Path:   "apple",
			Expect: false,
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			assert.Equal(c.Expect, IsAbs(c.Path))
		})
	}
}

func TestRWFileByte(t *testing.T) {
	var (
		dir     = "/tmp"
		pattern = "_*_"
		text    = "1"
		assert  = A.New(t)
	)

	f, err := TempFile(dir, pattern)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	defer func(name string) {
		_ = os.Remove(name)
	}(f.Name())

	assert.Nil(err)
	assert.NotNil(f)

	err = WriteFileByte(f.Name(), []byte(text), 0666)
	assert.Nil(err)

	b, err := ReadFileByte(f.Name())
	assert.Nil(err)
	assert.Equal(b, []byte(text))
}

func TestRWFileString(t *testing.T) {
	var (
		dir     = "/tmp"
		pattern = "_*_"
		text    = "1"
		assert  = A.New(t)
	)

	f, err := TempFile(dir, pattern)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	defer func(name string) {
		_ = os.Remove(name)
	}(f.Name())

	assert.Nil(err)
	assert.NotNil(f)

	err = WriteFile(f.Name(), text, 0666)
	assert.Nil(err)

	s, err := ReadFile(f.Name())
	assert.Nil(err)
	assert.Equal(s, text)
}

func TestRemoveFile(t *testing.T) {
	var (
		dir     = "/tmp"
		pattern = "_*_"
		assert  = A.New(t)
	)

	f, _ := TempFile(dir, pattern)
	_ = f.Close()

	err := Remove(f.Name(), false)
	assert.Nil(err)

	_, err = os.Stat(f.Name())
	assert.False(!errors.Is(err, os.ErrNotExist))
}

func TestWriteTempFile(t *testing.T) {
	var (
		dir     = "/tmp"
		pattern = "_*_"
		text    = "1"
		assert  = A.New(t)
	)

	f, _ := WriteTempFile(dir, pattern, text)
	assert.NotNil(f)

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(f)

}

func TestWriteTempDir(t *testing.T) {
	var (
		dir         = "/tmp"
		dirPattern  = "_*_"
		filePattern = "_txt_"
		text        = "1"
		assert      = A.New(t)
	)

	dirPath, filePath, err := WriteTempDir(dir, dirPattern, filePattern, text)
	assert.True(strings.HasPrefix(filePath, dirPath))
	assert.Nil(err)

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(dirPath)
}
