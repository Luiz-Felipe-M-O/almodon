package client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alan-b-lima/almodon/ui/web"
)

func About(glob *web.Glob) (http.Handler, error) {
	page := struct {
		Contributors []contributor
		Associates   []associate
		CopyYear     string
	}{
		Contributors: peers,
		Associates:   associates,
	}

	orgn := 2025
	curr := time.Now().Year()

	if orgn != curr {
		page.CopyYear = fmt.Sprintf("%d-%d", orgn, curr)
	} else {
		page.CopyYear = strconv.Itoa(orgn)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := glob.Parse("index", "about")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page, err := web.MakePage(tmpl, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page.ServeHTTP(w, r)
	}), nil
}

type (
	contributor struct {
		Name   string
		GitHub string
	}

	associate struct {
		Name string
		HRef string
		Icon string
	}
)

var peers = []contributor{
	{Name: "Alan Barbosa Lima", GitHub: "alan-b-lima"},
	{Name: "Breno Augusto Braga Oliveira", GitHub: "bragabreno"},
	{Name: "Juan Pablo Ferreira Costa", GitHub: "juan-ferreirax"},
	{Name: "Luan Filipe Oliveira de Carvalho", GitHub: "Luan-11"},
	{Name: "Luann Moreira Fernandes de Oliveira", GitHub: "LuannMFO"},
	{Name: "Lucas Rocha Oliveira", GitHub: "Lucas-Rocha-Oliveira"},
	{Name: "Luiz Felipe Melo Oliveira", GitHub: "Luiz-Felipe-M-O"},
	{Name: "Mateus Oliveira Silva", GitHub: "MateusSilva06"},
	{Name: "Otávio Gomes Calazans", GitHub: "otaviogomes03"},
	{Name: "Rafael Gomes Silva", GitHub: "rafleGomes"},
	{Name: "Vitor Moisés Vieira Sales", GitHub: "VitorMozer9"},
}

var associates = []associate{
	{Name: "UFVJM", HRef: "https://portal.ufvjm.edu.br/", Icon: "https://portal.ufvjm.edu.br/dicom/central-de-conteudo/identidade-visual/marcas-ufvjm/vertical-sem-assinatura-colorida.png"},
	{Name: "DECOM", HRef: "https://decom.ufvjm.edu.br/", Icon: "https://www.decom.ufvjm.edu.br/dc2020/wp-content/uploads/2020/03/logo_topo.png"},
}
