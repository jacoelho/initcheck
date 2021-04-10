package initcheck

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/packages"
)

//go:generate go run gen.go

func isInitFuncDeclaration(funcDecl *ast.FuncDecl) bool {
	return funcDecl.Name.Name == "init" && funcDecl.Recv.NumFields() == 0
}

func checkInitInPackage(pkg *packages.Package) {
	for _, fileAst := range pkg.Syntax {
		ast.Inspect(fileAst, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				return true
			}

			if isInitFuncDeclaration(funcDecl) {
				render(pkg.Fset, funcDecl)
			}

			return true
		})
	}
}

// render returns the pretty-print of the given node
func render(fset *token.FileSet, funcDecl *ast.FuncDecl) {
	name := fset.File(funcDecl.Pos()).Name()
	line := fset.Position(funcDecl.Pos()).Line

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, funcDecl); err != nil {
		panic(err)
	}

	fmt.Printf("%s:%d\n%s\n\n", name, line, buf.String())
}

// allImports returns all imports without duplication
func allImports(pkgs []*packages.Package) []*packages.Package {
	var all []*packages.Package
	var visit func(pkg *packages.Package)

	seen := make(map[*packages.Package]bool)

	visit = func(pkg *packages.Package) {
		if seen[pkg] {
			return
		}

		seen[pkg] = true

		var importsPaths []string
		for path := range pkg.Imports {
			importsPaths = append(importsPaths, path)
		}

		for _, path := range importsPaths {
			visit(pkg.Imports[path])
		}

		all = append(all, pkg)
	}

	for _, pkg := range pkgs {
		visit(pkg)
	}

	return all
}

func loadPackages(patterns []string) ([]*packages.Package, error) {
	config := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedDeps |
			packages.NeedModule | packages.NeedFiles | packages.NeedSyntax |
			packages.NeedTypes,
	}

	pkgs, err := packages.Load(config, patterns...)
	if err != nil {
		return nil, err
	}

	for _, p := range pkgs {
		for _, err := range p.Errors {
			return nil, err
		}
	}

	return pkgs, nil
}

// Run loads and check if imported packages contain an init function
func Run(patterns []string) error {
	pkgs, err := loadPackages(patterns)
	if err != nil {
		return err
	}

	all := allImports(pkgs)

	for _, pkg := range all {
		if isStdlib(pkg) {
			continue
		}

		checkInitInPackage(pkg)
	}

	return nil
}
