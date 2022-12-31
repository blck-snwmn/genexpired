package main

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func buildMethod(reciverName, reciverType string) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: reciverName,
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: reciverType,
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Expired",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							&ast.Ident{Name: "now"},
						},
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "time"},
							Sel: &ast.Ident{Name: "Time"},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "bool",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.NOT,
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X:   &ast.Ident{Name: reciverName},
										Sel: &ast.Ident{Name: "expireAt"},
									},
									Sel: &ast.Ident{Name: "Before"},
								},
								Args: []ast.Expr{
									&ast.Ident{Name: "now"},
								},
							},
						},
					},
				},
			},
		},
	}
}

func main() {
	var source string
	flag.StringVar(&source, "source", "", "")
	flag.Parse()

	if source == "" {
		panic("no source")
	}
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		panic("no source")
	}
	fmt.Println("do1")

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, source, nil, parser.AllErrors+parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// ast.Print(fset, node)
	// f := buildMethod("c", "Claim")
	exists := false
	astutil.Apply(node, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "Expired" && x.Recv != nil {
				var typ string
				switch t := x.Recv.List[0].Type.(type) {
				case *ast.Ident:
					typ = t.Name
				case *ast.StarExpr:
					typ = t.X.(*ast.Ident).Name
				}

				exists = true
				fmt.Println("in")
				c.Replace(buildMethod(strings.ToLower(typ)[0:1], typ))
			}
		}
		return true
	})
	if !exists {
		// node.Decls = append(node.Decls, f)
	}
	ff, err := os.Create(source)
	if err != nil {
		panic(err)
	}

	format.Node(ff, fset, node)
}
