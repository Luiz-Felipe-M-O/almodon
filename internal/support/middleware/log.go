package middleware

import (
	"io"
	"log"

	"github.com/alan-b-lima/ansi-escape-sequences"
)

const (
	_Info  = `INFO`
	_Warn  = `WARN`
	_Error = `ERROR`
)

type Logger struct {
	il   log.Logger
	ansi bool

	info  string
	warn  string
	error string
}

func NewLogger(w io.Writer, name string) *Logger {
	l := new(Logger)

	l.ansi = enableAnsi(w)

	l.il.SetOutput(w)
	l.il.SetFlags(log.Ldate | log.Ltime)
	l.il.SetPrefix(name + "> ")

	l.info = _Info
	l.warn = _Warn
	l.error = _Error

	return l
}

func (l *Logger) Error(v ...any)                 { print(&l.il, l.error, v...) }
func (l *Logger) Errorf(format string, v ...any) { printf(&l.il, l.error, format, v...) }

func (l *Logger) Warn(v ...any)                 { print(&l.il, l.warn, v...) }
func (l *Logger) Warnf(format string, v ...any) { printf(&l.il, l.warn, format, v...) }

func (l *Logger) Info(v ...any)                 { print(&l.il, l.info, v...) }
func (l *Logger) Infof(format string, v ...any) { printf(&l.il, l.info, format, v...) }

func print(log *log.Logger, level string, v ...any) {
	log.Print(append([]any{level + " "}, v...)...)
}

func printf(log *log.Logger, level string, format string, v ...any) {
	log.Printf(level+" "+format, v...)
}

func enableAnsi(w io.Writer) bool {
	f, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return false
	}

	err := ansi.EnableVirtualTerminal(f.Fd())
	return err == nil
}
