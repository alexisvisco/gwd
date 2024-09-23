package main

import (
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

	"github.com/alexisvisco/gwd/commands"
)

func main() {
	if os.Getenv("GWD_PROFILE") != "" {
		f, err := os.Create("gwd.cpu.pprof")
		if err != nil {
			panic(err)
		}

		defer func() {
			if err := pprof.StartCPUProfile(f); err != nil {
				panic(err)
			}
		}()
	}

	commands.Execute()
}
