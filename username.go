package extensions

import (
	"github.com/kaleocheng/goldmark-extensions/ast"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"regexp"
)

var usernameRegexp = regexp.MustCompile(`^@[-a-zA-Z0-9]*`)

type usernameParser struct {
}

var defaultUsernameParser = &usernameParser{}

// NewUsernameParser return a new InlineParser that parses
// username expressions.
func NewUsernameParser() parser.InlineParser {
	return defaultUsernameParser
}

func (s *usernameParser) Trigger() []byte {
	return []byte{'@'}
}

func (s *usernameParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, segment := block.PeekLine()
	m := usernameRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}
	block.Advance(m[1])
	node := ast.NewUsername()
	node.AppendChild(node, gast.NewTextSegment(text.NewSegment(segment.Start+1, segment.Start+m[1])))
	return node
}

func (s *usernameParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// UsernameHTMLRenderer is a renderer.NodeRenderer implementation that
// renders Username nodes.
type UsernameHTMLRenderer struct {
	html.Config
}

// NewUsernameHTMLRenderer returns a new UsernameHTMLRenderer.
func NewUsernameHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &UsernameHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *UsernameHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindUsername, r.renderUsername)
}

func (r *UsernameHTMLRenderer) renderUsername(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		w.WriteString("<username>")
	} else {
		w.WriteString("</username>")
	}
	return gast.WalkContinue, nil
}

type username struct {
}

// Username is an extension that allow you to use username expression like '@john' .
var Username = &username{}

func (e *username) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewUsernameParser(), 500),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewUsernameHTMLRenderer(), 500),
	))
}
