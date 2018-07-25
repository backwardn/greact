package xhtml

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
	"text/template"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Node struct {
	Type      html.NodeType
	DataAtom  atom.Atom
	Data      string
	Namespace string
	Attr      []html.Attribute
	Children  []*Node
}

func Clone(n *html.Node, e *Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ch := &Node{
			Type:      c.Type,
			DataAtom:  c.DataAtom,
			Data:      c.Data,
			Namespace: c.Namespace,
			Attr:      make([]html.Attribute, len(c.Attr)),
		}
		copy(ch.Attr, c.Attr)
		e.Children = append(e.Children, ch)
		Clone(c, ch)
	}
}

func makeNode(n *Node) *ast.UnaryExpr {
	return &ast.UnaryExpr{
		Op: token.AND,
		X: &ast.CompositeLit{
			Type: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "xhtml"},
				Sel: &ast.Ident{Name: "Node"},
			},
			Elts: []ast.Expr{
				&ast.KeyValueExpr{
					Key: &ast.Ident{Name: "Type"},
					Value: &ast.BasicLit{
						Kind:  token.INT,
						Value: fmt.Sprint(n.Type),
					},
				},
				makeDataAtomField(n.DataAtom),
				&ast.KeyValueExpr{
					Key: &ast.Ident{Name: "Data"},
					Value: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("%q", n.Data),
					},
				},
				&ast.KeyValueExpr{
					Key: &ast.Ident{Name: "Namespace"},
					Value: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("%q", n.Namespace),
					},
				},
				&ast.KeyValueExpr{
					Key:   &ast.Ident{Name: "Attr"},
					Value: mkAttr(n.Attr),
				},
				&ast.KeyValueExpr{
					Key:   &ast.Ident{Name: "Children"},
					Value: makeChildren(n.Children),
				},
			},
		},
	}
}

func makeDataAtomField(a atom.Atom) *ast.KeyValueExpr {
	v := a.String()
	if v == "" {
		v = "Atom(0)"
	} else {
		v = strings.Title(v)
	}
	return &ast.KeyValueExpr{
		Key: &ast.Ident{Name: "DataAtom"},
		Value: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "atom"},
			Sel: &ast.Ident{Name: v},
		},
	}
}

func makeChildren(nodes []*Node) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: &ast.ArrayType{
			Elt: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "xhtml"},
					Sel: &ast.Ident{Name: "Node"},
				},
			},
		},
		Elts: makeArray(nodes),
	}
}

func makeArray(nodes []*Node) []ast.Expr {
	if nodes != nil {
		var ls []ast.Expr
		for _, v := range mkArrayElems(nodes) {
			ls = append(ls, v)
		}
		return ls
	}
	return []ast.Expr{}
}

func mkArrayElems(n []*Node) []*ast.UnaryExpr {
	var ls []*ast.UnaryExpr
	for _, v := range n {
		ls = append(ls, makeNode(v))
	}
	return ls
}

func mkAttr(a []html.Attribute) *ast.CompositeLit {
	ls := []ast.Expr{}
	for _, v := range a {
		ls = append(ls, mkAttrItem(v))
	}
	return &ast.CompositeLit{
		Type: &ast.ArrayType{
			Elt: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "html",
				},
				Sel: &ast.Ident{
					Name: "Attribute",
				},
			},
		},
		Elts: ls,
	}
}

func mkAttrItem(a html.Attribute) *ast.CompositeLit {
	return &ast.CompositeLit{
		Elts: []ast.Expr{
			&ast.KeyValueExpr{
				Key: &ast.Ident{Name: "Namespace"},
				Value: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", a.Namespace),
				},
			},
			&ast.KeyValueExpr{
				Key: &ast.Ident{Name: "Key"},
				Value: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", a.Key),
				},
			},
			&ast.KeyValueExpr{
				Key: &ast.Ident{Name: "Val"},
				Value: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", a.Val),
				},
			},
		},
	}
}

const fnTpl = `package {{.ctx.Package}}

import (
	"github.com/gernest/vected/lib/xhtml"
	"golang.org/x/net/html/atom"
)

func ({{.ctx.Recv}} {{.ctx.StructName}})Render()*xhtml.Node{
	return {{.node}}
}
`

var tpl = template.Must(template.New("n").Parse(fnTpl))

type Context struct {
	Package    string
	Recv       string
	StructName string
}

// GenerateRenderMethod using the given context, this returns a new go file with
// the Render method attached to the struct defined in ctx.
func GenerateRenderMethod(n *Node, ctx *Context) ([]byte, error) {
	var buf bytes.Buffer
	node := makeNode(n)
	fset := token.NewFileSet()
	err := format.Node(&buf, fset, node)
	if err != nil {
		return nil, err
	}
	nstr := buf.String()
	buf.Reset()
	err = tpl.Execute(&buf, map[string]interface{}{
		"ctx":  ctx,
		"node": nstr,
	})
	if err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}
