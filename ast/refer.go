// Package ast defines AST nodes that represents extension's elements
package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
)

// A Refer struct represents refer url
type Refer struct {
	gast.BaseInline
	Title []byte
	URL   []byte
}

// Dump implements Node.Dump.
func (n *Refer) Dump(source []byte, level int) {
	m := map[string]string{
		"Title": fmt.Sprintf("%v", n.Title),
		"URL":   fmt.Sprintf("%v", n.URL),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindRefer is a NodeKind of the Refer node.
var KindRefer = gast.NewNodeKind("Refer")

// Kind implements Node.Kind.
func (n *Refer) Kind() gast.NodeKind {
	return KindRefer
}

// NewRefer returns a new Refer node.
func NewRefer(title []byte, url []byte) *Refer {
	return &Refer{
		Title: title,
		URL:   url,
	}
}
