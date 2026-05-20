package web

import (
	"embed"
	"html/template"
	"net/http"
	"net/url"
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

func init_toolkit() http.Handler {
	tmpl := template.Must(Base().Parse(MustText("toolkit")))
	page, err := Page(tmpl, nil)
	if err != nil {
		panic(err)
	}

	return page
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

func init_icons() http.Handler {
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

	page, err := Page(tmpl, display)
	if err != nil {
		panic(err)
	}
	return page
}
