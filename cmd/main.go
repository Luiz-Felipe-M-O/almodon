package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/alan-b-lima/almodon/internal/api/v1"
	"github.com/alan-b-lima/almodon/internal/support/middleware"
)

var StdOut = os.Stdout

func main() {
	log := middleware.NewLogger(StdOut, "")
	style := middleware.Styles()

	ln, err := net.Listen("tcp", ":4545")
	if err != nil {
		log.Error(err)
		return
	}
	defer ln.Close()

	api, err := api.New()
	if err != nil {
		log.Error(err)
		return
	}
	defer func() {
		if err := api.Close(); err != nil {
			log.Error(err)
		}
	}()

	srv := http.Server{Handler: middleware.LogTraffic(log, style, MakeMux(api))}
	done := EnableGracefulShutdown(func() {
		log.Info("Shutting server down...")
		srv.Shutdown(context.Background())
	})

	url := "http://" + strings.Replace(ln.Addr().String(), "[::]", "localhost", 1)
	log.Infof("Server listening at %s\n", style.HyperLink(url))

	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Error(err)
	}

	<-done
}

func MakeMux(api *api.Handler) *http.ServeMux {
	mux := new(http.ServeMux)

	fs := http.FileServer(http.Dir("../ui/web/dist/"))
	f := ServeFile("../ui/web/dist/index.html")

	mux.Handle("/", fs)
	mux.Handle("/{path}", f)
	mux.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("../ui/web/src/"))))
	mux.Handle("/api/", api)
	mux.HandleFunc("/terminate/{timeout}", Terminate)

	return mux
}

var Signals chan<- os.Signal

func EnableGracefulShutdown(fn func()) <-chan struct{} {
	signals := make(chan os.Signal, 1)
	Signals = signals

	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	done := make(chan struct{}, 1)

	go func() {
		<-signals
		fn()
		done <- struct{}{}
	}()

	return done
}

type ServeFile string

func (s ServeFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, string(s))
}

func Terminate(w http.ResponseWriter, r *http.Request) {
	ms, err := strconv.Atoi(r.PathValue("timeout"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	go func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		Signals <- syscall.SIGTERM
	}()

	w.WriteHeader(http.StatusNoContent)
}
