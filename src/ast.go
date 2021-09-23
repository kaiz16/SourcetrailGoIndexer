package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type AstVisitor struct {
	globalGenDecls map[*ast.GenDecl]bool
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
	}
}
func (a AstVisitor) index(filePath string) {
	// Create the AST by parsing filePath.
	// localDecls, globalDecls := make(map[string]int), make(map[string]int)
	f, err := parser.ParseFile(indexer.prog, filePath, nil, 0)
	if err != nil {
		panic(err)
	}

	v := newVisitor(f)
	ast.Walk(v, f)
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
			ident := v.getIdentFromNode(name)
			if ident != nil {
				indexer.registerLocalSymbol(ident)
			}
		}
	case *ast.RangeStmt:
		key := v.getIdentFromNode(d.Key)
		if key != nil {
			indexer.registerLocalSymbol(key)
		}

		val := v.getIdentFromNode(d.Value)
		if val != nil {
			indexer.registerLocalSymbol(val)
		}
	case *ast.FuncDecl:
		if d.Recv != nil {
			for _, ident := range v.getIdentNodesFromFields(d.Recv.List) {
				if ident != nil {
					indexer.registerLocalSymbol(ident)
				}
			}
		}
		for _, ident := range v.getIdentNodesFromFields(d.Type.Params.List) {
			if ident != nil {
				indexer.registerLocalSymbol(ident)
			}
		}
		if d.Type.Results != nil {
			for _, ident := range v.getIdentNodesFromFields(d.Type.Results.List) {
				if ident != nil {
					indexer.registerLocalSymbol(ident)
				}
			}
		}
	case *ast.GenDecl:
		if d.Tok != token.VAR {
			return v
		}
		for _, ident := range v.getIdentNodesFromSpecs(d.Specs) {
			if ident != nil {
				// Global variable
				if v.globalGenDecls[d] {
					indexer.registerGlobalVariable(d, ident)
				} else {
					indexer.registerLocalSymbol(ident)
				}
			}
		}
	}

	return v
}

func (v AstVisitor) getIdentFromNode(n ast.Node) *ast.Ident {
	ident, ok := n.(*ast.Ident)
	if !ok {
		return nil
	}
	if ident.Name == "_" || ident.Name == "" {
		return nil
	}
	if ident.Obj != nil && ident.Obj.Pos() == ident.Pos() {
		return ident
	}
	return nil
}

func (v AstVisitor) getIdentNodesFromFields(fs []*ast.Field) []*ast.Ident {
	nodes := make([]*ast.Ident, 0)
	for _, f := range fs {
		for _, name := range f.Names {
			nodes = append(nodes, v.getIdentFromNode(name))
		}
	}
	return nodes
}

func (v AstVisitor) getIdentNodesFromSpecs(fs []ast.Spec) []*ast.Ident {
	nodes := make([]*ast.Ident, 0)
	for _, spec := range fs {
		if value, ok := spec.(*ast.ValueSpec); ok {
			for _, name := range value.Names {
				if name.Name == "_" {
					continue
				}
				nodes = append(nodes, v.getIdentFromNode(name))
			}
		}
	}
	return nodes
}
