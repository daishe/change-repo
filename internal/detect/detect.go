package detect

import (
	"fmt"
	"os"
	"path/filepath"
)

type RepoVariant int

func (v RepoVariant) String() string {
	switch v {
	case NonRepo:
		return "NonRepo"
	case NormalRepo:
		return "NormalRepo"
	case BareRepo:
		return "BareRepo"
	case WorkTreeDir:
		return "WorkTreeDir"
	}
	return fmt.Sprintf("NonRepo(%d)", int(v))
}

const (
	NonRepo    RepoVariant = 0
	NormalRepo RepoVariant = 1 << iota
	BareRepo
	WorkTreeDir
)

func IsRepo(path string, expect RepoVariant) (RepoVariant, error) {
	if expect&NormalRepo != 0 {
		is, err := isNormalRepo(path)
		if err != nil {
			return NonRepo, err
		}
		if is {
			return NormalRepo, nil
		}
	}
	if expect&BareRepo != 0 {
		is, err := isBareRepo(path)
		if err != nil {
			return NonRepo, err
		}
		if is {
			return BareRepo, nil
		}
	}
	if expect&WorkTreeDir != 0 {
		is, err := isRepoWorkTree(path)
		if err != nil {
			return NonRepo, err
		}
		if is {
			return WorkTreeDir, nil
		}
	}
	return NonRepo, nil
}

var (
	gitDirectory     = ".git"
	gitFile          = ".git"
	headFile         = "HEAD"
	objectsDirectory = "objects"
)

func isNormalRepo(path string) (bool, error) {
	hasGitDir, err := repoFileCheck(filepath.Join(path, gitDirectory), true)
	if !hasGitDir || err != nil {
		return false, err
	}
	return true, nil
}

func isBareRepo(path string) (bool, error) {
	hasHeadFile, err := repoFileCheck(filepath.Join(path, headFile), false)
	if !hasHeadFile || err != nil {
		return false, err
	}
	hasObjectsDir, err := repoFileCheck(filepath.Join(path, objectsDirectory), true)
	if !hasObjectsDir || err != nil {
		return false, err
	}
	return true, nil
}

func isRepoWorkTree(path string) (bool, error) {
	hasGitFile, err := repoFileCheck(filepath.Join(path, gitFile), false)
	if !hasGitFile || err != nil {
		return false, err
	}
	return true, nil
}

func repoFileCheck(path string, isDir bool) (bool, error) {
	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("could not determine if path contains a Git repository: %w", err)
	}
	if isDir {
		return s.Mode().IsDir(), nil
	}
	return s.Mode().IsRegular(), nil
}
