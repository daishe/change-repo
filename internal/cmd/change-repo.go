// Copyright 2022-2026 Marek Dalewski

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
	"fmt"

	"github.com/spf13/cobra"

	"github.com/daishe/change-repo/internal/status"
)

func NewChangeRepoCmd(info *AppInfo) *cobra.Command {
	c := &cobra.Command{
		Use:   "change-repo [flags] [...directory]",
		Short: "A simple CLI utility to change the current working directory to\none containing Git repository.",
		Long:  "change-repo - a simple CLI utility to change the current working directory to\none containing Git repository.",
	}

	maxdepthFlag := c.Flags().Uint("maxdepth", 20, "controls recursion depth when scanning for Git repositories")
	showFlag := c.Flags().Bool("show", false, "controls wether to only show the path to the selected for Git repositories or navigation to it in a sub-shell")
	ignoreSoftErrorsFlag := c.Flags().Bool("ignore-soft-errors", false, "soft errors will not cause non zero exit code")
	versionFlag := c.Flags().Bool("version", false, "display version and copyright information")

	c.Run = func(c *cobra.Command, args []string) {
		if *versionFlag {
			displayVersion(c, info)
			status.Exit(c.Context(), 0)
		}

		baseDir := pickBaseDir(c, baseDirs(args))
		repos, nonRepos := []string(nil), []string(nil)
		scanForRepos(c, baseDir, *maxdepthFlag, &repos, &nonRepos)
		repoDir := pickRepoDir(c, baseDir, repos, nonRepos)

		if *ignoreSoftErrorsFlag {
			status.ClearErr(c.Context())
		}

		if *showFlag {
			fmt.Fprintf(c.OutOrStdout(), "%s\n", repoDir)
			return
		}
		changeDir(c, repoDir)
	}

	return c
}
