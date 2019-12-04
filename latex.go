package extensions

import (
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type latexDelimiterProcessor struct {
}

func (p *latexDelimiterProcessor) IsDelimiter(b byte) bool {
	return b == '$'
}

func (p *latexDelimiterProcessor) CanOpenCloser(opener, closer *parser.Delimiter) bool {
	return opener.Char == closer.Char
}

func (p *latexDelimiterProcessor) OnMatch(consumes int) gast.Node {
	isInline := true
	if consumes > 1 {
		isInline = false
	}
	return ast.NewLatex(isInline)
}

var defaultLatexDelimiterProcessor = &latexDelimiterProcessor{}

type latexParser struct {
}

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
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()
	node := parser.ScanDelimiter(line, before, 1, defaultLatexDelimiterProcessor)
	if node == nil {
		node = parser.ScanDelimiter(line, before, 2, defaultLatexDelimiterProcessor)
	}
	if node == nil {
		return nil
	}

	node.Segment = segment.WithStop(segment.Start + node.OriginalLength)
	block.Advance(node.OriginalLength)
	pc.PushDelimiter(node)
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
	before := "<latex>"
	end := "</latex>"
	if n.IsInline {
		before = "<latex-inline>"
		end = "</latex-inline>"
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
