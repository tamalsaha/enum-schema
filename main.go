package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"strings"
)

type Lookup struct {
	Target  string
	Types   map[string]types.Type
	comment *ast.CommentGroup
}

const hello = `
package main

import "fmt"

// append
func main() {
        // fmt
        fmt.Println("Hello, world")
        // main
        main, x := 1, 2
        // main
        print(main, x)
        // x
}
// x
`

var helloFile = hello

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", helloFile, parser.ParseComments)
	if err != nil {
		log.Fatal(err) // parse error
	}

	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("cmd/hello", fset, []*ast.File{f}, nil)
	if err != nil {
		log.Fatal(err) // type error
	}

	// Each comment contains a name.
	// Look up that name in the innermost scope enclosing the comment.
	for _, comment := range f.Comments {
		pos := comment.Pos()
		name := strings.TrimSpace(comment.Text())
		fmt.Printf("At %s,\t%q = ", fset.Position(pos), name)
		inner := pkg.Scope().Innermost(pos)
		if _, obj := inner.LookupParent(name, pos); obj != nil {
			fmt.Println(obj)
		} else {
			fmt.Println("not found")
		}
	}
}