/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/alexisvisco/gta/pkg/gta"
	"github.com/alexisvisco/gta/pkg/gta/diff"
	"github.com/spf13/cobra"
)

const localRef = ""

func view(_ *cobra.Command, args []string) error {
	previousRef := args[0]
	currentRef := localRef
	if len(args) == 2 {
		currentRef = args[1]
	}

	previousNoder, err := gta.GetTree(repository, previousRef)
	if err != nil {
		return err
	}

	if currentRef == localRef {
		packages, err := diff.LocalDiff(repository, previousNoder)
		if err != nil {
			return err
		}

		fmt.Println(packages.String())
	} else {
		currentNoder, err := gta.GetTree(repository, currentRef)
		if err != nil {
			return err
		}

		packages, err := diff.Diff(repository, previousNoder, currentNoder)
		if err != nil {
			return err
		}

		fmt.Println(packages.String())
	}

	return nil
}

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view <previous ref> [<current re>]",
	Short: "View packages that have changed between a previous ref and the current ref ",
	Long: "This command accept two arguments.\n" +
		"Each arguments can be either a tag, a branch or a specific commit hash.\n" +
		"If 'new' argument is not specified, gta will use your current uncommitted changes.\n" +
		"The command do nothing except printing the packages which should be tested.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("view called")
	},
	Args: cobra.RangeArgs(1, 2),
	RunE: view,
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
