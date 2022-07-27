package commands

import (
	"github.com/alexisvisco/gwd/pkg/diff"
	"github.com/alexisvisco/gwd/pkg/diff/packages"
	"github.com/alexisvisco/gwd/pkg/output"
	"github.com/alexisvisco/gwd/pkg/vars"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
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

const diffLong = `The diff command returns a list of modules that have been modified between the two references.
It also returns modules that have package that import a module that has been modified.

A reference can be a tag, branch or commit hash.

If the --current-ref flag is not set, the current reference is the current working directory state.
The --previous-ref (or -p) flag must be set (default is set to main branch).

An optional argument can be set which is the module name or path, if it is not set, all modified packages of this module are returned.

	$ gwd diff -> returns all modules that have been modified separated by a space
	$ gwd diff path_or_module_name -> returns all packages that have been modified for this module separated by newline

--verbose or -V option will show more information about the modules, modified packages and their imported package.
`

// diffCommand represents the runDiff command
var diffCommand = &cobra.Command{
	Use:     "diff [<module name or path>]",
	Short:   "Show modules that have been modified between a previous ref and the current ref",
	Long:    diffLong,
	RunE:    runDiff,
	Aliases: []string{"diff"},
}

func init() {
	diffCommand.Flags().StringVarP(
		&previousReference,
		"previous-ref",
		"p",
		"main",
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
