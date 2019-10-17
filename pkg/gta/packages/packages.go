package packages

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"gopkg.in/src-d/go-git.v4/utils/merkletrie"

	"github.com/alexisvisco/gta/pkg/gta"
)

type Packages map[string]*Details

func NewPackages() Packages {
	return Packages{}
}

type Details struct {
	Files      map[string]merkletrie.Action `json:"files"`
	ImportedBy map[string]int               `json:"imported_by"`
}

func (p Packages) addModifiedPackage(packageName, file string, action merkletrie.Action) {
	details, ok := p[packageName]
	if !ok {
		details = &Details{Files: make(map[string]merkletrie.Action), ImportedBy: make(map[string]int)}
		p[packageName] = details
	}

	details.Files[file] = action
}

func (p Packages) addImportModifiedPackage(importPackageName, packageName string) {
	details, ok := p[importPackageName]
	if ok {
		counter, ok := details.ImportedBy[packageName]
		if !ok {
			details.ImportedBy[packageName] = 1
		} else {
			details.ImportedBy[packageName] = counter + 1
		}
	}
}

func (p Packages) Human() {
	affectedPackages := make([]string, 0, len(p)*2)

	for packageName := range p {
		affectedPackages = append(affectedPackages, packageName)
		detail := p[packageName]

		for importByPackageName := range detail.ImportedBy {
			affectedPackages = append(affectedPackages, importByPackageName)
		}
	}

	sort.Strings(affectedPackages)
	buffer := bytes.NewBuffer(nil)
	for _, packageName := range affectedPackages {
		buffer.WriteString(packageName + "\n")
	}

	fmt.Println(buffer.String())
}

func (p Packages) HumanVerbose() {
	buffer := bytes.NewBuffer(nil)

	affectedPackages := make([]string, 0, len(p))
	for packageName := range p {
		affectedPackages = append(affectedPackages, packageName)
	}

	sort.Strings(affectedPackages)
	for _, packageName := range affectedPackages {
		buffer.WriteString(packageName + "\n")
		details := p[packageName]

		//----

		files := make([]string, len(details.Files))
		i := 0
		for file := range details.Files {
			files[i] = file
			i++
		}

		sort.Strings(files)
		buffer.WriteString("  files affected:\n")
		for _, file := range files {
			buffer.WriteString(fmt.Sprintf("    %s %s\n", actionToSymbol(details.Files[file]), file))
		}

		//----
		if len(details.ImportedBy) > 0 {
			buffer.WriteString("  imported by:\n")

			packagesList := make([]string, len(details.ImportedBy))
			i := 0
			for packageName := range details.ImportedBy {
				packagesList[i] = packageName
				i++
			}

			for _, packageName := range packagesList {
				buffer.WriteString(fmt.Sprintf("    %s (%d file(s))\n", packageName, details.ImportedBy[packageName]))
			}
		}
	}

	fmt.Println(buffer.String())
}

func actionToSymbol(action merkletrie.Action) string {
	switch action {
	case merkletrie.Delete:
		return "D"
	case merkletrie.Insert:
		return "A"
	case merkletrie.Modify:
		return "M"
	}
	return "?"
}

func (p Packages) Json() {
	jsonIndented, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println(gta.JsonError)
		return
	}

	fmt.Println(string(jsonIndented))
}
