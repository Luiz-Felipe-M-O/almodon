package doc

import (
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

	"github.com/alan-b-lima/almodon/ui/web"
)

type Doc struct {
	Title     string
	Path      string
	Descript  template.HTML
	EndPoints []EndPoint

	page http.Handler
}

var ErrDocNotFound = errors.New("resource doc not found")

func NewDoc(glob *web.Glob, title string, r io.Reader) (*Doc, error) {
	pkg, rc, err := resource(r)
	if err != nil {
		return nil, err
	}

	html, err := web.GoComment(rc.Doc)
	if err != nil {
		return nil, err
	}

	var eps []EndPoint
	for _, m := range rc.Methods {
		ep, err := NewEndPoint(m.Doc)
		if err != nil {
			if err == ErrNotRoute {
				continue
			}

			return nil, err
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
		Descript:  html,
		EndPoints: eps,
	}

	tmpl, err := glob.Clone().Func("lower", strings.ToLower).Parse("index", "doc/doc")
	if err != nil {
		return nil, err
	}

	doc.page, err = web.MakePage(tmpl, doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (d *Doc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d.page.ServeHTTP(w, r)
}

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
