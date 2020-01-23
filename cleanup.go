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

const (
	fmtRepositoryErr       string = "Error at `%s`: %s\n"
	fmtNoBranchesFound     string = "No gone branches found at `%s`.\n"
	fmtGoneBranchesHeading string = "Found gone branches at `%s`:\n"
	fmtRemovalSuccess      string = "\t- Deleted %s\n"
	fmtRemovalPreview      string = "\t- Will delete %s\n"
	fmtRemovalFailure      string = "\t- Failed to delete %s: %s\n"
)

// RepositoryPath describes the filesystem path for a repository.
type RepositoryPath string

// BranchesOptions are user-defined options for the `branch` command.
type BranchesOptions struct {
	HasMultipleRepos bool
	Force            bool
	DryRun           bool
	Exclude          string
}

// Branches is the entry point for the `branch` command and deletes all
// gone Git branches under a specified path. This path either has to be
// a repository or a root directory that contains multiple repositories.
//
// Returns an error if no Git repository can be found in the path.
func Branches(path string, options *BranchesOptions, w io.Writer) error {
	repositories, err := repositoryPaths(path, options.HasMultipleRepos)

	if err != nil {
		return err
	} else if len(repositories) == 0 {
		return errors.New("no Git repository found")
	}

	exclude := strings.Split(options.Exclude, ",")

	for _, repo := range repositories {
		deleted, err := deleteBranches(repo, options.DryRun, exclude)
		if err != nil {
			output := fmt.Sprintf(fmtRepositoryErr, repo, err.Error())
			_, _ = w.Write([]byte(output))
			continue
		}

		if len(deleted) == 0 {
			output := fmt.Sprintf(fmtNoBranchesFound, repo)
			_, _ = w.Write([]byte(output))
			continue
		}

		output := fmt.Sprintf(fmtGoneBranchesHeading, repo)
		_, _ = w.Write([]byte(output))

		for branch, err := range deleted {
			var output string

			switch {
			case options.DryRun:
				output = fmt.Sprintf(fmtRemovalPreview, branch)
			case err != nil:
				output = fmt.Sprintf(fmtRemovalFailure, branch, err.Error())
			default:
				output = fmt.Sprintf(fmtRemovalSuccess, branch)
			}

			_, _ = w.Write([]byte(output))
		}
	}

	return nil
}

// deleteBranches deletes branches in a repository that are considered
// gone. For determining these branches, `git branch -vv` will be used.
//
// Returns a map with the deleted branches as map keys and an error as
// value. If the error value is not nil for a key, the branch probably
// couldn't be deleted successfully. The second return value indicates
// if an error occurred when executing the `git branch -vv` command.
func deleteBranches(path RepositoryPath, dryRun bool, exclude []string) (map[string]error, error) {
	cmd := exec.Command("git", "branch", "-vv")
	cmd.Dir = string(path)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	deleted := make(map[string]error)

	for _, branch := range readBranchNames(out, searchExpr) {
		if isExcluded(branch, exclude) {
			continue
		}

		if !dryRun {
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

// readBranchNames reads Git branch names contained in a byte slice. It
// expects the output of `git branch -vv` as input. The accepted filter
// provides a simple way for only processing lines that contain a term.
//
// The Git output is expected to look like this:
// * master		34a234a [origin/master] Merged some features
//  feature/1	34a234a [origin/feature/1: gone] Implemented endpoints
//  feature/2	3fc2e37 [origin/feature/2: behind 71] Added CLI command
//
// Using a filter for gone branches, merely feature/1 will be returned.
//
// Returns a list of branch names that appeared in the byte sequence.
func readBranchNames(buf []byte, filter string) []string {
	branches := make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(buf))

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.Contains(line, filter) {
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

// repositoryPaths returns all paths to repositories contained in a given
// path.
//
// In the case that the provided path itself is a repository, it simply
// will return that path. However, in the case that the path is a parent
// directory containing multiple repositories, it returns all that paths.
//
// hasMultipleRepos indicates if the provided path is a parent directory.
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
