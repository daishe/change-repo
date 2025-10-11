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

package ui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
)

func SelectPrompt(msg string, list []string) (int, error) {
	if runtime.GOOS == "windows" {
		return selectPrompt_windows(msg, list)
	}

	component := promptui.Select{
		Label:             msg,
		Items:             list,
		Size:              20,
		HideHelp:          true,
		StartInSearchMode: true,
		Searcher: func(input string, i int) bool {
			return fuzzy.Match(input, list[i])
		},
		Keys: &promptui.SelectKeys{
			Prev:     promptui.Key{Code: readline.CharPrev, Display: "↑"},
			Next:     promptui.Key{Code: readline.CharNext, Display: "↓"},
			PageUp:   promptui.Key{Code: readline.CharBackward, Display: "←"},
			PageDown: promptui.Key{Code: readline.CharForward, Display: "→"},
			Search:   promptui.Key{Code: readline.CharCtrlW, Display: "^W"},
		},
	}

	idx, _, err := component.Run()
	if err != nil {
		return 0, err
	}
	return idx, nil
}

func selectPrompt_windows(msg string, list []string) (int, error) {
	for i, li := range list {
		fmt.Printf("  %d: %s\n", i+1, li)
	}
	fmt.Printf("%s: ", msg)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	fmt.Print("\n")
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = errors.New("invalid selection")
		}
		return -1, err
	}
	idx, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil || idx < 1 || idx > len(list) {
		return -1, errors.New("invalid selection")
	}
	return idx - 1, nil
}
