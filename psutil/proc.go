package psutil

import (
	"context"
	"fmt"
	"github.com/lafrinte/nops/fs"
	psutil "github.com/shirou/gopsutil/v3/process"
	"strconv"
	"strings"
	"time"
)

func killProcess(p *psutil.Process) error {
	children, err := p.Children()
	if err != nil {
		return fmt.Errorf("failed to get child-process: pid -> %d, err -> %s", p.Pid, err)
	}

	for _, c := range children {
		if err := c.Kill(); err != nil {
			return fmt.Errorf("failed kill child process: pid -> %d, ppid -> %d, err -> %s", c.Pid, p.Pid, err)
		}
	}

	return p.Kill()
}

/*
	GetProcStat

Description:

will try to get process's state if it is existed. when process does not exist, the function will return an error.

Returns:

	[]string: Status returns the process status. Return value could be one of these.
	          R: Running S: Sleep T: Stop I: Idle Z: Zombie W: Wait L: Lock The character
*/
func GetProcStat(pid int) ([]string, error) {
	p, err := psutil.NewProcess(int32(pid))
	if err != nil {
		return []string{}, fmt.Errorf("pid not exist: pid -> %d", pid)
	}

	return p.Status()
}

/*
	TerminateWithPidFile

Description:

will try reading pid file and terminate the process.

Args:

	timeout: max second waiting termination.
	force: whether using signal.SIGKILL when signal.SIGTERM failed.
*/
func TerminateProcess(pid int, timeout int, force bool) error {
	if timeout == 0 {
		timeout = 3
	}

	p, err := psutil.NewProcess(int32(pid))
	if err != nil {
		// pid not exist, return
		return nil
	}

	// pid exist, send signal.SIGTERM to process
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan int, 1)
	// TerminateWithContext will return directly when process get signal.
	go func(ctx context.Context, in chan int) {
		_ = p.TerminateWithContext(ctx)

		for {
			if ok, _ := p.IsRunning(); !ok {
				// process terminated
				in <- 0
				break
			}
		}
	}(ctx, ch)

	select {
	case <-ch:
		// signal.TERM succeed
		return nil
	case <-time.After(time.Duration(timeout) * time.Second):
		// timeout handler
		switch force {
		case true:
			if err := killProcess(p); err != nil {
				return fmt.Errorf("terminate and kill failed: %s", err)
			}
		case false:
			return fmt.Errorf("terminate failed: pid -> %d", pid)
		}
	case <-ctx.Done():
		return nil
	}

	return nil
}

/*
	TerminateWithPidFile

Description:

will try reading pid file and terminate the process.

Args:

	path: pid file path.
	timeout: max second waiting termination.
	force: whether using signal.SIGKILL when signal.SIGTERM failed.
*/
func TerminateWithPidFile(path string, timeout int, force bool) error {
	var (
		pid int
	)

	p, err := fs.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to open and read file: path -> %s", path)
	}

	pid, err = strconv.Atoi(strings.TrimSpace(p))
	if err != nil {
		return fmt.Errorf("pid -> %s reading from file -> %s same not a valid int", p, path)
	}

	return TerminateProcess(pid, timeout, force)
}
