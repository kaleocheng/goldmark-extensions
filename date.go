package extensions

import (
	"fmt"
	"github.com/kaleocheng/goldmark-extensions/ast"
	"github.com/kaleocheng/goldmark-extensions/utils"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"regexp"
)

var dateRegexp = regexp.MustCompile(`^@date\(.*\)`)

type dateParser struct {
}

var defaultDateParser = &dateParser{}

// NewDateParser return a new InlineParser that parses
// date expressions.
func NewDateParser() parser.InlineParser {
	return defaultDateParser
}

func (s *dateParser) Trigger() []byte {
	return []byte{'@'}
}

func (s *dateParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, segment := block.PeekLine()
	m := dateRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}

	var argv []string
	err := utils.Argvs(block.Value(text.NewSegment(segment.Start+6, segment.Start+m[1]-1)), &argv)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if len(argv) != 1 {
		return nil
	}

	block.Advance(m[1])

	value := []byte(argv[0])
	node := ast.NewDate(value)
	return node
}

func (s *dateParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// DateHTMLRenderer is a renderer.NodeRenderer implementation that
// renders Date nodes.
type DateHTMLRenderer struct {
	html.Config
}

// NewDateHTMLRenderer returns a new DateHTMLRenderer.
func NewDateHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &DateHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *DateHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindDate, r.renderDate)
}

func (r *DateHTMLRenderer) renderDate(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.Date)
	before := "<span class=\"date " + string(n.Value) + "\">"
	end := "</span>"
	if entering {
		w.WriteString(before)
	} else {
		w.WriteString(end)
	}
	return gast.WalkContinue, nil
}

type date struct {
}

// Date is an extension that allow you to use date expression like '$x^2$' .
var Date = &date{}

func (e *date) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewDateParser(), 400),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewDateHTMLRenderer(), 500),
	))
}
