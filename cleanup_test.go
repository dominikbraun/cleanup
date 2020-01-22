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

import "testing"

// Test_readBranchNames tests if readBranchNames returns the correct gone
// branch names for a given `git branch -vv` output.
func Test_readBranchNames(t *testing.T) {
	gitOutput := `
* master    34a234a [origin/master] Merged some features
  feature/1 34a234a [origin/feature/1: gone] Implemented endpoints
  feature/2 3fc2e37 [origin/feature/2: behind 71] Added CLI command`

	branches := readBranchNames([]byte(gitOutput), ": gone]")

	if len(branches) != 1 {
		t.Errorf("Expected %v branches, got %v", 1, len(branches))
	}

	if branches[0] != "feature/1" {
		t.Errorf("Expected branch %s, got %s", "feature/1", branches[0])
	}
}
