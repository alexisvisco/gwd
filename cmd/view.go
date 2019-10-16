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
	"errors"
	"fmt"

	"github.com/alexisvisco/gta/pkg/gta"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie/noder"
)

func view(_ *cobra.Command, args []string) error {
	oldRef := args[0]
	newRef := ""
	if len(args) == 2 {
		newRef = args[1]
	}

	oldNoder, newNoder, e := getNoders(oldRef, newRef)
	if e != nil {
		return e
	}

	changes, err := merkletrie.DiffTree(oldNoder, newNoder, gta.DiffTreeIsEquals)
	if err != nil {
		return errors.New("unable to compare changes between refs")
	}

	packages := make(map[string]bool)
	for _, ch := range changes {
		path := ch.To
		if path.String() == "" {
			path = ch.From
		}

		if path.IsDir() {
			packages[path.String()] = true
			continue
		}

		if len(path) == 1 {
			packages["."] = true
			continue
		}

		pathDir := path[len(path)-2]
		if pathDir.IsDir() {
			packages[pathDir.String()] = true
		}
	}

	for directories := range packages {
		fmt.Println(directories)
	}
	return nil
}

func getNoders(oldRef string, newRef string) (newNoder noder.Noder, oldNoder noder.Noder, err error) {
	oldNoder, err = gta.GetTree(repository, oldRef)
	if err != nil {
		return nil, nil, err
	}
	newNoder, err = gta.GetTree(repository, newRef)
	if err != nil {
		return nil, nil, err
	}
	return oldNoder, newNoder, nil
}

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view <old ref> [<new ref>]",
	Short: "View packages that have changed from a revision",
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
