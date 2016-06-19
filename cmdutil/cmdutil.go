package cmdutil

import (
	"io/ioutil"
	"os/exec"
	"time"
)

// CmdCache is a runnable command with cached stdout.
type CmdCache struct {
	Command string
	Args    []string
	last    time.Time
	maxAge  time.Duration
	cache   []byte
}

// MakeCmdCache returns a CmdCache initialized with the specified command and arguments.
func MakeCmdCache(maxAge time.Duration, command string, args ...string) *CmdCache {
	return &CmdCache{
		last:    time.Time{},
		Command: command,
		Args:    args,
		maxAge:  maxAge,
		cache:   nil,
	}
}

// Run gets the output of the command, using the cached value if the last run was less than maxAge ago.
func (cmd *CmdCache) Run() ([]byte, error) {
	if cmd.last.Add(cmd.maxAge).After(time.Now()) {
		return cmd.cache, nil
	}

	output, err := runCmd(cmd.Command, cmd.Args...)

	if err != nil {
		return nil, err
	}

	cmd.last = time.Now()
	cmd.cache = output
	return cmd.cache, nil
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
