package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jayemen/cmdsrv/cacheserver"
	"github.com/jayemen/cmdsrv/cmdcache"
)

func parseArgs() (c *cmdcache.Cmd, listen string) {
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
	c = cmdcache.New(maxAge, *cmdFlag, argSlice...)
	listen = *listenFlag
	return
}

func main() {
	cmd, listenSpec := parseArgs()
	server := cacheserver.New(cmd)
	go server.Start()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output, err := server.Run()

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
