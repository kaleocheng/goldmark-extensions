package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

// A At struct represents a At
type At struct {
	gast.BaseInline
}

// Dump implements Node.Dump.
func (n *At) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindAt is a NodeKind of the At node.
var KindAt = gast.NewNodeKind("At")

// Kind implements Node.Kind.
func (n *At) Kind() gast.NodeKind {
	return KindAt
}

// NewAt returns a new At node.
func NewAt() *At {
	return &At{}
}
