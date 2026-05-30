package web

import (
	"embed"
	"net/http"
)

type Toolkit struct {
	http.ServeMux
	glob *Glob
}

func NewToolkit(glob *Glob) http.Handler {
	tk := Toolkit{glob: glob}

	tk.Handle("/toolkit/", http.FileServerFS(toolkit))
	toolkit_example(glob, &tk.ServeMux)

	return &tk
}

func NewToolkitDyn(glob *Glob) http.Handler {
	tk := Toolkit{glob: glob}

	tk.Handle("/toolkit/", http.StripPrefix("/toolkit", http.FileServer(http.Dir("./ui/web/toolkit"))))
	toolkit_example(glob, &tk.ServeMux)

	return &tk
}

func toolkit_example(glob *Glob, mux *http.ServeMux) {
	tmpl, err := glob.Parse("index", "toolkit")
	if err != nil {
		return
	}

	page, err := MakePage(tmpl, nil)
	if err != nil {
		return
	}

	mux.Handle("/toolkit/{$}", page)
}

//go:embed toolkit
var toolkit embed.FS
