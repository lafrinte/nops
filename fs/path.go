package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func GetExecDir() string {
	dir, _ := os.Executable()
	return filepath.Dir(dir)
}

func GetDirName(p string) string {
	return filepath.Dir(p)
}

func GetBaseName(p string) string {
	return filepath.Base(p)
}

func GetStat(p string) (os.FileInfo, error) {
	return os.Stat(p)
}

func IsExist(p string) bool {
	_, err := GetStat(p)
	return !errors.Is(err, os.ErrNotExist)
}

func IsFile(p string) bool {
	s, err := GetStat(p)
	if err != nil {
		return false
	}
	return s.Mode().IsRegular()
}

func IsDir(p string) bool {
	s, err := GetStat(p)
	if err != nil {
		return false
	}

	return s.Mode().IsDir()
}

func IsAbs(p string) bool {
	return IsExist(p) && filepath.IsAbs(p)
}

func TempFile(dir string, pattern string) (*os.File, error) {
	return os.CreateTemp(dir, pattern)
}

func TempDir(dir string, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

func ReadFileByte(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func ReadFile(p string) (string, error) {
	buf, err := ReadFileByte(p)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func WriteFileByte(path string, buf []byte, perm os.FileMode) error {
	dir := GetDirName(path)
	if !IsExist(dir) {
		if err := os.MkdirAll(dir, 0750); err != nil {
			return fmt.Errorf("dir %s does not exist and has no permission to create it", dir)
		}
	}

	if IsExist(dir) && !IsDir(dir) {
		return fmt.Errorf("dir %s does not a valid directory", dir)
	}

	if !IsExist(path) {
		if _, err := os.Create(path); err != nil {
			return fmt.Errorf("%s does not exist and has not permission to create it", path)
		}
	}

	if !IsFile(path) {
		return fmt.Errorf("%s is not a valid directory", path)
	}

	return os.WriteFile(path, buf, perm)
}

func WriteFile(path string, s string, perm os.FileMode) error {
	return WriteFileByte(path, []byte(s), perm)
}

func WriteTempFile(dir string, pattern string, s string) (string, error) {
	f, err := TempFile(dir, pattern)
	if err != nil {
		return "", err
	}

	_, err = f.Write([]byte(s))
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func WriteTempDir(dir string, dirPattern string, filePattern string, s string) (string, string, error) {
	dir, err := TempDir(dir, dirPattern)
	if err != nil {
		return "", "", err
	}

	path, err := WriteTempFile(dir, filePattern, s)
	if err != nil {
		return dir, "", err
	}

	return dir, path, nil
}

func Remove(path string, force bool) error {
	if force {
		return os.RemoveAll(path)
	}

	return os.Remove(path)
}

func Walk(path string) ([]string, error) {
	var ps []string

	if !IsDir(path) {
		return []string{}, nil
	}

	err := filepath.Walk(path,
		func(p string, i os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			ps = append(ps, p)

			return nil
		},
	)

	return ps, err
}

func WalkRegex(path string, pattern string) ([]string, error) {
	var (
		ps []string
		r  *regexp.Regexp
	)

	r, err := regexp.Compile(pattern)
	if err != nil {
		return ps, err
	}

	tps, err := Walk(path)
	if err != nil {
		return ps, err
	}

	for _, p := range tps {
		if r.MatchString(p) {
			ps = append(ps, p)
		}
	}

	return ps, nil
}

func WalkGlob(path string) ([]string, error) {
	return filepath.Glob(path)
}

func WalkGlobs(paths ...string) ([]string, error) {
	var ps []string

	for _, path := range paths {
		p, err := WalkGlobs(path)
		if err != nil {
			return ps, err
		}

		ps = append(ps, p...)
	}

	return ps, nil
}
