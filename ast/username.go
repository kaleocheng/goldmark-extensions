package ast

import (
	gast "github.com/yuin/goldmark/ast"
)

// A Username struct represents a username
type Username struct {
	gast.BaseInline
}

// Dump implements Node.Dump.
func (n *Username) Dump(source []byte, level int) {
	gast.DumpHelper(n, source, level, nil, nil)
}

// KindUsername is a NodeKind of the Username node.
var KindUsername = gast.NewNodeKind("Username")

// Kind implements Node.Kind.
func (n *Username) Kind() gast.NodeKind {
	return KindUsername
}

// NewUsername returns a new Username node.
func NewUsername() *Username {
	return &Username{}
}
