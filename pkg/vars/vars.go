package vars

import "github.com/go-git/go-git/v5"

// Repository is the current git repository
var Repository *git.Repository

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
