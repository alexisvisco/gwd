package packages

import (
	"path/filepath"

	"github.com/alexisvisco/gwd/pkg/vars"
	"github.com/pkg/errors"
)

// ImportPath is a type that represents a package name.
// A package name is the module name suffixed by the path of the package.
type ImportPath string

// GetImportPathFromPath try to get the import path from a path.
// moduleName is the name of the module in which the path is.
// modulePath is the path of the module in which the path is.
// Example:  the current path here is pkg/diff/packages/types.go it will give github.com/alexisvisco/gwd/pkg/diff/packages/types
func GetImportPathFromPath(moduleName, modulePath, path string, isDir bool) (ImportPath, error) {
	packageNameValue := ""
	if isDir {
		packageNameValue = path
	} else {
		packageNameValue = filepath.Dir(path)
	}

	if moduleName, ok := vars.ModulePathToModuleName[packageNameValue]; ok {
		return ImportPath(moduleName), nil
	}

	rel, err := filepath.Rel(modulePath, packageNameValue)
	if err != nil {
		return "", errors.Wrapf(err, "unable to retrieve package name from path %q with module %q and module path %q", packageNameValue, moduleName, modulePath)
	}

	return ImportPath(filepath.Join(moduleName, rel)), nil
}
