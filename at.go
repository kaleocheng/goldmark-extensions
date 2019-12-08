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

var atRegexp = regexp.MustCompile(`^@[-a-zA-Z0-9]*`)

type atParser struct {
}

var defaultAtParser = &atParser{}

// NewAtParser return a new InlineParser that parses
// at expressions.
func NewAtParser() parser.InlineParser {
	return defaultAtParser
}

func (s *atParser) Trigger() []byte {
	return []byte{'@'}
}

func (s *atParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, segment := block.PeekLine()
	m := atRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}
	block.Advance(m[1])
	node := ast.NewAt()
	node.AppendChild(node, gast.NewTextSegment(text.NewSegment(segment.Start+1, segment.Start+m[1])))
	return node
}

func (s *atParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// AtHTMLRenderer is a renderer.NodeRenderer implementation that
// renders At nodes.
type AtHTMLRenderer struct {
	html.Config
}

// NewAtHTMLRenderer returns a new AtHTMLRenderer.
func NewAtHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &AtHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *AtHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindAt, r.renderAt)
}

func (r *AtHTMLRenderer) renderAt(w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		w.WriteString("<span class=\"at\">@")
	} else {
		w.WriteString("</span>")
	}
	return gast.WalkContinue, nil
}

type at struct {
}

// At is an extension that allow you to use at expression like '@john' .
var At = &at{}

func (e *at) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewAtParser(), 500),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewAtHTMLRenderer(), 500),
	))
}
