package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"

	"github.com/alexisvisco/gta/pkg/gta/parsing"
	"github.com/alexisvisco/gta/pkg/gta/vars"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.1",
	Use:     "gta",
	Short:   "Use gta to tests only packages which changes from revision.",
	Example: `- gta view master feature/ok - will show the packages affected
- gta run master feature/ok  - will tests the packages affected`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if vars.ModuleName == "" {
			return errors.New("module-name flag should be set because go.mod is not available or not readable")
		}

		repo, err := git.PlainOpen(".")
		if err != nil {
			return err
		}

		vars.Repository = repo
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentFlags().StringVarP(
		&vars.ModuleName,
		"module-name",
		"m", parsing.GetModuleName(),
		"module name is used to known the import names for this project",
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
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
