package packages

import (
	"github.com/alexisvisco/gwd/pkg/utils"
)

// Modified represent a list of changes for a package
// The key is the package import path
// The value is a struct containing the list of files affected and the list of packages that import this package
type Modified map[ImportPath]*Details

func NewChanges() Modified {
	return Modified{}
}

type Details struct {
	// Files is a map of file name to action (insert, delete, modified in a git sense)
	// The action is a string that can be converted to merkletrie.Action
	// The key is the file path
	Files              utils.StringSet            `json:"files"`
	ImportedImportPath map[ImportPath]*ImportedBy `json:"imported_by"`
}

type ImportedBy struct {
	ModuleName string `json:"module_name"`
	ModulePath string `json:"module_path"`
	Counter    int    `json:"counter"`
}

// AddModifiedPackage will add a file that have been modified
func (p Modified) AddModifiedPackage(importPath ImportPath, file string) {
	details, ok := p[importPath]
	if !ok {
		details = &Details{Files: utils.NewStringSet(file), ImportedImportPath: make(map[ImportPath]*ImportedBy)}
		p[importPath] = details
	}
}

// AddImportPathWhichImportModifiedPackage will add the import path that import a package which have been modified
// moduleName, modulePath are the module name and path of the import path that import a modified package
// importedByImportPath is the package import path which import a package modified
// importPathChanged is the package import path that have been modified
func (p Modified) AddImportPathWhichImportModifiedPackage(moduleName, modulePath string, importPathChanged, importedByImportPath ImportPath) {
	details, ok := p[importPathChanged]
	if ok {
		importedBy, ok := details.ImportedImportPath[importedByImportPath]
		if !ok {
			importedBy = &ImportedBy{
				ModuleName: moduleName,
				ModulePath: modulePath,
				Counter:    0,
			}
			details.ImportedImportPath[importedByImportPath] = importedBy
		}
		importedBy.Counter++
	}
}
