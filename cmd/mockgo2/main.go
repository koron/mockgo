package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

var (
	verbose   bool
	pkgdir    string
	typeNames []string
)

func run() error {
	flag.StringVar(&pkgdir, "pkgdir", ".", "target package directory")
	flag.BoolVar(&verbose, "verbose", false, "show verbose/debug messages")
	flag.Parse()
	typeNames = flag.Args()
	if len(typeNames) == 0 {
		return errors.New("no types. require one or more types to mock")
	}

	pkg, err := parseDir(pkgdir)
	if err != nil {
		return err
	}

	_ = pkg
	for _, n := range typeNames {
		_ = n
	}
	return nil
}

func getOnePackage(pkgs map[string]*ast.Package) (*ast.Package, error) {
	if n := len(pkgs); n > 1 {
		return nil, fmt.Errorf("found %d packages, expected just one", n)
	}
	for _, p := range pkgs {
		return p, nil
	}
	return nil, errors.New("no packages found, expected one")
}

func astFiles(pkg *ast.Package) []*ast.File {
	files := make([]*ast.File, 0, len(pkg.Files))
	for _, f := range pkg.Files {
		files = append(files, f)
	}
	return files
}

func parseDir(path string) (*ast.Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, func(fi os.FileInfo) bool {
		if strings.HasSuffix(fi.Name(), "_test.go") {
			return false
		}
		return true
	}, 0)
	if err != nil {
		return nil, err
	}
	pkg, err := getOnePackage(pkgs)
	if err != nil {
		return nil, err
	}

	conf := types.Config{Importer: importer.Default()}
	_, err = conf.Check(path, fset, astFiles(pkg), nil)
	if err != nil {
		return nil, err
	}

	return pkg, nil
}
