package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/alan-b-lima/almodon/internal/api"
	"github.com/alan-b-lima/almodon/internal/server"
	"github.com/alan-b-lima/almodon/internal/support/middleware"
)

const (
	almodon = `Almodon ` + version + ` ` + runtime.GOOS + `/` + runtime.GOARCH + copyright
	version = `v0.0.1`
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "version" {
		fmt.Println(almodon)
		return
	}

	addr := ":4545"
	if len(os.Args) >= 2 {
		addr = os.Args[1]
	}

	api, err := api.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := api.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	var mux http.ServeMux

	server, err := server.New(addr, Logger(&mux))
	if err != nil {
		fmt.Println(err)
		return
	}

	mux.Handle("/", api)
	mux.HandleFunc("/terminate", func(w http.ResponseWriter, r *http.Request) {
		shutdown(server)
	})

	go SignalShutdown(server)

	fmt.Printf("server listening at http://%s\n", strings.Replace(server.Addr().String(), "[::]", "localhost", 1))
	if err := server.Serve(); err != nil {
		fmt.Println(err)
		return
	}

	<-server.Done()

	fmt.Printf("%s %s\n", time.Now().Format(time.DateTime), "server powering off...")
	time.Sleep(time.Second)
}

func SignalShutdown(server *server.Server) {
	<-server.Signal()
	fmt.Print("\r")

	if err := shutdown(server); err != nil {
		fmt.Println(err)
	}
}

func shutdown(server *server.Server) error {
	fmt.Printf("%s %s\n", time.Now().Format(time.DateTime), "starting server shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	errs := make(chan error, 1)
	go func() {
		errs <- server.Shutdown(ctx)
	}()

	select {
	case <-ctx.Done():
		server.CancelOngoing()
		time.Sleep(5 * time.Second)

	case <-server.Signal():

	case err := <-errs:
		if err == nil {
			return nil
		}
	}

	fmt.Printf("%s %s\n", time.Now().Format(time.DateTime), "starting forceful server shutdown...")
	return server.ForceShutdown()
}

func Logger(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.NewResponseWriterWithStatus(w)

		h.ServeHTTP(rw, r)

		fmt.Printf("%s %d %s %s %s\n", time.Now().Format(time.DateTime), rw.StatusCode(), r.RemoteAddr, r.Method, r.URL)
	}
}

const copyright = `
Copyright (C) 2026 Alan Lima

Esse programa está licenciado sob a Licença GPLv3 (GNU General Public License
versão 3). Este programa é distribuído na esperança de ser útil, mas SEM
NENHUMA GARANTIA; sem mesmo a garantia implícita de COMERCIALIZAÇÃO ou
ADEQUAÇÃO A UM PROPÓSITO ESPECÍFICO. Consulte a Licença GPLv3 para mais
detalhes.

Você deve ter recebido uma cópia da Licença GPLv3 junto com este programa. Se
não, veja <https://www.gnu.org/licenses/>.`
