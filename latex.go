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
	"net/url"
	"os"
	"regexp"
)

var latexBlockRegexp = regexp.MustCompile(`^\$\$[^$]*\$\$`)
var latexInlineRegexp = regexp.MustCompile(`^\$[^$]*\$`)

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
	line, segment := block.PeekLine()
	isInline := false
	m := latexBlockRegexp.FindSubmatchIndex(line)
	if m == nil {
		m = latexInlineRegexp.FindSubmatchIndex(line)
		isInline = true
	}
	if m == nil {
		return nil
	}

	block.Advance(m[1])

	var value []byte
	if isInline {
		value = block.Value(text.NewSegment(segment.Start+1, segment.Start+m[1]-1))
	} else {
		value = block.Value(text.NewSegment(segment.Start+2, segment.Start+m[1]-2))
	}
	node := ast.NewLatex(isInline, value)

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

// LocalLatexJSON is the response from local Latex server
type LocalLatexJSON struct {
	SpeakText string `json:"speakText"`
	SVG       string `json:"svg"`
	Width     string `json:"width"`
	Height    string `json:"height"`
	Style     string `json:"style"`
}

func (r *LatexHTMLRenderer) renderLatex(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.Latex)
	value := url.QueryEscape(string(n.Value))

	if useLocalServer, ok := os.LookupEnv("HUGO_LATEX_USE_LOCAL"); ok {
		// use local server to render latex
		u, err := url.Parse(useLocalServer)
		if err != nil {
			fmt.Println(err)
			return gast.WalkContinue, nil
		}
		latexRes := &LocalLatexJSON{}
		query := "/?q=" + value
		if n.IsInline {
			query += "&inline=true"
		}

		err = utils.GetJSON(u.String()+query, latexRes)
		if err != nil {
			fmt.Println(err)
			return gast.WalkContinue, nil
		}

		if entering {
			w.WriteString(latexRes.SVG)
		}
		return gast.WalkContinue, nil
	}

	// use external service to get latex img
	before := "<figure><img src=\"https://math.now.sh?from="
	end := "\"/></figure>"
	if n.IsInline {
		before = "<img src=\"https://math.now.sh?inline="
		end = "\"/>"
	}
	if entering {
		w.WriteString(before)
		w.WriteString(value)
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
