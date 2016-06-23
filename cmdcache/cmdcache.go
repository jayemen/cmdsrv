package cmdcache

import (
	"bytes"
	"os/exec"
	"time"
)

// Cmd is a runnable command with cached output.
type Cmd struct {
	Command string
	Args    []string
	lastRun time.Time
	maxAge  time.Duration
	cache   bytes.Buffer
}

// New returns a CmdCache initialized with the specified command and arguments.
func New(maxAge time.Duration, command string, args ...string) *Cmd {
	return &Cmd{
		Command: command,
		Args:    args,
		lastRun: time.Time{},
		maxAge:  maxAge,
		cache:   bytes.Buffer{},
	}
}

// Run gets the output of the command, using the cached value if the last run was less than maxAge ago.
func (c *Cmd) Run() ([]byte, error) {
	if c.lastRun.Add(c.maxAge).After(time.Now()) {
		return c.cache.Bytes(), nil
	}

	err := c.refreshCache()

	if err != nil {
		return nil, err
	}

	c.lastRun = time.Now()
	return c.cache.Bytes(), nil
}

func (c *Cmd) refreshCache() error {
	cmd := exec.Command(c.Command, c.Args...)

	reader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	c.cache.Reset()

	_, err = c.cache.ReadFrom(reader)
	if err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
