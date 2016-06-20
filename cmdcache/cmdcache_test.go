package cmdcache

import (
	"testing"
	"time"

	"github.com/jayemen/cmdsrv/testutil"
)

func TestExecuteCommand(t *testing.T) {
	util := testutil.Wrap(t)
	cmd := New(1*time.Second, "echo", "-n", "this is a test")
	output, err := cmd.Run()
	util.AssertNil(err)
	util.AssertEqual("this is a test", string(output))
}
