// Package ast defines AST nodes that represents extension's elements
package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
)

// A Latex struct represents a strikethrough of GFM text.
type Latex struct {
	gast.BaseInline
	IsInline bool
	Value    []byte
}

// Dump implements Node.Dump.
func (n *Latex) Dump(source []byte, level int) {
	m := map[string]string{
		"Inline": fmt.Sprintf("%v", n.IsInline),
		"Vaule":  fmt.Sprintf("%v", n.Value),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindLatex is a NodeKind of the Latex node.
var KindLatex = gast.NewNodeKind("Latex")

// Kind implements Node.Kind.
func (n *Latex) Kind() gast.NodeKind {
	return KindLatex
}

// NewLatex returns a new Latex node.
func NewLatex(isInline bool, value []byte) *Latex {
	return &Latex{
		IsInline: isInline,
		Value:    value,
	}
}
