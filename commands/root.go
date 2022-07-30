package commands

import (
	"fmt"
	"github.com/alexisvisco/gwd/pkg/diff/modules"
	"github.com/alexisvisco/gwd/pkg/diff/packages"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexisvisco/gwd/pkg/output"
	"github.com/alexisvisco/gwd/pkg/parsing"
	"github.com/alexisvisco/gwd/pkg/utils"
	"github.com/spf13/cobra"

	"github.com/alexisvisco/gwd/pkg/vars"
)

func preRun(_ *cobra.Command, _ []string) error {
	var err error

	if !utils.PathExist(vars.GoWorkFileName) {
		return fmt.Errorf("go workspace file %q not found", vars.GoWorkFileName)
	}

	vars.GoWorkModulePaths, err = parsing.GetWorkspaceModulePaths(vars.GoWorkFileName)
	if err != nil {
		return err
	}
	if len(vars.GoWorkModulePaths) == 0 {
		return fmt.Errorf("no modules found in go workspace file %q", vars.GoWorkFileName)
	}

	for _, modulePath := range vars.GoWorkModulePaths {
		goModPath := filepath.Join(modulePath, "go.mod")

		if !utils.PathExist(goModPath) {
			return fmt.Errorf("go module %q not found", goModPath)
		}

		moduleName, err := parsing.GetModuleName(goModPath)
		if err != nil {
			return err
		}

		modulePath = strings.Trim(modulePath, "./")

		vars.ModulePathToModuleName[modulePath] = moduleName
		vars.ModuleNameToModulePath[moduleName] = modulePath
	}

	return nil
}

func run(_ *cobra.Command, args []string) error {
	if err := vars.LoadFilesChanged(); err != nil {
		return err
	}

	modulesChanged, err := modules.FromFilesChanged(vars.FilesChanged)
	if err != nil {
		return err
	}

	if len(args) == 1 {
		for _, mod := range modulesChanged.Modules {
			if args[0] == mod.ModulePath || args[0] == mod.ModuleName {
				output.Print(output.StringArray(lo.Map(lo.Keys(mod.PackagesModified), func(t packages.ImportPath, i int) string {
					return string(t)
				})))
				return nil
			}
		}
		return errors.New("module not found or no diff")
	}

	output.Print(modulesChanged)
	return nil
}

var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "gwd",
	Short:   "List go modules that differ based on a list of files provided in stdin, it also returns modules that have package that import a module that has been modified",
	Long: `The gwd command use the go.work file to know the list of go modules the project has.
The foal of gwd is to list the modules that have been modified or modules that import a module that has been modified.
This command is useful to know which modules have been modified in a project for example in a CI/CD system to lint, test, build only modules that changed.`,
	Example:           `git diff "v0.0.1" --name-only | gwd --stdin`,
	PersistentPreRunE: preRun,
	RunE:              run,
	SilenceUsage:      true,
	SilenceErrors:     true,
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(
		&vars.GoWorkFileName,
		"go-work",
		"w", "go.work",
		"golang workspace file name",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&vars.OutputJson,
		"json",
		"j", false,
		"output of commands will be json format",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&vars.OutputVerbose,
		"verbose",
		"v", false,
		"output of commands will be in an human verbose format",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&vars.FilesChangedFromStdin,
		"stdin",
		"i", false,
		"read stdin, used in conjuncture with git diff --name-only for example",
	)

	rootCmd.PersistentFlags().StringVarP(
		&vars.FilesChangedFromFile,
		"file",
		"f", "",
		"read a file, each line must be a file path",
	)

	if err := rootCmd.Execute(); err != nil {
		output.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
