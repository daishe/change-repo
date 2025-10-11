// Copyright 2022 Marek Dalewski

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cmd //nolint:testpackage // this package needs access to some private global values

import (
	"bytes"
	"context"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
)

func emptyCommand() *cobra.Command {
	c := &cobra.Command{}
	c.SetContext(context.Background())
	c.SetErr(&bytes.Buffer{})
	c.SetIn(&bytes.Buffer{})
	c.SetOut(&bytes.Buffer{})
	return c
}

func p(path string) string {
	return strings.ReplaceAll(path, "/", string(os.PathSeparator))
}

func TestScanForReposProcedure(t *testing.T) { //nolint:paralleltest // this test overrides some global values and should not be run in parallel
	tests := []struct {
		name     string
		path     string
		maxdepth uint
		repos    []string
		nonRepos []string
	}{
		{
			name:     "testdata",
			path:     p("./testdata"),
			maxdepth: 20,
			repos:    []string{p("testdata/loc0/repo0"), p("testdata/loc0/repo1"), p("testdata/loc0/subset/repo0"), p("testdata/loc1/repo0")},
			nonRepos: []string{p("testdata/loc0/nonrepo"), p("testdata/loc0/subset"), p("testdata/loc0/subset/nonrepo")},
		},
	}

	for _, td := range tests { //nolint:paralleltest // this test overrides some global values and should not be run in parallel
		t.Run(td.name, func(t *testing.T) {
			originalGitDirectory := gitDirectory
			gitDirectory = "git-data"
			defer func() { gitDirectory = originalGitDirectory }()

			repos := []string(nil)
			nonRepos := []string(nil)
			scanForRepos(emptyCommand(), td.path, 20, &repos, &nonRepos)
			sort.Strings(repos)
			sort.Strings(nonRepos)

			if diff := cmp.Diff(td.repos, repos); diff != "" {
				t.Errorf("repositories set mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(td.nonRepos, nonRepos); diff != "" {
				t.Errorf("non repositories set mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
