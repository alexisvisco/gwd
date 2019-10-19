package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alexisvisco/gta/pkg/gta/diff"
	"github.com/alexisvisco/gta/pkg/gta/vars"
)

func run(_ *cobra.Command, _ []string) error {
	packages, err := diff.Diff(vars.Repository, previousReference, currentReference)
	if err != nil {
		return err
	}

	packagesList := make([]string, 0, len(packages))

	checkAndAppend := func(dir string) {
		if _, err := os.Stat(dir); !os.IsNotExist(err) {
			packagesList = append(packagesList, dir)
		}
	}
	for packageName, detail := range packages {
		checkAndAppend("./" + packageName)
		for packageNameImportedBy := range detail.ImportedBy {
			checkAndAppend("./" + packageNameImportedBy)
		}
	}

	cmd := exec.Command(
		"sh",
		"-c",
		strings.ReplaceAll(command, "${packages}", strings.Join(packagesList, " ")),
	)

	bytes, err := cmd.Output()
	fmt.Print(string(bytes))
	return err
}

// changesCommand represents the changes command
var testCommand = &cobra.Command{
	Use:   "test",
	Short: "Run go test with packages changes between a previous ref and the current ref ",
	Long: "This command accept 3 flags.\n" +
		"If --current-ref is omitted, gta will use your current uncommitted changes.\n" +
		"Inspects the git history to determine which files changed between \n" +
		"the previous red and a feature branch, and uses this information to determine\n" +
		"which packages must be tested for a given build (including packages that import\n" +
		"the changed package).\n" +
		"\nTest each affected packages with the 'go test' command.",
	RunE:    run,
	Aliases: []string{"run"},
}

func init() {
	testCommand.LocalFlags().StringVarP(
		&previousReference,
		"previous-ref",
		"p",
		"master",
		"set the previous reference to diff with current one.\nIt can be a tag, branch or commit hash",
	)

	testCommand.LocalFlags().StringVarP(
		&currentReference,
		"current-ref",
		"c",
		"",
		"set the current reference to diff with previous one.\nIt can be a tag, branch or commit hash",
	)

	testCommand.PersistentFlags().StringVar(
		&command,
		"command",
		"go test ${packages}",
		"execute the command with replacing the ${packages} placeholder",
	)

	rootCmd.AddCommand(testCommand)
}
