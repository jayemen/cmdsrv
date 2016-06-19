package cmdutil

import (
	"io/ioutil"
	"os/exec"
	"sync"
	"time"
)

// CmdCache is a runnable command with cached stdout.
type CmdCache struct {
	Command string
	Args    []string
	last    time.Time
	cache   []byte
	lock    sync.Mutex
}

// MakeCmdCache returns a CmdCache initialized with the specified command and arguments.
func MakeCmdCache(command string, args ...string) *CmdCache {
	return &CmdCache{
		last:    time.Time{},
		Command: command,
		Args:    args,
		cache:   nil,
	}
}

// Run gets the output of the command, using the cached value if the last run was less than maxAge ago.
func (cmd *CmdCache) Run(maxAge time.Duration) ([]byte, error) {
	cmd.lock.Lock()

	if cmd.last.Add(maxAge).After(time.Now()) {
		cp := make([]byte, len(cmd.cache))
		copy(cp, cmd.cache)
		cmd.lock.Unlock()
		return cp, nil
	}

	cmd.lock.Unlock()

	output, err := runCmd(cmd.Command, cmd.Args...)

	if err != nil {
		return nil, err
	}

	cmd.lock.Lock()
	cmd.last = time.Now()
	cmd.cache = output
	cp := make([]byte, len(cmd.cache))
	copy(cp, cmd.cache)
	cmd.lock.Unlock()
	return cp, nil
}

func outputOf(cmd *exec.Cmd) ([]byte, error) {
	reader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	output, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return output, nil
}

func runCmd(program string, args ...string) ([]byte, error) {
	output, err := outputOf(exec.Command(program, args...))
	if err != nil {
		return nil, err
	}

	return output, nil
}
