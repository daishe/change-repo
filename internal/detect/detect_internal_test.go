package detect

import "testing"

func OverrideGitDirectory(t *testing.T, new string) {
	t.Helper()
	src := gitDirectory
	t.Cleanup(func() { gitDirectory = src })
	gitDirectory = new
}

func OverrideGitFile(t *testing.T, new string) {
	t.Helper()
	src := gitFile
	t.Cleanup(func() { gitFile = src })
	gitFile = new
}

func OverrideHeadFile(t *testing.T, new string) {
	t.Helper()
	src := headFile
	t.Cleanup(func() { headFile = src })
	headFile = new
}

func OverrideObjectsDirectory(t *testing.T, new string) {
	t.Helper()
	src := objectsDirectory
	t.Cleanup(func() { objectsDirectory = src })
	objectsDirectory = new
}
