package find_test

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/daishe/change-repo/internal/find"
	"github.com/google/go-cmp/cmp"
)

type Loc struct {
	Path  string
	IsDir bool
	Depth uint
}

var TestdataLocs = []Loc{
	{Path: "testdata/dir0", IsDir: true, Depth: 1},
	{Path: "testdata/dir0/dir0-0", IsDir: true, Depth: 2},
	{Path: "testdata/dir0/dir0-1", IsDir: true, Depth: 2},
	{Path: "testdata/dir0/file0-0", IsDir: false, Depth: 2},
	{Path: "testdata/dir1", IsDir: true, Depth: 1},
	{Path: "testdata/dir1/dir1-0", IsDir: true, Depth: 2},
	{Path: "testdata/dir1/dir1-0/dir1-0-0", IsDir: true, Depth: 3},
	{Path: "testdata/dir1/dir1-0/file1-0-0", IsDir: false, Depth: 3},
	{Path: "testdata/dir2", IsDir: true, Depth: 1},
	{Path: "testdata/dir2/file2-0", IsDir: false, Depth: 2},
	{Path: "testdata/dir2/file2-1", IsDir: false, Depth: 2},
	{Path: "testdata/file0", IsDir: false, Depth: 1},
}

func p(path string) string {
	return strings.ReplaceAll(path, "/", string(os.PathSeparator))
}

func TestFindDirs(t *testing.T) {
	search := find.Find{
		Root:     "testdata",
		Maxdepth: 2,
		OnError: func(err error) bool {
			t.Errorf("unexpected error: %v", err)
			return true
		},
	}

	want := []string{}
	for _, loc := range TestdataLocs {
		if !loc.IsDir || loc.Depth > search.Maxdepth {
			continue
		}
		want = append(want, p(loc.Path))
	}
	slices.Sort(want)

	gotMatched, gotNonMatched := search.AllDirs()
	slices.Sort(gotMatched)
	slices.Sort(gotNonMatched)

	if diff := cmp.Diff(want, gotMatched); diff != "" {
		t.Errorf("bad matched entries (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{}, gotNonMatched); diff != "" {
		t.Errorf("bad non-matched entries (-want +got):\n%s", diff)
	}
}
