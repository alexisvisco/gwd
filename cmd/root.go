package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

var repository *git.Repository

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.1",
	Use:     "gta",
	Short:   "Use gta to tests only packages which changes from revision.",
	Example: `- gta view master feature/ok - will show the packages affected
- gta run master feature/ok - will tests the packages affected`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		repo, err := git.PlainOpen(".")
		if err != nil {
			return err
		}
		repository = repo
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
