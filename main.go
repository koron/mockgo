package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/koron-go/srcdom"
	"github.com/koron/mockgo/internal/common"
	"github.com/koron/mockgo/internal/mock1"
	"github.com/koron/mockgo/internal/mock2"
	"github.com/koron/mockgo/internal/mock3"
	"golang.org/x/tools/imports"
)

type mockTypeGenerator func(w io.Writer, mockTag, mockTypn, mockPkgn string, typ *srcdom.Type, pkg *srcdom.Package) error

type errs []error

func (e *errs) Append(err error) {
	*e = append(*e, err)
}

func (e errs) Error() string {
	if len(e) == 0 {
		return "no errors"
	}
	b := &strings.Builder{}
	fmt.Fprintln(b, "found some errors:")
	for i, err := range e {
		fmt.Fprintf(b, "#%d - %v\n", i+1, err)
	}
	return b.String()
}

func path2pkgname(path string) (string, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Base(p), nil
}

func mockFilename(typn string) string {
	// if mocktype name ends with "mock", truncate it for filename.
	base := strings.TrimSuffix(strings.ToLower(typn), "mock")
	if forTest {
		return base + "_mock_test.go"
	}
	return base + "_mock.go"
}

func generateMockType(outdir, mockTypn string, applyFormat bool, typ *srcdom.Type, pkg *srcdom.Package) error {
	pkgn, err := path2pkgname(outdir)
	if err != nil {
		return err
	}

	fname := mockFilename(mockTypn)
	fpath := filepath.Join(outdir, fname)
	f, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	bw := bufio.NewWriter(f)

	var w io.Writer = bw
	var bb *bytes.Buffer

	if applyFormat {
		bb = &bytes.Buffer{}
		w = bb
	}

	verbosef("writing %s for %s mock (%s)", fpath, typ.Name, mockTypn)
	err = mockTypeGen(w, "mock", mockTypn, pkgn, typ, pkg)
	if err != nil {
		f.Close()
		os.Remove(fpath)
		return err
	}

	if bb != nil {
		b, err := imports.Process(fname, bb.Bytes(), nil)
		if err != nil {
			f.Close()
			os.Remove(fpath)
			return err
		}
		bw.Write(b)
	}

	err = bw.Flush()
	if err != nil {
		f.Close()
		os.Remove(fpath)
		return err
	}
	err = f.Sync()
	if err != nil {
		f.Close()
		os.Remove(fpath)
		return err
	}
	err = f.Close()
	if err != nil {
		os.Remove(fpath)
		return err
	}
	return nil
}

func generateMockTypeAll(outdir string, typnames []string, pkg *srcdom.Package) error {
	var errs errs
	for _, typn := range typnames {
		var mockTypn string
		if n := strings.IndexRune(typn, ':'); n >= 0 {
			typn, mockTypn = typn[:n], typn[n+1:]
		}
		typ, ok := pkg.Type(typn)
		if !ok {
			err := fmt.Errorf("not found type:%s, skipped", typn)
			errs.Append(err)
			log.Print(err)
			continue
		}
		if mockTypn == "" {
			mockTypn = typ.Name
			if mockSuffix {
				mockTypn += "Mock"
			}
		}
		err := generateMockType(outdir, mockTypn, !noFormat, typ, pkg)
		if err != nil {
			err2 := fmt.Errorf("failed to generate mock for %s: %s", typ.Name, err)
			errs.Append(err2)
			log.Print(err2)
			continue
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

var (
	verbose    bool
	forTest    bool
	mockSuffix bool
	mockRev    int
	noFormat   bool
	version    bool

	mockTypeGen mockTypeGenerator
)

func determieMockTypeGenerator(mockRev int) error {
	switch mockRev {
	case 1:
		mockTypeGen = mock1.Generate
	case 2:
		mockTypeGen = mock2.Generate
	case 3:
		mockTypeGen = mock3.Generate
	default:
		return fmt.Errorf("unknow mock revision: %d", mockRev)
	}
	return nil
}

func gen() error {
	var (
		pkgname  string
		outdir   string
		typnames []string
	)
	flag.BoolVar(&forTest, "fortest", false, "generate mock for plain test, without +mock")
	flag.BoolVar(&mockSuffix, "mocksuffix", false, "add `Mock` suffix to generated mock types")
	flag.IntVar(&mockRev, "revision", 1, "mock revision (1-3)")
	flag.BoolVar(&noFormat, "noformat", false, "suppress goimports on generation mock code")
	flag.StringVar(&outdir, "outdir", ".", "output directory")
	flag.StringVar(&pkgname, "package", "", "package name")
	flag.BoolVar(&verbose, "verbose", false, "show verbose/debug messages to stderr")
	flag.BoolVar(&version, "version", false, "show version end exit")
	flag.Parse()

	typnames = flag.Args()
	common.ForTest = forTest

	if version {
		showVersion()
		return nil
	}

	// check options
	if pkgname == "" {
		return errors.New("need -package option")
	}
	if len(typnames) == 0 {
		return errors.New("need one or more type names")
	}
	if err := determieMockTypeGenerator(mockRev); err != nil {
		return err
	}

	// read source files, build srcdom.
	path := filepath.ToSlash(pkgname)
	if !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "../") {
		path = filepath.Join(build.Default.GOPATH, "src", pkgname)
	}
	pkg, err := srcdom.Read(path)
	if err != nil {
		return err
	}

	err = generateMockTypeAll(outdir, typnames, pkg)
	if err != nil {
		return err
	}
	verbosef("complete successfully")
	return nil
}

func verbosef(msg string, args ...interface{}) {
	if !verbose {
		return
	}
	log.Printf(msg, args...)
}

func showVersion() {
	fmt.Printf("mockgo version %s\n", Version)
}

func main() {
	err := gen()
	if err != nil {
		log.Fatal(err)
	}
}
