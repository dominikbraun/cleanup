// Package main provides the cleanup executable and its implementation.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

const (
	gitDir     string = ".git"
	searchExpr string = ": gone]"
)

// RepositoryPath describes the filesystem path for a repository.
type RepositoryPath string

// BranchesOptions are user-defined options for the `branches` command
// which will be instantiated and populated by the CLI.
type BranchesOptions struct {
	HasMultipleRepos bool
	Force            bool
	DryRun           bool
	Exclude          string
	Where            string
	AndWhere         string
}

// VersionOptions are user-defined options for the `version` command
// which will be instantiated and populated by the CLI.
type VersionOptions struct {
	Quiet bool
}

// RunBranches is the entry point for the `branch` command and deletes
// all gone Git branches under a specified path. This path either has
// to be a repository or a root directory that contains multiple repos.
//
// Returns an error if no Git repository can be found in the path.
func RunBranches(path string, options *BranchesOptions, w io.Writer) error {
	repositories, err := repositoryPaths(path, options.HasMultipleRepos)
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		return errors.New("no Git repository found")
	}

	exclude := strings.Split(options.Exclude, ",")
	exclude = append(exclude, "master")

	for _, repo := range repositories {
		deleted, err := deleteBranches(repo, options, exclude)
		if err != nil {
			output := fmt.Sprintf("Error in repository %s: %s\n", repo, err.Error())
			_, _ = w.Write([]byte(output))
			continue
		}

		if len(deleted) == 0 {
			output := fmt.Sprintf("No gone branches found in repository %s.\n", repo)
			_, _ = w.Write([]byte(output))
			continue
		}

		output := fmt.Sprintf("Found gone branches in repository %s:\n", repo)
		_, _ = w.Write([]byte(output))

		for branch, err := range deleted {
			var output string

			switch {
			case options.DryRun:
				output = fmt.Sprintf("\t- Will delete branch %s\n", branch)
			case err != nil:
				output = fmt.Sprintf("\t- Failed to delete branch %s: %s\n", branch, err.Error())
			default:
				output = fmt.Sprintf("\t- Deleted branch %s\n", branch)
			}

			_, _ = w.Write([]byte(output))
		}
	}

	return nil
}

// Version displays version information for cleanup.
func Version(options *VersionOptions, w io.Writer) error {
	var output string

	switch {
	case options.Quiet:
		output = fmt.Sprintf("%s\n", version)
	default:
		output = fmt.Sprintf("cleanup version %s\n", version)
	}

	_, _ = w.Write([]byte(output))

	return nil
}

// deleteBranches deletes branches in a repository that are considered
// gone. For determining these branches, `git branch -vv` will be used.
//
// Returns a map of deleted branch names mapped against an error value.
// If the error value is not nil, the corresponding branch couldn't be
// deleted successfully. The second return value indicates if an error
// occurred when running the `git branch -vv` command.
func deleteBranches(path RepositoryPath, options *BranchesOptions, exclude []string) (map[string]error, error) {
	cmd := exec.Command("git", "branch", "-vv")
	cmd.Dir = string(path)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	branches := readBranchNames(out, func(line string) bool {
		if options.Where != "" {
			return strings.Contains(line, options.Where)
		}
		if options.AndWhere != "" {
			return strings.Contains(line, searchExpr) && strings.Contains(line, options.AndWhere)
		}
		return strings.Contains(line, searchExpr)
	})

	deleted := make(map[string]error)

	for _, branch := range branches {
		if isExcluded(branch, exclude) {
			continue
		}

		if !options.DryRun {
			cmd := exec.Command("git", "branch", "-d", branch)
			cmd.Dir = string(path)

			_, err := cmd.Output()
			deleted[branch] = err
		} else {
			deleted[branch] = nil
		}
	}

	return deleted, nil
}

// readBranchNames reads Git branch names contained in a byte slice,
// which is expected to be the output of `git branch -vv`. Each line
// is tested against a filter function and will only be processed if
// it passes that filter.
//
// The Git output is expected to look like this:
//
// * master		34a234a [origin/master] Merged some features
//  feature/1	34a234a [origin/feature/1: gone] Implemented endpoints
//  feature/2	3fc2e37 [origin/feature/2: behind 71] Added CLI command
//
// Returns a list of branch names that appeared in the byte sequence.
func readBranchNames(buf []byte, filter func(string) bool) []string {
	branches := make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(buf))

	for scanner.Scan() {
		line := scanner.Text()

		if !filter(line) {
			continue
		}

		for i, c := range line {
			if i > 1 && string(c) == " " {
				branches = append(branches, line[2:i])
				break
			}
		}
	}

	return branches
}

// repositoryPaths returns all repository paths contained in a path.
//
// If the provided path itself is a repository, it will be included in
// the returned path slice. If the hasMultipleRepos flag is `true`, all
// direct parent directories that are Git repositories will be returned.
func repositoryPaths(path string, hasMultipleRepos bool) ([]RepositoryPath, error) {
	paths := make([]RepositoryPath, 0)

	isRepo, err := isRepository(RepositoryPath(path))
	if err != nil {
		return nil, err
	}
	if isRepo {
		paths = append(paths, RepositoryPath(path))
	}

	if !hasMultipleRepos {
		return paths, nil
	}

	content, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, repo := range content {
		if !repo.IsDir() {
			continue
		}

		currentPath := RepositoryPath(path + "/" + repo.Name())

		isRepo, err := isRepository(currentPath)
		if err != nil {
			return nil, err
		}
		if isRepo {
			paths = append(paths, currentPath)
		}
	}

	return paths, nil
}

// isRepository checks if a given path is a Git repository, meaning that
// it contains a `.git` directory.
func isRepository(path RepositoryPath) (bool, error) {
	content, err := ioutil.ReadDir(string(path))
	if err != nil {
		return false, err
	}

	for _, item := range content {
		if item.IsDir() && item.Name() == gitDir {
			return true, nil
		}
	}

	return false, nil
}

// isExcluded checks if a branch is contained in an slice of excluded
// branches. Whitespaces will be skipped when comparing the branches.
func isExcluded(branch string, exclude []string) bool {
	for _, e := range exclude {
		if branch == strings.Trim(e, " ") {
			return true
		}
	}

	return false
}
