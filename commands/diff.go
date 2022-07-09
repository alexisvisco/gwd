package commands

import (
	"github.com/alexisvisco/gwd/pkg/diff/packages"
	"github.com/alexisvisco/gwd/pkg/output"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/alexisvisco/gwd/pkg/diff"
	"github.com/alexisvisco/gwd/pkg/vars"
)

func runDiff(_ *cobra.Command, args []string) error {
	modules, err := diff.Diff(vars.Repository, previousReference, currentReference)
	if err != nil {
		return err
	}

	if len(args) == 1 {
		for _, mod := range modules.Modules {
			if args[0] == mod.ModulePath || args[0] == mod.ModuleName {
				output.Print(output.StringArray(lo.Map(lo.Keys(mod.PackagesModified), func(t packages.ImportPath, i int) string {
					return string(t)
				})))
				return nil
			}
		}
		return errors.New("module not found or no diff")
	}

	output.Print(modules)
	return nil
}

// diffCommand represents the runDiff command
var diffCommand = &cobra.Command{
	Use:   "diff [<module name or path>]",
	Short: "Show modules and packages diff between a previous ref and the current ref ",
	Long: "This command accept 2 flags.\n" +
		"If --current-ref is omitted, gta will use your current uncommitted runDiff.\n" +
		"The command do nothing except printing the packages which should be tested.\n" +
		"If <module name or path> is specified, it will only print the packages diff for this module.",
	RunE:    runDiff,
	Aliases: []string{"diff"},
}

func init() {
	diffCommand.Flags().StringVarP(
		&previousReference,
		"previous-ref",
		"p",
		"master",
		"set the previous reference to diff with current one.\nIt can be a tag, branch or commit hash",
	)

	diffCommand.Flags().StringVarP(
		&currentReference,
		"current-ref",
		"c",
		"",
		"set the current reference to diff with previous one.\nIt can be a tag, branch or commit hash",
	)

	rootCmd.AddCommand(diffCommand)
}
