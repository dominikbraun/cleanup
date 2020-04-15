package main

import (
	"strings"
	"testing"
)

// Test_readBranchNames tests if readBranchNames returns the correct gone
// branch names for a given `git branch -vv` output.
func Test_readBranchNames(t *testing.T) {
	gitOutput := `
* master    34a234a [origin/master] Merged some features
  feature/1 34a234a [origin/feature/1: gone] Implemented endpoints
  feature/2 3fc2e37 [origin/feature/2: behind 71] Added CLI command`

	branches := readBranchNames([]byte(gitOutput), func(line string) bool {
		return strings.Contains(line, searchExpr)
	})

	if len(branches) != 1 {
		t.Errorf("Expected %v branches, got %v", 1, len(branches))
	}

	if branches[0] != "feature/1" {
		t.Errorf("Expected branch %s, got %s", "feature/1", branches[0])
	}
}

// Test_isExcluded tests if isExcluded checks excluded branches correctly.
func Test_isExcluded(t *testing.T) {
	type assertion struct {
		branch   string
		exclude  []string
		expected bool
	}

	assertions := []assertion{
		{
			branch:   "feature/1",
			exclude:  []string{" feature/0 ", "feature/2", "my-fix"},
			expected: false,
		},
		{
			branch:   "ci-setup",
			exclude:  []string{"feature/1", "another-fix", "ci-setup"},
			expected: true,
		},
		{
			branch:   "feature/2",
			exclude:  []string{"feature/1", " feature/2", "feature/3 "},
			expected: true,
		},
	}

	for _, a := range assertions {
		if result := isExcluded(a.branch, a.exclude); result != a.expected {
			t.Errorf("%s: expected %v, got %v", a.branch, a.expected, result)
		}
	}
}
