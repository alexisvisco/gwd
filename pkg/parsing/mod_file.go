package parsing

import (
	"github.com/pkg/errors"
	"golang.org/x/mod/modfile"
	"io/ioutil"
)

// GetModuleName retrieve the module name from a go file
func GetModuleName(path string) (string, error) {
	goModFile, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read go.mod file %q", path)
	}

	work, err := modfile.Parse(path, goModFile, nil)
	if err != nil {
		return "", errors.Wrapf(err, "unable to parse go.mod file %q", path)
	}

	return work.Module.Mod.Path, nil
}
