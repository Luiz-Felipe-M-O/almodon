package web

import (
	"embed"
	"net/http"
)

type Toolkit struct {
	http.ServeMux
}

type Script struct {
	http.Handler
}

//go:embed toolkit script
var source embed.FS

func NewToolkit(glob *Glob) *Toolkit {
	var tk Toolkit
	tk.Handle("/toolkit/", http.FileServerFS(source))

	if glob != nil {
		toolkit_example(glob, &tk.ServeMux)
	}
	return &tk
}

func NewToolkitDyn(glob *Glob) *Toolkit {
	var tk Toolkit
	tk.Handle("/toolkit/", http.StripPrefix("/toolkit", http.FileServer(http.Dir("./ui/web/toolkit"))))

	if glob != nil {
		toolkit_example(glob, &tk.ServeMux)
	}
	return &tk
}

func NewScript() *Script {
	return &Script{Handler: http.FileServerFS(source)}
}

func NewScriptDyn() *Script {
	return &Script{
		Handler: http.StripPrefix("/script", http.FileServer(http.Dir("./ui/web/script"))),
	}
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
