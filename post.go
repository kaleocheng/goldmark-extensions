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
	"path/filepath"
	"regexp"
	"strings"
)

var postRegexp = regexp.MustCompile(`^@post\(".*"\)`)

type postParser struct {
}

var defaultPostParser = &postParser{}

// NewPostParser return a new InlineParser that parses
// post expressions.
func NewPostParser() parser.InlineParser {
	return defaultPostParser
}

func (s *postParser) Trigger() []byte {
	return []byte{'@'}
}

func (s *postParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	line, segment := block.PeekLine()
	m := postRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}

	block.Advance(m[1])

	url := block.Value(text.NewSegment(segment.Start+7, segment.Start+m[1]-2))
	cleanedTitle := strings.Split(filepath.Clean(string(url)), "/")
	title := []byte(strings.Title(strings.ToLower(strings.Replace(cleanedTitle[len(cleanedTitle)-1], "-", " ", -1))))
	node := ast.NewPost(title, url)

	return node
}

func (s *postParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

// PostHTMLRenderer is a renderer.NodeRenderer implementation that
// renders Post nodes.
type PostHTMLRenderer struct {
	html.Config
}

// NewPostHTMLRenderer returns a new PostHTMLRenderer.
func NewPostHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &PostHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *PostHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindPost, r.renderPost)
}

func (r *PostHTMLRenderer) renderPost(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.Post)
	title := string(n.Title)
	url := string(n.URL)
	before := "<a href=\""
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

type post struct {
}

// Post is an extension that allow you to use post expression like '$x^2$' .
var Post = &post{}

func (e *post) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewPostParser(), 400),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewPostHTMLRenderer(), 500),
	))
}
