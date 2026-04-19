package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/alan-b-lima/almodon/internal/api"
	"github.com/alan-b-lima/almodon/internal/server"
	"github.com/alan-b-lima/almodon/internal/support/middleware"
)

func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "version":
			fmt.Println(version)
			return

		case "copyright", "legal":
			fmt.Println(legal)
			return

		case "help":
			fmt.Println(help)
			return
		}
	}

	addr := ":4545"
	if len(os.Args) >= 2 {
		addr = os.Args[1]
	}

	if err := Main(addr); err != nil {
		fmt.Println(err)
	}
}

func Main(addr string) error {
	api, err := api.New()
	if err != nil {
		return err
	}
	defer func() {
		if err := api.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	log := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	var mux http.ServeMux

	server, err := server.New(addr, Logger(log, &mux))
	if err != nil {
		return err
	}

	mux.Handle("/", api)
	mux.HandleFunc("/terminate", func(w http.ResponseWriter, r *http.Request) {
		go shutdown(server)
		w.WriteHeader(http.StatusNoContent)
	})

	go SignalShutdown(server)

	fmt.Println("server listening at http://" + strings.Replace(server.Addr().String(), "[::]", "localhost", 1))
	if err := server.Serve(); err != nil {
		return err
	}

	<-server.Done()

	fmt.Printf("%s %s\n", time.Now().Format(time.DateTime), "server powering off...")
	return nil
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

func Logger(log *log.Logger, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := middleware.NewResponseWriterWithStatus(w)

		h.ServeHTTP(rw, r)

		log.Printf("%d %s %s %s\n", rw.StatusCode(), r.RemoteAddr, r.Method, r.URL)
	}
}

const (
	version = `Almodon ` + tag + ` ` + runtime.GOOS + `/` + runtime.GOARCH
	legal   = version + copyright
	help    = version + "\n\n" + `Work in Progress`
)

const tag = `v0.0.1`

const copyright = `
Copyright (C) 2026 Alan Lima

Esse programa está licenciado sob a Licença GPLv3 (GNU General Public License
versão 3). Este programa é distribuído na esperança de ser útil, mas SEM
NENHUMA GARANTIA; sem mesmo a garantia implícita de COMERCIALIZAÇÃO ou
ADEQUAÇÃO A UM PROPÓSITO ESPECÍFICO. Consulte a Licença GPLv3 para mais
detalhes.

Você deve ter recebido uma cópia da Licença GPLv3 junto com este programa. Se
não, veja <https://www.gnu.org/licenses/>.`
