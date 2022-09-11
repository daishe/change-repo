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
	versionFlag := c.Flags().Bool("version", false, "display version and copyright information")

	c.RunE = func(c *cobra.Command, args []string) error {
		if *versionFlag {
			displayVersion(c, info)
			status.Exit(c.Context(), 0)
		}

		baseDir := pickBaseDir(c, baseDirs(args))
		respos, nonRespos := []string(nil), []string(nil)
		scanForRepos(c, baseDir, *maxdepthFlag, &respos, &nonRespos)
		repoDir := pickRepoDir(c, baseDir, respos, nonRespos)
		changeDir(c, repoDir)

		return nil
	}

	return c
}
