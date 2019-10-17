package vars

import "gopkg.in/src-d/go-git.v4"

// Repository is the current git repository
var Repository *git.Repository

// ModuleName represent the module path used to known which import is
// part of the current project
var ModuleName string

// OutputJson is used to print result of command as json
var OutputJson bool

// OutputVerbose is used to print result of command as a verbose result
// with more detailed information
var OutputVerbose bool
