// Copyright 2019 The cleanup authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main provides the cleanup executable and its implementation.
package main

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

const version string = "UNSPECIFIED"

// main builds the CLI commands and executes the desired sub-command.
func main() {
	var branchesOptions BranchesOptions
	var versionOptions VersionOptions

	cleanup := &cobra.Command{
		Use:     "cleanup",
		Aliases: []string{"git-cleanup"},
		Short:   `💫 Remove gone Git branches with ease.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	branches := &cobra.Command{
		Use:   "branches <PATH>",
		Short: `Delete local branches that are gone on the remote`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Branches(args[0], &branchesOptions, os.Stdout)
		},
	}

	branches.Flags().BoolVarP(&branchesOptions.HasMultipleRepos, "has-multiple-repos",
		"m", false, `Delete branches in sub-repositories`)
	branches.Flags().BoolVarP(&branchesOptions.Force, "force",
		"f", false, `Force the deletion, ignoring warnings`)
	branches.Flags().BoolVarP(&branchesOptions.DryRun, "dry-run",
		"d", false, `Preview the branches without deleting them`)
	branches.Flags().StringVarP(&branchesOptions.Exclude, "exclude",
		"e", "", `Exclude one or more branches from deletion`)

	version := &cobra.Command{
		Use:   "version",
		Short: `Display version information`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Version(&versionOptions, os.Stdout)
		},
	}

	version.Flags().BoolVarP(&versionOptions.Quiet, "quiet",
		"q", false, `Only print the version number`)

	cleanup.AddCommand(branches)
	cleanup.AddCommand(version)

	if err := cleanup.Execute(); err != nil {
		log.Fatal(err)
	}
}
