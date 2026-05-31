package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alan-b-lima/almodon/internal/domain/auth"
	"github.com/alan-b-lima/almodon/internal/domain/user"
	"github.com/alan-b-lima/almodon/pkg/uuid"

	"github.com/alan-b-lima/almodon/ui/web"
)

func main() {
	users := &UserView{
		Users: []user.Result{
			{
				UUID:    uuid.NewUUIDv7(),
				SIAPE:   "0000000",
				Name:    "Alan Barbosa Lima",
				Email:   "alan-lima.al@ufvjm.edu.br",
				Role:    auth.Maintainer,
				Logged:  true,
				Created: time.Now(),
				Updated: time.Now(),
			},
			{
				UUID:    uuid.NewUUIDv7(),
				SIAPE:   "0000001",
				Name:    "Breno Augusto Braga Oliveira",
				Email:   "breno.augusto@ufvjm.edu.br",
				Role:    auth.Admin,
				Logged:  false,
				Created: time.Now(),
				Updated: time.Now(),
			},
			{
				UUID:    uuid.NewUUIDv7(),
				SIAPE:   "0000002",
				Name:    "Luiz Felipe Melo Oliveira",
				Email:   "luiz.melo@ufvjm.edu.br",
				Role:    auth.Promoted,
				Logged:  true,
				Created: time.Now(),
				Updated: time.Now(),
			},
			{
				UUID:    uuid.NewUUIDv7(),
				SIAPE:   "0000003",
				Name:    "Juan Pablo Ferreira Costa",
				Email:   "juan-pablo.jp@ufvjm.edu.br",
				Role:    auth.User,
				Logged:  false,
				Created: time.Now(),
				Updated: time.Now(),
			},
		},
	}

	http.Handle("/", users)
	http.Handle("/toolkit/", web.ToolkitDyn())

	fmt.Println("server listening at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

type UserView struct {
	Users []user.Result
}

func (u *UserView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl, err := web.Base().Parse(`
		{{ define "main" }}
		<main>
		  <ol>
			  {{ range . }}
		    <li>
				  <ul>
				    <li>{{ .SIAPE }}</li>
			      <li>{{ .Name }}</li>
			      <li>{{ .Email }}</li>
					</ul>
				</li>
				{{ end }}
			</ol>			
		</main>
		{{ end }}
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if err := tmpl.Execute(w, u.Users); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}