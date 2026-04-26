package find

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	NoDescend = errors.New("no descend") //nolint:errname,staticcheck // The name do not start with 'Err' intentionally.
)

const DefaultMaxDepth uint = 20

type Find struct {
	Root     string
	Maxdepth uint
	OnError  func(error) bool
	Cond     func(string, fs.FileMode) (bool, error)
}

func (s Find) onError(err error) bool {
	if err == nil {
		return true
	}
	if s.OnError == nil {
		return true
	}
	if s.OnError(err) {
		return true
	}
	return false
}

func (s Find) cond(path string, mode fs.FileMode) (bool, error) {
	if s.Cond == nil {
		return true, nil
	}
	return s.Cond(path, mode)
}

func (s Find) AllDirs() ([]string, []string) {
	match, nonMatch := []string{}, []string{}
	s.allDirs(&match, &nonMatch)
	return match, nonMatch
}

func (s Find) allDirs(match, nonMatch *[]string) bool {
	if s.Maxdepth == 0 {
		return true
	}

	entries, err := os.ReadDir(s.Root)
	if err != nil {
		return s.onError(err)
	}

	matchFound, nonMatchCache := false, []string(nil)

	for _, e := range entries {
		path := filepath.Join(s.Root, e.Name())
		isDir, err := isDir(path, e.Type())
		if err != nil {
			if !s.onError(err) {
				return false
			}
			continue
		}
		if !isDir {
			continue
		}

		ok, err := s.cond(path, e.Type())
		if err != nil && err != NoDescend { //nolint:errorlint // Direct comparison with NoDescend is required, as it should not be wrapped.
			if !s.onError(err) {
				return false
			}
			continue
		}
		if ok {
			matchFound = true
			*match = append(*match, path)
		} else {
			nonMatchCache = append(nonMatchCache, path)
		}

		if err == NoDescend { //nolint:errorlint // Direct comparison with NoDescend is required, as it should not be wrapped.
			continue
		}
		next := s
		next.Root = path
		next.Maxdepth -= 1
		if !next.allDirs(match, nonMatch) {
			return false
		}
	}

	if matchFound {
		*nonMatch = append(*nonMatch, nonMatchCache...)
	}
	return true
}

func isDir(path string, mode fs.FileMode) (bool, error) {
	if mode.IsDir() {
		return true, nil
	}

	if mode.Type()&os.ModeSymlink == 0 {
		return false, nil
	}
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false, err
	}
	stat, err := os.Stat(target)
	if err != nil {
		return false, err
	}
	return stat.Mode().IsDir(), nil
}
