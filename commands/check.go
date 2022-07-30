package commands

import (
	"github.com/alexisvisco/gwd/pkg/diff/modules"
	"github.com/alexisvisco/gwd/pkg/output"
	"github.com/alexisvisco/gwd/pkg/vars"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func runCheck(_ *cobra.Command, args []string) error {
	if err := vars.LoadFilesChanged(); err != nil {
		return err
	}

	modulesChanged, err := modules.FromFilesChanged(vars.FilesChanged)
	if err != nil {
		return err
	}

	for _, mod := range modulesChanged.Modules {
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
	rootCmd.AddCommand(checkCommand)
}
