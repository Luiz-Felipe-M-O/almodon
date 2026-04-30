package doc

import (
	"bytes"
	_ "embed"
	"errors"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/alan-b-lima/almodon/ui/web"
)

type Doc struct {
	Title     string
	EndPoints []EndPoint

	bytes []byte
}

var ErrDocNotFound = errors.New("resource doc not found")

func New(title string, r io.Reader) (*Doc, error) {
	resource, err := resource(r)
	if err != nil {
		return nil, err
	}

	var eps []EndPoint
	for _, m := range resource.Methods {
		ep, ok := NewEndPoint(m.Doc)
		if !ok {
			continue
		}

		eps = append(eps, ep)
	}

	slices.SortFunc(eps, func(a, b EndPoint) int {
		return methods[a.Method] - methods[b.Method]
	})

	doc := Doc{
		Title:     title,
		EndPoints: eps,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, doc); err != nil {
		return nil, err
	}

	doc.bytes = buf.Bytes()
	return &doc, nil
}

func (d *Doc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, "doc.html", time.Time{}, d.ReadSeeker())
}

func (d *Doc) ReadSeeker() io.ReadSeeker {
	return bytes.NewReader(d.bytes)
}

var tmpl = template.Must(parse_chain(
	web.DefaultTemplate().Funcs(template.FuncMap{"lower": strings.ToLower}),
	`{{ define "head" }}<link rel="stylesheet" href="/toolkit/style/doc.css">{{ end }}`,
	doc_text,
))

//go:embed template.html
var doc_text string

func parse_chain(tmpl *template.Template, texts ...string) (*template.Template, error) {
	for _, text := range texts {
		var err error

		tmpl, err = tmpl.Parse(text)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

var methods = map[string]int{
	"GET":    1,
	"POST":   2,
	"PATCH":  3,
	"PUT":    4,
	"DELETE": 5,
}

func resource(r io.Reader) (*doc.Type, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "http.go", r, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	d, err := doc.NewFromFiles(fset, []*ast.File{file}, "")
	if err != nil {
		return nil, err
	}

	var rc *doc.Type
	for _, ty := range d.Types {
		if ty.Name == "Resource" {
			rc = ty
			break
		}
	}
	if rc == nil {
		return nil, ErrDocNotFound
	}

	return rc, nil
}
