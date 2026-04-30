package web

import (
	"embed"
	"net/http"
)

//go:embed toolkit
var toolkit embed.FS

var (
	toolkit_fs     = http.FileServerFS(toolkit)
	toolkit_fs_dyn = http.StripPrefix("/toolkit", http.FileServer(http.Dir("./ui/web/toolkit/")))
)

func Toolkit(w http.ResponseWriter, r *http.Request) {
	toolkit_fs.ServeHTTP(w, r)
}

func ToolkitDyn(w http.ResponseWriter, r *http.Request) {
	toolkit_fs_dyn.ServeHTTP(w, r)
}
