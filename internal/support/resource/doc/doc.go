package doc

import (
	"bytes"
	_ "embed"
	"errors"
	"go/ast"
	"go/doc"
	"go/doc/comment"
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
	Path      string
	Descript  template.HTML
	EndPoints []EndPoint

	buf bytes.Reader
}

var ErrDocNotFound = errors.New("resource doc not found")

func New(title string, r io.Reader) (*Doc, error) {
	pkg, rc, err := resource(r)
	if err != nil {
		return nil, err
	}

	var p comment.Parser
	var b strings.Builder

	if err := parse_content(&b, p.Parse(rc.Doc).Content); err != nil {
		return nil, err
	}
	root := template.HTML(b.String())

	var eps []EndPoint
	for _, m := range rc.Methods {
		ep, ok := NewEndPoint(m.Doc)
		if !ok {
			continue
		}

		eps = append(eps, ep)
	}
	if len(eps) == 0 {
		return nil, ErrDocNotFound
	}

	slices.SortFunc(eps, func(a, b EndPoint) int {
		return methods[a.Method] - methods[b.Method]
	})

	doc := Doc{
		Title:     title,
		Path:      pkg,
		Descript:  root,
		EndPoints: eps,
	}

	var buf bytes.Buffer
	if err := doc_tmpl.Execute(&buf, doc); err != nil {
		return nil, err
	}

	doc.buf = *bytes.NewReader(buf.Bytes())
	return &doc, nil
}

func (d *Doc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeContent(w, r, ".html", time.Time{}, &d.buf)
}

var doc_tmpl = template.Must(
	web.Base().
		Funcs(template.FuncMap{"lower": strings.ToLower}).
		Parse(web.MustText("doc/doc")),
)

var methods = map[string]int{
	"GET":    1,
	"POST":   2,
	"PATCH":  3,
	"PUT":    4,
	"DELETE": 5,
}

func resource(r io.Reader) (string, *doc.Type, error) {
	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "http.go", r, parser.ParseComments)
	if err != nil {
		return "", nil, err
	}

	d, err := doc.NewFromFiles(fset, []*ast.File{file}, "")
	if err != nil {
		return "", nil, err
	}

	var rc *doc.Type
	for _, ty := range d.Types {
		if ty.Name == "Resource" {
			rc = ty
			break
		}
	}
	if rc == nil {
		return "", nil, ErrDocNotFound
	}

	return d.Name, rc, nil
}
