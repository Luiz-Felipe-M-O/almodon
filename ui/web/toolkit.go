package web

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

//go:embed toolkit
var toolkit embed.FS

func Toolkit() http.Handler {
	return &toolkit_mux_emb
}

func ToolkitDyn() http.Handler {
	return &toolkit_mux_dyn
}

var (
	toolkit_mux_emb http.ServeMux
	toolkit_mux_dyn http.ServeMux
)

func init() {
	toolkit_mux_emb.Handle("/", http.FileServerFS(toolkit))
	toolkit_mux_dyn.Handle("/", http.StripPrefix("/toolkit", http.FileServer(http.Dir("./ui/web/toolkit"))))

	toolkit_handler := init_toolkit()
	toolkit_mux_emb.Handle("/toolkit/{$}", toolkit_handler)
	toolkit_mux_dyn.Handle("/toolkit/{$}", toolkit_handler)

	icon_handler := init_icons()
	toolkit_mux_emb.Handle("/toolkit/assets/icons/{$}", icon_handler)
	toolkit_mux_dyn.Handle("/toolkit/assets/icons/{$}", icon_handler)
}

func init_toolkit() http.HandlerFunc {
	tmpl := template.Must(Base().Parse(MustText("toolkit")))

	var b bytes.Buffer
	if err := tmpl.Execute(&b, nil); err != nil {
		panic(err)
	}

	reader := bytes.NewReader(b.Bytes())
	build := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, ".html", build, reader)
	}
}

var icons = map[string]string{
	"about.svg":        "Sobre",
	"account.svg":      "Conta",
	"box.svg":          "Caixa",
	"close.svg":        "Fechar",
	"double-arrow.svg": "Retornar",
	"gear.svg":         "Engrenagem",
	"group.svg":        "Grupo",
	"history.svg":      "Histórico",
	"home.svg":         "Início",
	"login.svg":        "Entrar",
	"logout.svg":       "Sair",
	"menu.svg":         "Menu",
	"refresh.svg":      "Recarregar",
}

type icon struct {
	Src, Name  string
	Documented bool
}

const (
	icon_dir   = "toolkit/assets/icons"
	icon_route = "/" + icon_dir
)

func init_icons() http.HandlerFunc {
	dir, err := toolkit.ReadDir(icon_dir)
	if err != nil {
		panic(err)
	}

	display := make([]icon, 0, len(dir))

	for _, file := range dir {
		url, err := url.JoinPath(icon_route, file.Name())
		if err != nil {
			panic(err)
		}

		if name, in := icons[file.Name()]; in {
			display = append(display, icon{
				Src:        url,
				Name:       name,
				Documented: true,
			})
		} else {
			display = append(display, icon{
				Src:        url,
				Name:       file.Name(),
				Documented: false,
			})
		}
	}

	tmpl := template.Must(Base().Parse(MustText("toolkit/icons")))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, display); err != nil {
		panic(err)
	}

	reader := bytes.NewReader(buf.Bytes())
	build := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, ".html", build, reader)
	}
}
