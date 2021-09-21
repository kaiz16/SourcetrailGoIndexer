package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type AstVisitor struct {
	pkgDecl map[*ast.GenDecl]bool
	locals  map[string]int
	globals map[string]int
}

func newVisitor(f *ast.File) AstVisitor {
	decls := make(map[*ast.GenDecl]bool)
	for _, decl := range f.Decls {
		if v, ok := decl.(*ast.GenDecl); ok {
			decls[v] = true
		}
	}

	return AstVisitor{
		decls,
		make(map[string]int),
		make(map[string]int),
	}
}
func (a AstVisitor) index(filePath string) {
	// Create the AST by parsing filePath.
	localDecls, globalDecls := make(map[string]int), make(map[string]int)
	f, err := parser.ParseFile(indexer.prog, filePath, nil, 0)
	if err != nil {
		panic(err)
	}

	v := newVisitor(f)
	ast.Walk(v, f)
	for k, v := range v.locals {
		localDecls[k] += v
	}
	for k, v := range v.globals {
		globalDecls[k] += v
	}
	// Inspect the AST and print all identifiers and literals.
	// ast.Inspect(f, func(n ast.Node) bool {
	// 	var s string
	// 	switch x := n.(type) {
	// 	case *ast.BasicLit:
	// 		s = x.Value
	// 	case *ast.Ident:
	// 		s = x.Name
	// 	}
	// 	if s != "" {
	// 		fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
	// 	}
	// 	return true
	// })
}

func (v AstVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	switch d := n.(type) {
	case *ast.AssignStmt:
		if d.Tok != token.DEFINE {
			return v
		}
		for _, name := range d.Lhs {
			v.localDecl(name)
		}
	case *ast.RangeStmt:
		v.localDecl(d.Key)
		v.localDecl(d.Value)
	case *ast.FuncDecl:
		if d.Recv != nil {
			v.localDeclList(d.Recv.List)
		}
		v.localDeclList(d.Type.Params.List)
		if d.Type.Results != nil {
			v.localDeclList(d.Type.Results.List)
		}
	case *ast.GenDecl:
		if d.Tok != token.VAR {
			return v
		}
		for _, spec := range d.Specs {
			if value, ok := spec.(*ast.ValueSpec); ok {
				for _, name := range value.Names {
					if name.Name == "_" {
						continue
					}
					if v.pkgDecl[d] {
						indexer.registerGlobalVariable(d, name)
					} else {
						v.locals[name.Name]++
					}
				}
			}
		}
	}

	return v
}

func (v AstVisitor) localDecl(n ast.Node) {
	ident, ok := n.(*ast.Ident)
	if !ok {
		return
	}
	if ident.Name == "_" || ident.Name == "" {
		return
	}
	if ident.Obj != nil && ident.Obj.Pos() == ident.Pos() {
		v.locals[ident.Name]++
	}
}

func (v AstVisitor) localDeclList(fs []*ast.Field) {
	for _, f := range fs {
		for _, name := range f.Names {
			v.localDecl(name)
		}
	}
}
