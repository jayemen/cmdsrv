package cmdutil

import (
	"io/ioutil"
	"os/exec"
	"time"
)

// CmdCache holds a command, its arguments, its cached output, and the last time it was run.
type CmdCache struct {
	last    time.Time
	cache   []byte
	Command string
	Args    []string
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

// Run executes the configured command, returning standard output. If the last run was less than maxAge ago, then instead returns the cached result from the previous run.
func (cmd *CmdCache) Run(maxAge time.Duration) ([]byte, error) {
	if cmd.last.Add(maxAge).After(time.Now()) {
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
