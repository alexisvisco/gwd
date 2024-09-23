package vars

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// FilesChanged represent the list of files changed
var FilesChanged []string

// GoWorkFileName are the list of modules described by the go.mod
var GoWorkFileName string

// GoWorkModulePaths is the list of path of different modules
var GoWorkModulePaths []string

// ModulePathToModuleName is a map where the key is the module path and the value the module name
var ModulePathToModuleName = map[string]string{}

// ModuleNameToModulePath is a map where the key is the module name and the value the module path
var ModuleNameToModulePath = map[string]string{}

// OutputJson is used to print result of command as json
var OutputJson bool

// OutputVerbose is used to print result of command as a verbose result
// with more detailed information
var OutputVerbose bool

var FilesChangedFromStdin bool
var FilesChangedFromFile string

func LoadFilesChanged() (err error) {
	if FilesChangedFromStdin {
		FilesChanged, err = readFilesChangedFromStdin()
		if err != nil {
			return errors.Wrap(err, "unable to read stdin")
		}
	} else if FilesChangedFromFile != "" {
		content, err := ioutil.ReadFile(FilesChangedFromFile)
		if err != nil {
			return errors.Wrapf(err, "unable to read file %s", FilesChangedFromFile)
		}
		FilesChanged = strings.Split(string(content), "\n")
	}

	return errors.Wrap(err, "unable to load files changed: you must choose between stdin or file check with --help")
}

func readFilesChangedFromStdin() ([]string, error) {
	reader := bufio.NewReader(os.Stdin)

	buf := bytes.NewBuffer(nil)
	for {
		data := make([]byte, 4<<20)
		amount, err := reader.Read(data)
		buf.WriteString(string(data[:amount]))
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "failed to read files changed")
		}
	}

	return strings.Split(buf.String(), "\n"), nil
}
