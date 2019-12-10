// Package ast defines AST nodes that represents extension's elements
package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
)

// A Date struct represents date url
type Date struct {
	gast.BaseInline
	Value []byte
}

// Dump implements Node.Dump.
func (n *Date) Dump(source []byte, level int) {
	m := map[string]string{
		"Value": fmt.Sprintf("%v", n.Value),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindDate is a NodeKind of the Date node.
var KindDate = gast.NewNodeKind("Date")

// Kind implements Node.Kind.
func (n *Date) Kind() gast.NodeKind {
	return KindDate
}

// NewDate returns a new Date node.
func NewDate(value []byte) *Date {
	return &Date{
		Value: value,
	}
}
