package parsing

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/mod/modfile"
	"io/ioutil"
)

// GetWorkspaceModulePaths retrieve the module paths from a workspace file
func GetWorkspaceModulePaths(path string) ([]string, error) {
	goWorkFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read go.work file %q", path)
	}

	work, err := modfile.ParseWork(path, goWorkFile, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse go.work file %q", path)
	}

	return lo.Map(work.Use, func(t *modfile.Use, i int) string {
		return t.Path
	}), nil
}
