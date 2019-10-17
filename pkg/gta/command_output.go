package gta

import "github.com/alexisvisco/gta/pkg/gta/vars"

const JsonError = `{"error": "unable to output json from this command"}`

type CommandOutput interface {
	Human()
	HumanVerbose()
	Json()
}

func Output(output CommandOutput) {
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
