package packages

import (
	"os"
	"path"
	"strings"

	"github.com/MichaelTJones/walk"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"

	"github.com/alexisvisco/gta/pkg/gta/parsing"
	"github.com/alexisvisco/gta/pkg/gta/vars"
)

func FromChanges(changes merkletrie.Changes) Packages {
	packages := NewPackages()

	setModifiedPackages(changes, packages)
	setImportModifiedPackages(packages)

	return packages
}

func setModifiedPackages(changes merkletrie.Changes, packages Packages) {
	for _, ch := range changes {
		action, err := ch.Action()
		if err != nil {
			continue
		}

		path := ch.To
		if path.String() == "" {
			path = ch.From
		}

		if path.IsDir() {
			pathString := path.String()
			packages.addModifiedPackage(pathString, pathString, action)
			continue
		}

		if len(path) == 1 {
			packages.addModifiedPackage(".", "./"+path.String(), action)
			continue
		}

		pathDir := path[:len(path)-1]

		if pathDir.IsDir() {
			packages.addModifiedPackage(pathDir.String(), path.String(), action)
		}
	}
}

func setImportModifiedPackages(packages Packages) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	_ = walk.Walk(dir, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		packageName := getPackageName(filepath)
		if _, ok := packages[packageName]; ok {
			return nil
		}

		imports := parsing.GetImports(filepath)
		for _, i := range imports {
			if !strings.HasPrefix(i, vars.ModuleName) {
				continue
			}

			importPackageName := getPackageNameFromImport(i)
			if _, ok := packages[importPackageName]; ok {
				packages.addImportModifiedPackage(importPackageName, packageName)
			}
		}

		return nil
	})
}

func getPackageName(p string) string {
	dir := path.Dir(p)
	splitDir := strings.SplitN(dir, vars.ModuleName, 2)
	if len(splitDir) == 2 {
		dir = strings.Trim(splitDir[1], "/")
		if dir == "" {
			dir = "."
		}
		return dir
	} else {
		return ""
	}
}

func getPackageNameFromImport(i string) string {
	split := strings.SplitN(i, vars.ModuleName, 2)
	if len(split) == 2 {
		pkg := strings.Trim(split[1], "/")
		return pkg
	} else {
		return ""
	}
}
