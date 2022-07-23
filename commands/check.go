package commands

import (
	"github.com/alexisvisco/gwd/pkg/diff"
	"github.com/alexisvisco/gwd/pkg/output"
	"github.com/alexisvisco/gwd/pkg/vars"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func runCheck(_ *cobra.Command, args []string) error {
	modules, err := diff.Diff(vars.Repository, previousReference, currentReference)
	if err != nil {
		return err
	}

	for _, mod := range modules.Modules {
		if mod.ModulePath == args[0] || mod.ModuleName == args[0] {
			output.Print(output.String(mod.ModulePath))
			return nil
		}
	}

	return errors.New("module not found")
}

var checkCommand = &cobra.Command{
	Use:     "check <module or path name>",
	Short:   "Return success exit code if the module is modified, arg 1 must be one of the module name or path",
	RunE:    runCheck,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"c"},
}

func init() {
	checkCommand.Flags().StringVarP(
		&previousReference,
		"previous-ref",
		"p",
		"master",
		"set the previous reference to diff with current one.\nIt can be a tag, branch or commit hash",
	)

	checkCommand.Flags().StringVarP(
		&currentReference,
		"current-ref",
		"c",
		"",
		"set the current reference to diff with previous one.\nIt can be a tag, branch or commit hash",
	)

	rootCmd.AddCommand(checkCommand)
}
