package main

import (
	"github.com/alexisvisco/gwd/commands"
	"os"
	"runtime/pprof"
)
import _ "net/http/pprof"

func main() {
	f, err := os.Create("gwd.cpu.pprof")
	if err != nil {
		panic(err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}

	commands.Execute()

	pprof.StopCPUProfile()
}
