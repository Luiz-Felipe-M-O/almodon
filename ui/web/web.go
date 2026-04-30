package web

import (
	_ "embed"
	"html/template"
)

//go:embed dist/almodon.html
var almodon string

var almodon_tmpl = template.Must(template.New("almodon").Parse(almodon))

func DefaultTemplate() *template.Template {
	return template.Must(almodon_tmpl.Clone())
}
