package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

func parseDir(path string) (*ast.Package, error) {
	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, path, func(fi os.FileInfo) bool {
		if strings.HasSuffix(fi.Name(), "_test.go") {
			return false
		}
		return true
	}, 0)
	if err != nil {
		return nil, err
	}
	if n := len(pkgs); n > 1 {
		return nil, fmt.Errorf("found %d packages, expected just one", n)
	}
	for _, p := range pkgs {
		return p, nil
	}
	return nil, errors.New("no packages found, expected one")
}
