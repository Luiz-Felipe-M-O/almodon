package server

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	listener net.Listener
	server   http.Server

	base   context.Context
	cancel context.CancelFunc

	signals chan os.Signal
	active  chan struct{}
}

func New(addr string, handler http.Handler, logger *slog.Logger) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	var log *log.Logger
	if logger != nil {
		log = slog.NewLogLogger(logger.Handler(), slog.LevelError)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	base, cancel := context.WithCancel(context.Background())

	s := Server{
		listener: ln,
		server: http.Server{
			Handler: handler,
			BaseContext: func(net.Listener) context.Context {
				return base
			},
			ErrorLog: log,
		},

		base:   base,
		cancel: cancel,

		signals: signals,
		active:  make(chan struct{}),
	}

	return &s, nil
}

// Addr returns the address the server is listening on.
func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

// Serve starts the server, this function only stops on error or server
// termination. On error, it returns the error direclty, on termination, it
// returns nil.
func (s *Server) Serve() error {
	if err := s.server.Serve(s.listener); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}

		return err
	}

	return nil
}

func (s *Server) Done() <-chan struct{} {
	return s.active
}

// Signal returns a channel that will be closed on the first signal received.
func (s *Server) Signal() <-chan os.Signal {
	return s.signals
}

// CancelOngoing cancels the base context for all requests.
func (s *Server) CancelOngoing() {
	s.cancel()
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	s.release()
	return nil
}

func (s *Server) ForceShutdown() error {
	defer s.release()

	return s.server.Close()
}

func (s *Server) release() error {
	defer close(s.active)

	signal.Stop(s.signals)
	s.cancel()

	return s.listener.Close()
}
