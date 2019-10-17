package cmd

import (
	"github.com/spf13/cobra"

	"github.com/alexisvisco/gta/pkg/gta"
	"github.com/alexisvisco/gta/pkg/gta/diff"
	"github.com/alexisvisco/gta/pkg/gta/vars"
)

func changes(_ *cobra.Command, _ []string) error {
	packages, err := diff.Diff(vars.Repository, previousReference, currentReference)
	if err != nil {
		return err
	}

	gta.Output(packages)

	return nil
}

// changesCommand represents the changes command
var changesCommand = &cobra.Command{
	Use:   "changes",
	Short: "Show packages changes between a previous ref and the current ref ",
	Long: "This command accept 2 flags.\n" +
		"If --current-ref is omitted, gta will use your current uncommitted changes.\n" +
		"The command do nothing except printing the packages which should be tested.",
	RunE:    changes,
	Aliases: []string{"diff"},
}

func init() {
	changesCommand.LocalFlags().StringVarP(
		&previousReference,
		"previous-ref",
		"p",
		"master",
		"set the previous reference to diff with current one.\nIt can be a tag, branch or commit hash",
	)

	changesCommand.LocalFlags().StringVarP(
		&currentReference,
		"current-ref",
		"c",
		"",
		"set the current reference to diff with previous one.\nIt can be a tag, branch or commit hash",
	)

	rootCmd.AddCommand(changesCommand)
}
