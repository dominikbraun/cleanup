// Package main provides the cleanup executable and its implementation.
package main

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

// version will be populated with the current Git tag when building the
// cleanup binary using the CI/CD pipeline.
var version string = "UNSPECIFIED"

// main builds the CLI commands and executes the desired sub-command.
func main() {
	var branchesOptions BranchesOptions
	var versionOptions VersionOptions

	cleanupCmd := &cobra.Command{
		Use:     "cleanup",
		Aliases: []string{"git-cleanupCmd"},
		Short:   `💫 Remove gone Git branches with ease.`,
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	branchesCmd := &cobra.Command{
		Use:   "branches <PATH>",
		Short: `Delete local branches that are gone on the remote`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunBranches(args[0], &branchesOptions, os.Stdout)
		},
	}

	branchesCmd.Flags().BoolVarP(&branchesOptions.HasMultipleRepos, "has-multiple-repos",
		"m", false, `Delete branches in sub-repositories`)
	branchesCmd.Flags().BoolVarP(&branchesOptions.Force, "force",
		"f", false, `Force the deletion, ignoring warnings`)
	branchesCmd.Flags().BoolVarP(&branchesOptions.DryRun, "dry-run",
		"d", false, `Preview the branches without deleting them`)
	branchesCmd.Flags().StringVarP(&branchesOptions.Exclude, "exclude",
		"e", "", `Exclude one or more branches from deletion`)
	branchesCmd.Flags().StringVarP(&branchesOptions.Where, "where",
		"w", "", `Delete all branches whose output contain a given string`)
	branchesCmd.Flags().StringVar(&branchesOptions.AndWhere, "and-where",
		"", `Delete all gone branches whose output contain a given string`)

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: `Display version information`,
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Version(&versionOptions, os.Stdout)
		},
	}

	versionCmd.Flags().BoolVarP(&versionOptions.Quiet, "quiet",
		"q", false, `Only print the version number`)

	cleanupCmd.AddCommand(branchesCmd)
	cleanupCmd.AddCommand(versionCmd)

	if err := cleanupCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
