// Package ast defines AST nodes that represents extension's elements
package ast

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
)

// A Post struct represents a post url in blog
type Post struct {
	gast.BaseInline
	Title []byte
	URL   []byte
}

// Dump implements Node.Dump.
func (n *Post) Dump(source []byte, level int) {
	m := map[string]string{
		"Title": fmt.Sprintf("%v", n.Title),
		"URL":   fmt.Sprintf("%v", n.URL),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindPost is a NodeKind of the Post node.
var KindPost = gast.NewNodeKind("Post")

// Kind implements Node.Kind.
func (n *Post) Kind() gast.NodeKind {
	return KindPost
}

// NewPost returns a new Post node.
func NewPost(title []byte, url []byte) *Post {
	return &Post{
		Title: title,
		URL:   url,
	}
}
