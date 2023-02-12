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

package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/daishe/change-repo/internal/status"
	"github.com/daishe/change-repo/internal/ui"
)

const (
	changeRepoPathEnv  = "CHANGE_REPO_PATH"
	changeRepoShellEnv = "CHANGE_REPO_SHELL"
	shellEnv           = "SHELL"
)

type AppInfo struct {
	AppVersion string
	CommitHash string
}

func displayVersion(c *cobra.Command, info *AppInfo) {
	fmt.Fprintf(c.OutOrStdout(), "Version of application: %s, commit: %s\n", info.AppVersion, info.CommitHash)
	fmt.Fprintf(c.OutOrStdout(), "\n")
	fmt.Fprintf(c.OutOrStdout(), "Copyright 2022 Marek Dalewski. License GPLv3+: GNU General Public License version 3 or later\n")
	fmt.Fprintf(c.OutOrStdout(), "\n")
	fmt.Fprintf(c.OutOrStdout(), "You should have received a copy of the GNU General Public License along with this program. If not, see <https://gnu.org/licenses/gpl.html>.\n")
	fmt.Fprintf(c.OutOrStdout(), "This is free software: you are free to change and redistribute it. This program comes with ABSOLUTELY NO WARRANTY.\n")
}

func baseDirs(args []string) (res []string) {
	if len(args) > 0 {
		return args
	}

	fromEnv := func(envName string) []string {
		parts := strings.Split(os.Getenv(envName), string(os.PathListSeparator))
		res := make([]string, 0, len(parts))
		for _, p := range parts {
			if p = filepath.Clean(p); p != "" {
				res = append(res, p)
			}
		}
		return res
	}

	res = append(res, fromEnv(changeRepoPathEnv)...)
	if len(res) == 0 {
		res = append(res, ".") // assume current working directory
	}
	return res
}

func pickBaseDir(c *cobra.Command, baseDirs []string) string {
	baseDirs = stringsSortedUnique(baseDirs)
	if len(baseDirs) == 0 {
		status.CheckErr(c.Context(), "no directories to search for Git repositories within them")
	}

	baseDir := baseDirs[0]
	if len(baseDirs) > 1 {
		baseIdx, err := ui.SelectPrompt("Select set of repositories", baseDirs)
		status.CheckErr(c.Context(), err)
		baseDir = baseDirs[baseIdx]
		fmt.Fprintf(c.OutOrStdout(), "> %s\n", baseDir)
	}

	baseDir, err := filepath.Abs(baseDir)
	status.CheckErr(c.Context(), err)
	return baseDir
}

func scanForRepos(c *cobra.Command, root string, maxdepth uint, repos *[]string, nonRepos *[]string) {
	if maxdepth == 0 {
		return
	}

	someRepoFound, nonReposCache := false, []string(nil)

	entries, err := os.ReadDir(root)
	status.ShowErr(c.Context(), err)

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		innerDir := filepath.Join(root, e.Name())
		isGitRepo, err := isGitRepo(innerDir)
		status.ShowErr(c.Context(), err)

		if isGitRepo {
			someRepoFound = true
			*repos = append(*repos, innerDir)
		} else {
			nonReposCache = append(nonReposCache, innerDir)
			scanForRepos(c, innerDir, maxdepth-1, repos, nonRepos)
		}
	}

	if someRepoFound {
		*nonRepos = append(*nonRepos, nonReposCache...)
	}
}

var gitDirectory string = ".git"

func isGitRepo(path string) (bool, error) {
	s, err := os.Stat(filepath.Join(path, gitDirectory))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("could not determine if path contains a Git repository: %w", err)
	}
	return s.IsDir(), nil
}

func pickRepoDir(c *cobra.Command, root string, repos []string, nonRepos []string) string {
	repos = stringsSortedUnique(pathsFromRoot(repos, root))
	nonRepos = stringsSortedUnique(pathsFromRoot(nonRepos, root))

	if len(nonRepos) > 0 {
		repos = append(repos, "<locations, that do not contain Git repository>")
	}

	pickNonRepo := func() string {
		nonRepoIdx, err := ui.SelectPrompt("Select location", nonRepos)
		status.CheckErr(c.Context(), err)
		fmt.Fprintf(c.OutOrStdout(), "> %s\n", nonRepos[nonRepoIdx])
		return filepath.Join(root, nonRepos[nonRepoIdx])
	}

	pickRepo := func() string {
		repoIdx, err := ui.SelectPrompt("Select repository", repos)
		status.CheckErr(c.Context(), err)
		fmt.Fprintf(c.OutOrStdout(), "> %s\n", repos[repoIdx])

		if len(nonRepos) > 0 && repoIdx == len(repos)-1 { // non repo request
			return pickNonRepo()
		}

		return filepath.Join(root, repos[repoIdx])
	}

	if len(repos) == 0 {
		if len(nonRepos) == 0 {
			status.CheckErr(c.Context(), "no Git repositories found")
		}
		fmt.Fprintf(c.OutOrStdout(), "No Git repositories found\n")
		return pickNonRepo()
	}
	return pickRepo()
}

func changeDir(c *cobra.Command, dir string) {
	run(c, dir, shellExecutable())
}

func stringsSortedUnique(x []string) []string {
	m := map[string]struct{}{}
	for _, v := range x {
		m[v] = struct{}{}
	}

	x = make([]string, 0, len(m))
	for v := range m {
		x = append(x, v)
	}

	sort.Strings(x)
	return x
}

func pathsFromRoot(x []string, root string) []string {
	x = append([]string(nil), x...)
	for i, v := range x {
		v = strings.TrimPrefix(v, root)
		v = strings.TrimPrefix(v, string(os.PathSeparator))
		v = strings.TrimPrefix(v, "/") // forced Unix path separator
		x[i] = filepath.Clean(v)
	}
	return x
}

func shellExecutable() string {
	if v, ok := os.LookupEnv(changeRepoShellEnv); ok {
		return v
	}
	if v, ok := os.LookupEnv(shellEnv); ok {
		return v
	}
	return "sh"
}

func run(c *cobra.Command, dir string, name string, args ...string) {
	cmd := exec.CommandContext(c.Context(), name, args...)
	cmd.Dir = dir
	cmd.Stderr = c.ErrOrStderr()
	cmd.Stdin = c.InOrStdin()
	cmd.Stdout = c.OutOrStdout()
	if err := cmd.Run(); err != nil {
		if ee := (*exec.ExitError)(nil); errors.As(err, &ee) {
			status.Exit(c.Context(), ee.ExitCode())
		}
		status.CheckErr(c.Context(), fmt.Errorf("executing %q: %w", name, err))
	}
}
