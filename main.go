package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jayemen/cmdsrv/cmdutil"
)

func parseArgs() (cmdCache *cmdutil.CmdCache, listen string) {
	cmdFlag := flag.String("cmd", "", "command to execute")
	argsFlag := flag.String("args", "", "command-line arguments")
	listenFlag := flag.String("listen", ":7777", "listen configuration")
	maxAgeFlag := flag.Int("cache-time", 5, "seconds the command output can be cached")

	flag.Parse()

	if *cmdFlag == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	argSlice := []string{}

	if *argsFlag != "" {
		argSlice = strings.Split(*argsFlag, " ")
	}

	maxAge := time.Duration(*maxAgeFlag) * time.Second
	cmdCache = cmdutil.MakeCmdCache(maxAge, *cmdFlag, argSlice...)
	listen = *listenFlag
	return
}

type cmdResponse struct {
	output []byte
	err    error
}

type cmdServer struct {
	cmd *cmdutil.CmdCache
	ch  chan chan cmdResponse
}

func makeServer(cmd *cmdutil.CmdCache) *cmdServer {
	server := &cmdServer{
		cmd: cmd,
		ch:  make(chan chan cmdResponse),
	}

	return server
}

func (s *cmdServer) start() {
	for {
		reply := <-s.ch
		output, err := s.cmd.Run()
		reply <- cmdResponse{output, err}
	}
}

func (s *cmdServer) runCmd() (output []byte, err error) {
	reply := make(chan cmdResponse)
	s.ch <- reply
	response := <-reply
	return response.output, response.err
}

func main() {
	cmd, listenSpec := parseArgs()
	server := makeServer(cmd)
	go server.start()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output, err := server.runCmd()

		if err != nil {
			_, innerErr := w.Write([]byte(err.Error()))

			if innerErr != nil {
				log.Fatal(err)
			}

			log.Print("Err running command: ", err)
			return
		}

		_, err = w.Write(output)
		if err != nil {
			log.Fatal(err)
		}
	})

	log.Fatal(http.ListenAndServe(listenSpec, mux))
}
