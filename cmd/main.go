package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/alan-b-lima/almodon/internal/almodon"
	"github.com/alan-b-lima/almodon/internal/server"
	"github.com/alan-b-lima/almodon/internal/support/resource"
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

	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	if err := Main(addr, logger); err != nil {
		logger.Error(err.Error())
	}
}

func Main(addr string, log *slog.Logger) (reterr error) {
	api, err := almodon.New()
	if err != nil {
		return err
	}
	defer func() {
		if err := api.Close(); err != nil {
			reterr = err
		}
	}()

	var mux http.ServeMux
	handler := Logger(log, &mux)

	server, err := server.New(addr, handler, log)
	if err != nil {
		return err
	}

	mux.Handle("/", api)
	mux.HandleFunc("/terminate", func(w http.ResponseWriter, r *http.Request) {
		go Shutdown(log, server)
		w.WriteHeader(http.StatusNoContent)
	})

	go SignalShutdown(log, server)

	log.Info("server listening at http://" + strings.Replace(server.Addr().String(), "[::]", "localhost", 1))
	if err := server.Serve(); err != nil {
		return err
	}

	<-server.Done()

	log.Info("server powering off...")
	return nil
}

func SignalShutdown(log *slog.Logger, server *server.Server) {
	<-server.Signal()
	fmt.Print("\r")

	if err := Shutdown(log, server); err != nil {
		log.Error(err.Error())
	}
}

func Shutdown(log *slog.Logger, server *server.Server) error {
	log.Info("starting server shutdown...")

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

	log.Warn("starting forceful server shutdown...")
	return server.ForceShutdown()
}

func Logger(log *slog.Logger, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := resource.NewResponseWriter(w)
		start := time.Now()

		h.ServeHTTP(rw, r)

		elapsed := time.Since(start)

		level := slog.LevelInfo
		if rw.StatusCode()/100 == 5 {
			level = slog.LevelError
		}

		log.LogAttrs(r.Context(), level, "completed request",
			slog.Int("status", rw.StatusCode()),
			slog.String("remote", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.String()),
			slog.Duration("elapsed", elapsed),
			slog.Any("error", rw.Error()),
		)
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
