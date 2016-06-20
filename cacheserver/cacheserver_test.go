package cacheserver

import (
	"testing"
	"time"

	"github.com/jayemen/cmdsrv/cmdcache"
	"github.com/jayemen/cmdsrv/testutil"
)

func TestServer(t *testing.T) {
	util := testutil.Wrap(t)
	cmd := cmdcache.New(1*time.Second, "echo", "-n", "this is a test")
	server := New(cmd)
	go server.Start()
	output, err := server.Run()
	util.AssertNil(err)
	util.AssertEqual("this is a test", string(output))
}
