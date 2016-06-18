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

func parseArgs() (cmdCache cmdutil.CmdCache, listen string, maxAge time.Duration) {
	cmdFlag := flag.String("cmd", "", "the command name")
	argsFlag := flag.String("args", "", "the command arguments")
	listenFlag := flag.String("listen", ":7777", "the listen config")
	maxAgeFlag := flag.Int("cache-time", 5, "the number of seconds the command output can be cached")

	flag.Parse()

	if *cmdFlag == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	cmdCache = cmdutil.MakeCmdCache(*cmdFlag, strings.Split(*argsFlag, " ")...)
	maxAge = time.Duration(*maxAgeFlag) * time.Second
	listen = *listenFlag
	return
}

func main() {
	cmd, listenSpec, maxAge := parseArgs()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output, err := cmd.Run(maxAge)
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
