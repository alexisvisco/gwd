package commands

import (
	"fmt"
	"github.com/alexisvisco/gwd/pkg/output"
	"github.com/alexisvisco/gwd/pkg/parsing"
	"github.com/alexisvisco/gwd/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"

	"github.com/alexisvisco/gwd/pkg/vars"
)

var rootCmd = &cobra.Command{
	Version: "0.1",
	Use:     "gwd",
	Short:   "Will check git diff of modules from go.work file and only returns modules and packages which runDiff from revision.",
	Example: `- gta diff master feature/ok - will show the packages affected
- gta run master feature/ok  - will tests the packages affected`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

		repo, err := git.PlainOpen(".")
		if err != nil {
			return err
		}

		vars.Repository = repo
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(
		&vars.GoWorkFileName,
		"go-work-file-name",
		"W", "go.work",
		"golang workspace file name",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&vars.OutputJson,
		"json",
		"J", false,
		"output of commands will be json format",
	)

	rootCmd.PersistentFlags().BoolVarP(
		&vars.OutputVerbose,
		"verbose",
		"V", false,
		"output of commands will be in an human verbose format",
	)

	if err := rootCmd.Execute(); err != nil {
		output.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
