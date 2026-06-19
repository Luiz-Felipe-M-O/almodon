package web

import (
	"embed"
	"net/http"
)

type Toolkit struct{ http.ServeMux }

//go:embed toolkit
var source embed.FS

func NewToolkit(glob *Glob) *Toolkit {
	var tk Toolkit
	tk.Handle("/toolkit/", http.FileServerFS(source))
	style_example(glob, &tk.ServeMux)
	return &tk
}

func NewToolkitDyn(glob *Glob) *Toolkit {
	var tk Toolkit
	tk.Handle("/toolkit/", http.StripPrefix("/toolkit", http.FileServer(http.Dir("./ui/web/toolkit"))))
	style_example(glob, &tk.ServeMux)
	return &tk
}

func style_example(glob *Glob, mux *http.ServeMux) {
	if glob == nil {
		return
	}

	tmpl, err := glob.Parse("index", "toolkit")
	if err != nil {
		return
	}

	page, err := MakePage(tmpl, nil)
	if err != nil {
		return
	}

	mux.Handle("/toolkit/style/{$}", page)
}
