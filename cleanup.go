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
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

const (
	gitDir     string = ".git"
	searchExpr string = "]"
)

// RepositoryPath describes the filesystem path for a repository.
type RepositoryPath string

// BranchesOptions are user-defined options for the `branch` command.
type BranchesOptions struct {
	HasMultipleRepos bool
	Force            bool
}

// Branches is the entry point for the `branch` command and deletes all
// gone Git branches under a specified path. This path either has to be
// a repository or a root directory that contains multiple repositories.
//
// Returns an error if no Git repository can be found in the path.
func Branches(path string, options *BranchesOptions) error {
	repositories, err := repositoryPaths(path, options.HasMultipleRepos)
	if err != nil {
		return err
	}

	if len(repositories) == 0 {
		return errors.New("no Git repository found")
	}

	for _, repo := range repositories {
		deleted, err := deleteBranches(repo)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("Deleted from repository at `%s`:\n", repo)

		for _, branch := range deleted {
			fmt.Printf("- %s\n", branch)
		}
	}

	return nil
}

// deleteBranches deletes branches in a repository that are considered
// gone. For determining these branches, `git branch -vv` will be used.
//
// Returns a list of all branches that have been deleted successfully.
func deleteBranches(path RepositoryPath) ([]string, error) {
	cmd := exec.Command("git", "branch", "-vv")
	cmd.Dir = string(path)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	deleted := make([]string, 0)

	for _, branch := range readBranchNames(out, searchExpr) {
		cmd := exec.Command("git", "branch", "-d", branch)
		cmd.Dir = string(path)

		_, _ = cmd.Output()
		deleted = append(deleted, branch)
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
//  feature/12	3fc2e37 [origin/feature/2: behind 71] Added CLI command
//
// Using a filter for gone branches, merely feature/1 will be returned.
//
// Returns a list of branch names that appeared in the byte sequence.
func readBranchNames(buf []byte, filter string) []string {
	branches := make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewBuffer(buf))

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, filter) {
			for i, c := range line {
				if i > 1 && string(c) == " " {
					branches = append(branches, line[2:i])
					break
				}
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
// directory containing multiple repositories, all paths will be returned.
//
// hasMultipleRepos indicates if the provided path is a parent directory.
func repositoryPaths(path string, hasMultipleRepos bool) ([]RepositoryPath, error) {
	paths := make([]RepositoryPath, 0)

	if !hasMultipleRepos {
		isRepo, err := isRepository(RepositoryPath(path))
		if err != nil {
			return nil, err
		}

		if isRepo {
			paths = append(paths, RepositoryPath(path))
			return paths, nil
		}
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
