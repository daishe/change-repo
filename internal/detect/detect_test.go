package detect_test

import (
	"os"
	"strings"
	"testing"

	"github.com/daishe/change-repo/internal/detect"
)

type Loc struct {
	Variant detect.RepoVariant
	Path    string
}

var TestdataLocs = []Loc{
	{Variant: detect.NormalRepo, Path: "testdata/loc0/repo0"},
	{Variant: detect.NormalRepo, Path: "testdata/loc0/repo1"},
	{Variant: detect.NormalRepo, Path: "testdata/loc0/subset/repo0"},
	{Variant: detect.NormalRepo, Path: "testdata/loc1/repo0"},
	{Variant: detect.BareRepo, Path: "testdata/loc0/bare0"},
	{Variant: detect.WorkTreeDir, Path: "testdata/loc0/worktree0"},
	{Variant: detect.NonRepo, Path: "testdata/loc0/nonrepo"},
	{Variant: detect.NonRepo, Path: "testdata/loc0/subset"},
	{Variant: detect.NonRepo, Path: "testdata/loc0/nonrepo-bare-without-head"},
	{Variant: detect.NonRepo, Path: "testdata/loc0/nonrepo-bare-without-objects"},
	{Variant: detect.NonRepo, Path: "testdata/loc0/subset/nonrepo"},
}

func OverrideTestdataNames(t *testing.T) {
	t.Helper()
	detect.OverrideGitDirectory(t, "git-data")
	detect.OverrideGitFile(t, "git-file")
	detect.OverrideHeadFile(t, "head-file")
	detect.OverrideObjectsDirectory(t, "objects-data")
}

func p(path string) string {
	return strings.ReplaceAll(path, "/", string(os.PathSeparator))
}

func TestIsRepo(t *testing.T) { //nolint:paralleltest // this test overrides some global values and should not be run in parallel
	OverrideTestdataNames(t)

	for _, loc := range TestdataLocs {
		verdict, err := detect.IsRepo(p(loc.Path), detect.NormalRepo|detect.BareRepo|detect.WorkTreeDir)
		if err != nil {
			t.Errorf("unexpected error while checking if %q location is a repository: %v", loc.Path, err)
			continue
		}
		if verdict != loc.Variant {
			t.Errorf("invalid location %q qualification (expected: %s, got %s)", loc.Path, loc.Variant, verdict)
		}
	}
}
