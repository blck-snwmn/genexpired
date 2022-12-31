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
	m := map[string]bool{}
	astutil.Apply(node, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name == "Expired" && x.Recv != nil {
				var structName string
				switch t := x.Recv.List[0].Type.(type) {
				case *ast.Ident:
					structName = t.Name
				case *ast.StarExpr:
					structName = t.X.(*ast.Ident).Name
				}
				m[structName] = true
				fmt.Println("in")
				c.Replace(buildMethod(strings.ToLower(structName)[0:1], structName))
			}
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); !ok {
				return true
			}
			structName := x.Name.Name
			_, ok := m[structName]
			if !ok {
				// なければ追加する
				m[structName] = false
			}
			fmt.Println(structName)
		}
		return true
	})
	for structName, v := range m {
		if !v {
			node.Decls = append(node.Decls, buildMethod(strings.ToLower(structName)[0:1], structName))
		}
	}
	ff, err := os.Create(source)
	if err != nil {
		panic(err)
	}

	format.Node(ff, fset, node)
}
