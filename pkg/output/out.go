package output

import (
	"github.com/alexisvisco/gwd/pkg/vars"
)

const JsonError = `{"error": "unable to output json from this command"}`

type CommandOutput interface {
	Human()
	HumanVerbose()
	Json()
}

func Print(output CommandOutput) {
	if vars.OutputJson {
		output.Json()
		return
	}
	if vars.OutputVerbose {
		output.HumanVerbose()
		return
	}

	output.Human()
}

func Error(err error) {
	if err != nil {
		Print(&errOutput{err.Error()})
	}
}
