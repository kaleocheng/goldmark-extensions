package extensions

import (
	"encoding/json"
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

var referRegexp = regexp.MustCompile(`^@refer\(.*\)`)

type referParser struct {
}

var defaultReferParser = &referParser{}

// NewReferParser return a new InlineParser that parses
// refer expressions.
func NewReferParser() parser.InlineParser {
	return defaultReferParser
}

func (s *referParser) Trigger() []byte {
	return []byte{'@'}
}

func (s *referParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, segment := block.PeekLine()
	m := referRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}

	block.Advance(m[1])

	var argvBytes []byte
	var argv []string
	var title []byte
	var url []byte
	argvBytes = append(argvBytes, []byte("[")...)
	argvBytes = append(argvBytes, block.Value(text.NewSegment(segment.Start+7, segment.Start+m[1]-1))...)
	argvBytes = append(argvBytes, []byte("]")...)
	_ = json.Unmarshal(argvBytes, &argv)
	length := len(argv)

	if length == 2 {
		title = []byte(argv[0])
		url = []byte(argv[1])
		node := ast.NewRefer(title, url)
		return node
	}

	if length == 1 {
		title = []byte(argv[0])
		url = []byte(argv[0])
		node := ast.NewRefer(title, url)
		return node
	}

	return nil
}

func (s *referParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// ReferHTMLRenderer is a renderer.NodeRenderer implementation that
// renders Refer nodes.
type ReferHTMLRenderer struct {
	html.Config
}

// NewReferHTMLRenderer returns a new ReferHTMLRenderer.
func NewReferHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &ReferHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *ReferHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindRefer, r.renderRefer)
}

func (r *ReferHTMLRenderer) renderRefer(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.Refer)
	title := string(n.Title)
	url := string(n.URL)
	before := "<a class=\"references\" href=\""
	end := "</a>"
	if entering {
		w.WriteString(before)
		w.WriteString(url)
		w.WriteString("\">")
		w.WriteString(title)
	} else {
		w.WriteString(end)
	}
	return gast.WalkContinue, nil
}

type refer struct {
}

// Refer is an extension that allow you to use refer expression like '$x^2$' .
var Refer = &refer{}

func (e *refer) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewReferParser(), 400),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewReferHTMLRenderer(), 500),
	))
}
