package vars

import "gopkg.in/src-d/go-git.v4"


// Repository is the current git repository
var Repository *git.Repository

// ModuleName represent the module path used to known which import is
// part of the current project
var ModuleName = ""
