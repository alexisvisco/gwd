package output

import (
	"encoding/json"
	"fmt"
	"strings"
)

type errOutput struct {
	Error string `json:"error"`
}

func (e errOutput) Human() {
	fmt.Println("Error:", e.Error)
}

func (e errOutput) HumanVerbose() {
	fmt.Println("Error:", e.Error)
}

func (e errOutput) Json() {
	j, err := json.Marshal(e)
	if err != nil {
		fmt.Println(JsonError)
		return
	}

	fmt.Println(string(j))
}

type String string

func (s String) Human() {
	fmt.Println(s)
}

func (s String) HumanVerbose() {
	fmt.Println(s)
}

func (s String) Json() {
	j, err := json.Marshal(s)
	if err != nil {
		fmt.Println(JsonError)
		return
	}

	fmt.Println(string(j))
}

type StringArray []string

func (p StringArray) Human() {
	fmt.Println(strings.Join(p, "\n"))
}

func (p StringArray) HumanVerbose() {
	p.Human()
}

func (p StringArray) Json() {
	marshal, err := json.Marshal(p)
	if err != nil {
		fmt.Println(JsonError)
		return
	}

	fmt.Println(string(marshal))
}
