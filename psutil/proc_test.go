package psutil

import (
	"fmt"
	psutil "github.com/shirou/gopsutil/v3/process"
	A "github.com/stretchr/testify/assert"
	"nops/fs"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

var TempDir string

func MockTestDir() {
	dir, err := fs.TempDir("/tmp", "proc_test")
	if err != nil {
		panic(err)
	}

	TempDir = dir
}

func MockCleanTestDir() {
	_ = fs.Remove(TempDir, true)
}

func MockShellProcessAndGetPid(script string, out chan int) {
	f, _ := fs.TempFile(TempDir, "_*_")
	_, _ = f.WriteString(script)

	cmd := exec.Command("/bin/bash", f.Name())
	if err := cmd.Start(); err == nil {
		out <- cmd.Process.Pid
		_ = cmd.Wait()
	}
}

func TestKillProcess(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name   string
		Script string
		Nil    bool
	}{
		{
			Name: "killing main-process succeed",
			Script: `#/usr/bin/env bash
sleep 100
exit 0`,
			Nil: true,
		},
		{
			Name: "killing main-process and its' children",
			Script: `#/usr/bin/env bash
for i in $(seq 1 2|xargs); do
  (sleep 100) &
done
wait
exit 0`,
			Nil: true,
		},
	}

	MockTestDir()

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			in := make(chan int, 1)
			go func() {
				MockShellProcessAndGetPid(c.Script, in)
			}()

		LOOP:
			for {
				select {
				case pid := <-in:
					// process is still running
					p, _ := psutil.NewProcess(int32(pid))

					fmt.Println(p.Status())
					err := killProcess(p)
					switch c.Nil {
					case true:
						assert.Nil(err)
					case false:
						assert.NotNil(err)
					}
					break LOOP
				default:
					time.Sleep(time.Millisecond * 100)
				}
			}
		})
	}

	MockCleanTestDir()
}

func TestTerminateProcess(t *testing.T) {
	assert := A.New(t)

	cases := []struct {
		Name    string
		Pid     int
		Script  string
		Timeout int
		Nil     bool
		Force   bool
	}{
		{
			Name: "process is not exist",
			Pid:  1000000,
			Nil:  true,
		},
		{
			Name: "terminate main-process succeed",
			Script: `#/usr/bin/env bash
sleep 100
exit 0`,
			Nil: true,
		},
		{
			Name: "terminate main-process timeout",
			Script: `#/usr/bin/env bash
trap "echo 'get TERM signal but ignored.'" TERM
sleep 100
exit 0`,
			Nil: false,
		},
		{
			Name: "terminate main-process force",
			Script: `#/usr/bin/env bash
trap "echo 'get TERM or INT signal but ignored.'" TERM INT
sleep 10
exit 0`,
			Force: true,
			Nil:   true,
		},
		{
			Name: "killing main-process and its' children",
			Script: `#/usr/bin/env bash
for i in $(seq 1 2|xargs); do
  (sleep 10) &
done
wait
exit 0`,
			Timeout: 15,
			Nil:     true,
		},
	}

	MockTestDir()

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			in := make(chan int, 1)
			if c.Script != "" && c.Pid == 0 {
				go func() {
					MockShellProcessAndGetPid(c.Script, in)
				}()

			LOOP:
				for {
					select {
					case pid := <-in:
						// process is still running
						err := TerminateProcess(pid, c.Timeout, c.Force)
						switch c.Nil {
						case true:
							assert.Nil(err)
						case false:
							assert.NotNil(err)
						}
						break LOOP
					default:
						time.Sleep(time.Millisecond * 100)
					}
				}
			}

			if c.Pid > 0 {
				switch c.Nil {
				case true:
					assert.Nil(TerminateProcess(c.Pid, c.Timeout, c.Force))
				case false:
					assert.NotNil(TerminateProcess(c.Pid, c.Timeout, c.Force))
				}
			}
		})
	}

	MockCleanTestDir()
}

func TestTerminateWithPidFile(t *testing.T) {
	var (
		assert  = A.New(t)
		pattern = "_*_.pid"
		timeout = 5
		script  = `#/usr/bin/env bash
trap "echo 'get TERM or INT signal but ignored.'" TERM INT
sleep 10
exit 0`
		in = make(chan int, 1)
	)

	MockTestDir()

	go func() {
		MockShellProcessAndGetPid(script, in)
	}()

	select {
	case pid := <-in:
		// process is still running

		path, err := fs.WriteTempFile(TempDir, pattern, strconv.Itoa(pid))
		assert.Nil(err)
		assert.NotNil(path)

		err = TerminateWithPidFile(path, timeout, true)
		assert.Nil(err)
	case <-time.After(time.Duration(timeout+5) * time.Second):
	}

	MockCleanTestDir()
}
