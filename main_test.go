package main

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron-go/srcdom"
)

func TestMockFilename(t *testing.T) {
	forTest = false
	for i, tc := range []struct {
		typn string
		want string
	}{
		{"foo", "foo_mock.go"},
		{"FOO", "foo_mock.go"},
		{"foomock", "foo_mock.go"},
		{"FOOMOCK", "foo_mock.go"},
		{"foo_mock", "foo__mock.go"},
	} {
		got := mockFilename(tc.typn)
		if got != tc.want {
			t.Errorf("failed #%d %+v: got=%s", i, tc, got)
		}
	}
}

type GenOptions struct {
	ForTest    bool
	MockSuffix bool
	MockRev    int
	NoFormat   bool

	Outdir    string
	Package   string
	Verbose   bool
	TypeNames []string
}

func newGenOptions(srcPkg, dstDir string, mockRev int, typs ...string) GenOptions {
	return GenOptions{
		MockRev:   mockRev,
		Outdir:    dstDir,
		Package:   srcPkg,
		TypeNames: typs,
	}
}

func (opts GenOptions) apply() {
	forTest = opts.ForTest
	mockSuffix = opts.MockSuffix
	mockRev = opts.MockRev
	noFormat = opts.NoFormat
	//outdir = opts.Outdir
	//pkgname = opts.Package
	//verbose = opts.Verbose
}

var muGen sync.Mutex

func runGen(opts GenOptions) error {
	muGen.Lock()
	defer muGen.Unlock()

	opts.apply()
	pkgname := opts.Package
	outdir := opts.Outdir
	typnames := opts.TypeNames

	// check options
	if pkgname == "" {
		return errors.New("need -package option")
	}
	if len(typnames) == 0 {
		return errors.New("need one or more type names")
	}
	if err := determieMockTypeGenerator(mockRev); err != nil {
		return fmt.Errorf("failed to determine mock: %w", err)
	}

	// read source files, build srcdom.
	path := filepath.ToSlash(pkgname)
	if !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "../") {
		path = filepath.Join(build.Default.GOPATH, "src", pkgname)
	}
	pkg, err := srcdom.Read(path)
	if err != nil {
		return fmt.Errorf("failed to read code: %w", err)
	}

	err = os.MkdirAll(outdir, 0777)
	if err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}
	err = generateMockTypeAll(outdir, typnames, pkg)
	if err != nil {
		return fmt.Errorf("failed to generation: %w", err)
	}
	return nil
}

func readFile(path string) ([]string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), "\n"), nil
}

func compareFile(t *testing.T, wantDir, gotDir string, filename string) {
	t.Helper()
	want, err := readFile(filepath.Join(wantDir, filename))
	if err != nil {
		t.Errorf("failed to read want file: %s", err)
		return
	}
	got, err := readFile(filepath.Join(gotDir, filename))
	if err != nil {
		t.Errorf("failed to read got file: %s", err)
		return
	}
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("unmatch files %s: -want +got\n%s", filename, d)
		return
	}
}

func TestMockTypeGen1(t *testing.T) {
	outdir := filepath.Join(t.TempDir(), "mock1_gen1")
	opts := newGenOptions("./testdata/pkg1", outdir, 1, "Foo")
	err := runGen(opts)
	if err != nil {
		t.Error(err)
	}
	compareFile(t, "./testdata/mock1_gen1", outdir, "foo_mock.go")
}

func TestMockTypeGen2(t *testing.T) {
	outdir := filepath.Join(t.TempDir(), "mock1_gen2")
	opts := newGenOptions("./testdata/pkg1", outdir, 2, "Foo")
	err := runGen(opts)
	if err != nil {
		t.Error(err)
	}
	compareFile(t, "./testdata/mock1_gen2", outdir, "foo_mock.go")
}

func TestMockTypeGen3(t *testing.T) {
	outdir := filepath.Join(t.TempDir(), "mock1_gen3")
	opts := newGenOptions("./testdata/pkg1", outdir, 3, "Foo")
	err := runGen(opts)
	if err != nil {
		t.Error(err)
	}
	compareFile(t, "./testdata/mock1_gen3", outdir, "foo_mock.go")
}
