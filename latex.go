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
)

var latexRegexp = regexp.MustCompile(`^\$\$?(.|\n)+?\$\$?`)
var defaultLatexParser = &latexParser{}

// NewLatexParser return a new InlineParser that parses
// latex expressions.
func NewLatexParser() parser.InlineParser {
	return defaultLatexParser
}

func (s *latexParser) Trigger() []byte {
	return []byte{'$'}
}

func (s *latexParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {

	line, segment := block.PeekLine()
	m := latexRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}
	block.Advance(m[1])
	node := ast.NewLatex()
	node.AppendChild(node, gast.NewTextSegment(text.NewSegment(segment.Start+1, segment.Start+m[1])))
	return node
}

func (s *latexParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// LatexHTMLRenderer is a renderer.NodeRenderer implementation that
// renders Latex nodes.
type LatexHTMLRenderer struct {
	html.Config
}

// NewLatexHTMLRenderer returns a new LatexHTMLRenderer.
func NewLatexHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &LatexHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *LatexHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindLatex, r.renderLatex)
}

func (r *LatexHTMLRenderer) renderLatex(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.Latex)
	before := "<figure><img src=\"https://math.now.sh?from="
	end := "\"/></figure>"
	if n.IsInline {
		before = "<img src=\"https://math.now.sh?inline="
		end = "\"/>"
	}
	if entering {
		w.WriteString(before)
	} else {
		w.WriteString(end)
	}
	return gast.WalkContinue, nil
}

type latex struct {
}

// Latex is an extension that allow you to use latex expression like '$x^2$' .
var Latex = &latex{}

func (e *latex) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewLatexParser(), 500),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewLatexHTMLRenderer(), 500),
	))
}
