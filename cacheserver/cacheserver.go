package cacheserver

import "github.com/jayemen/cmdsrv/cmdcache"

// Server serializes access to a cmdcache.Cmd, in a threadsafe manner.
type Server struct {
	cmd *cmdcache.Cmd
	ch  chan chan response
}

type response struct {
	output []byte
	err    error
}

// New returns a new server instance.
func New(cmd *cmdcache.Cmd) *Server {
	server := &Server{
		cmd: cmd,
		ch:  make(chan chan response),
	}

	return server
}

// Start begins a blocking loop that handles server requests.
func (s *Server) Start() {
	for {
		reply := <-s.ch
		output, err := s.cmd.Run()
		reply <- response{output, err}
	}
}

// Run executes the configured command, and returns its result. This is thread-safe.
func (s *Server) Run() (output []byte, err error) {
	reply := make(chan response)
	s.ch <- reply
	response := <-reply
	return response.output, response.err
}
