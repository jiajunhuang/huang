package main

import (
	"flag"

	"github.com/jiajunhuang/huang/cli"
	"github.com/jiajunhuang/huang/master"
	"github.com/jiajunhuang/huang/worker"
)

var (
	runAsMaster = flag.Bool("master", false, "run as master")
	runAsWorker = flag.Bool("worker", false, "run as worker")
)

func main() {
	flag.Parse()

	if *runAsMaster {
		master.Main()
	} else if *runAsWorker {
		worker.Main()
	} else { // by default, we're running as cli
		cli.Main()
	}
}
