package modules

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/alexisvisco/gwd/pkg/diff/packages"
	"github.com/alexisvisco/gwd/pkg/output"
	"github.com/alexisvisco/gwd/pkg/parsing"
	"github.com/alexisvisco/gwd/pkg/vars"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/samber/lo"
)

type Modules struct {
	Modules []*Changes `json:"modules"`
}

type Changes struct {
	ModuleName                 string              `json:"module_name"`
	ModulePath                 string              `json:"module_path"`
	PackagesModified           packages.Modified   `json:"packages_modified"`
	Reason                     ModuleChangeReason  `json:"reason"`
	ModuleDependencyImportPath packages.ImportPath `json:"module_dependency_import_path,omitempty"`
}

type cache struct {
	modulePathToChanges map[string]*Changes
	modifiedImportPath  map[packages.ImportPath]packages.Modified
}

func FromChanges(changes merkletrie.Changes) (*Modules, error) {
	modules := &Modules{Modules: make([]*Changes, 0, len(vars.ModuleNameToModulePath))}
	cache := &cache{
		modulePathToChanges: make(map[string]*Changes),
		modifiedImportPath:  make(map[packages.ImportPath]packages.Modified),
	}
	err := detectImportPathThatHaveChanged(changes, modules, cache)
	if err != nil {
		return nil, err
	}

	if err := detectImportPathImported(modules, cache); err != nil {
		return nil, err
	}

	return modules, nil
}

type ModuleChangeReason string

const (
	ModuleChangeReasonFile             = ModuleChangeReason("ModuleChangeReasonFile")
	ModuleChangeReasonModuleDependency = ModuleChangeReason("ModuleChangeReasonModuleDependency")
)

func detectImportPathThatHaveChanged(changes merkletrie.Changes, modules *Modules, c *cache) error {
	for _, ch := range changes {
		action, err := ch.Action()
		if err != nil {
			continue
		}

		path := ch.To
		if path.String() == "" {
			path = ch.From
		}

		moduleName, modulePath := getModuleFromFilePath(path.String())

		if moduleName == "" {
			// the changes are not in a go module, so we need to ignore it
			continue
		}

		changes, ok := c.modulePathToChanges[modulePath]
		if !ok {
			changes = &Changes{
				ModuleName:       moduleName,
				ModulePath:       modulePath,
				PackagesModified: packages.NewChanges(),
				Reason:           ModuleChangeReasonFile,
			}
			c.modulePathToChanges[modulePath] = changes
		}

		packageName, err := packages.GetImportPathFromPath(moduleName, modulePath, path.String(), path.IsDir())
		if err != nil {
			return err
		}
		c.modifiedImportPath[packageName] = changes.PackagesModified

		changes.PackagesModified.AddModifiedPackage(packageName, path.String(), action)

	}

	for _, change := range c.modulePathToChanges {
		modules.Modules = append(modules.Modules, change)
	}

	return nil
}

// detectImportPathImported loop over all modules and parse all go files to check whether they import a package
// that have been modified.
func detectImportPathImported(modules *Modules, cache *cache) error {
	for modulePathWhichCanImportModifiedPackage, moduleNameWhichCanImportModifiedPackage := range vars.ModulePathToModuleName {
		err := filepath.WalkDir(modulePathWhichCanImportModifiedPackage, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
				return nil
			}

			importPathWhichImportModifiedPackage, err := packages.GetImportPathFromPath(
				moduleNameWhichCanImportModifiedPackage,
				modulePathWhichCanImportModifiedPackage,
				path,
				false,
			)

			if err != nil {
				return err
			}

			for _, importPathWhichIsModifiedPackage := range parsing.GetImports(path) {
				packageChange, ok := cache.modifiedImportPath[importPathWhichIsModifiedPackage]
				if !ok {
					continue
				}

				packageChange.AddImportPathWhichImportModifiedPackage(
					moduleNameWhichCanImportModifiedPackage,
					modulePathWhichCanImportModifiedPackage,
					importPathWhichIsModifiedPackage,
					importPathWhichImportModifiedPackage,
				)

				changes, ok := cache.modulePathToChanges[modulePathWhichCanImportModifiedPackage]
				if !ok {
					changes = &Changes{
						ModuleName:                 moduleNameWhichCanImportModifiedPackage,
						ModulePath:                 modulePathWhichCanImportModifiedPackage,
						PackagesModified:           packages.NewChanges(),
						Reason:                     ModuleChangeReasonModuleDependency,
						ModuleDependencyImportPath: importPathWhichIsModifiedPackage,
					}
					cache.modulePathToChanges[modulePathWhichCanImportModifiedPackage] = changes
					modules.Modules = append(modules.Modules, changes)
				}
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func getModuleFromFilePath(filePath string) (moduleName, modulePath string) {
	for modulePath, moduleName = range vars.ModulePathToModuleName {
		if strings.Contains(filePath, modulePath) {
			return
		}
	}
	return "", ""
}

func (m Modules) Human() {
	if len(m.Modules) == 0 {
		return
	}

	fmt.Println(strings.Join(lo.Map(m.Modules, func(t *Changes, i int) string {
		return t.ModulePath
	}), " "))
}

func (m Modules) HumanVerbose() {
	for _, mod := range m.Modules {
		fmt.Printf(
			"%q%s\n", mod.ModuleName,
			lo.Ternary(mod.Reason == ModuleChangeReasonModuleDependency,
				fmt.Sprintf(" no files changed but the package imports %q", mod.ModuleDependencyImportPath),
				":"))
		for _, ch := range mod.PackagesModified {
			for s := range ch.Files {
				fmt.Printf(" - %q\n", s)
				if len(ch.ImportedImportPath) > 0 {
					fmt.Printf("   imported by:\n")
					for importPath, detail := range ch.ImportedImportPath {
						fmt.Printf("   ∟ %q (%d times) module %q\n", importPath, detail.Counter, detail.ModuleName)
					}
				}
			}
		}
	}
}

func (m Modules) Json() {
	marshal, err := json.Marshal(m)
	if err != nil {
		fmt.Println(output.JsonError)
		return
	}
	fmt.Println(string(marshal))
}
