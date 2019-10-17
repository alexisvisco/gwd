package packages

import (
	"os"
	"path"
	"strings"

	"github.com/MichaelTJones/walk"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"

	"github.com/alexisvisco/gta/pkg/gta/parsing"
	"github.com/alexisvisco/gta/pkg/gta/vars"
)

func FromChanges(changes merkletrie.Changes, wt *git.Worktree) Packages {
	packages := make(Packages)
	containModification := PresenceReason{Reason: ReasonModification}

	for _, ch := range changes {
		path := ch.To
		if path.String() == "" {
			path = ch.From
		}

		if path.IsDir() {
			packages[path.String()] = containModification
			continue
		}

		if len(path) == 1 {
			packages["."] = containModification
			continue
		}

		pathDir := path[:len(path)-1]

		if pathDir.IsDir() {
			packages[pathDir.String()] = containModification
		}
	}

	dir, err := os.Getwd()
	if err != nil {
		return packages
	}

	packagesWhichContainModifiedPackage := make(map[string]PresenceReason)

	err = walk.Walk(dir, func(filepath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		packageName := getPackageName(filepath)
		if _, ok := packages[packageName] ; ok {
			return nil
		}
		if _, ok := packagesWhichContainModifiedPackage[packageName] ; ok {
			return nil
		}


		imports := parsing.GetImports(filepath)

		for _, i := range imports {
			if !strings.HasPrefix(i, vars.ModuleName) {
				continue
			}

			importPackageName := getPackageName(i)

			packagesWhichContainModifiedPackage[packageName] = PresenceReason{
				Reason:          ReasonContainModifiedPackage,
				ModifiedPackage: &importPackageName,
			}
		}

		for k, v := range packagesWhichContainModifiedPackage {
			packages[k] = v
		}

		return nil
	})


	return packages
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
