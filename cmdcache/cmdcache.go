package cmdcache

import (
	"io/ioutil"
	"os/exec"
	"time"
)

// Cmd is a runnable command with cached output.
type Cmd struct {
	Command string
	Args    []string
	lastRun time.Time
	maxAge  time.Duration
	cache   []byte
}

// New returns a CmdCache initialized with the specified command and arguments.
func New(maxAge time.Duration, command string, args ...string) *Cmd {
	return &Cmd{
		Command: command,
		Args:    args,
		lastRun: time.Time{},
		maxAge:  maxAge,
		cache:   nil,
	}
}

// Run gets the output of the command, using the cached value if the last run was less than maxAge ago.
func (c *Cmd) Run() ([]byte, error) {
	if c.lastRun.Add(c.maxAge).After(time.Now()) {
		return c.cache, nil
	}

	output, err := outputOf(exec.Command(c.Command, c.Args...))

	if err != nil {
		return nil, err
	}

	c.lastRun = time.Now()
	c.cache = output
	return c.cache, nil
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
