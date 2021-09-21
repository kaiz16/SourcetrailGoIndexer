package main

import (
	"flag"
	"go/token"
	"os"
)

var indexer Indexer
var visitor AstVisitor

func main() {

	path, _ := os.Getwd()
	pkgPath := path + "/../test"
	debug := false
	flag.StringVar(&pkgPath, "pkgPath", pkgPath, "The absolute path for target package. Redirect to the example folder by default.\n")
	flag.BoolVar(&debug, "debug", false, "Print log or not.")
	flag.Parse()

	indexer.DatabasePath = pkgPath + "/cg.srctrldb"
	fset := token.NewFileSet() // positions are relative to fset

	indexer.prog = fset
	indexer.Open()
	defer indexer.Close()
	visitor.index(pkgPath + "/main.go")
	// Resolve call graphs
	cg(pkgPath)
}
