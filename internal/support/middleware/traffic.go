package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/alan-b-lima/ansi-escape-sequences"
)

func LogTraffic(log *Logger, s Style, handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := NewW(w)

		handler.ServeHTTP(rw, r)

		status := rw.StatusCode()
		pen := s.StatusCodePen(status)

		var barr [32]byte
		b := bytes.NewBuffer(barr[:0])
		pen.Writer = b

		fmt.Fprintf(&pen, " %03d ", status)
		if 500 <= status && status < 599 {
			log.Errorf("%s %s %s %s\n", b.String(), r.RemoteAddr, r.Method, r.URL)
		} else {
			log.Infof("%s %s %s %s\n", b.String(), r.RemoteAddr, r.Method, r.URL)
		}
	}
}

type Style struct {
	hyperlink func(string) string

	pens map[int]ansi.Pen
	no   ansi.Pen

	enabled bool
}

func Styles() (s Style) {
	s.enabled = ansi.EnableVirtualTerminal(os.Stdout.Fd()) == nil
	s.no.SetStyle(false)

	if s.enabled {
		var Success ansi.Pen
		var Redirect ansi.Pen
		var ClientError ansi.Pen
		var ServerError ansi.Pen

		Success.BGColor(ansi.RGBFromHex(0x0ed145))
		Success.FGColor(ansi.RGBFromHex(0xffffff))

		Redirect.BGColor(ansi.RGBFromHex(0x4b53cc))
		Redirect.FGColor(ansi.RGBFromHex(0xffffff))

		ClientError.BGColor(ansi.RGBFromHex(0xea1d1d))
		ClientError.FGColor(ansi.RGBFromHex(0xffffff))

		ServerError.BGColor(ansi.RGBFromHex(0x88001b))
		ServerError.FGColor(ansi.RGBFromHex(0xffffff))

		s.pens = map[int]ansi.Pen{
			2: Success,
			3: Redirect,
			4: ClientError,
			5: ServerError,
		}
		s.hyperlink = hyperlink

		return s
	}

	s.pens = map[int]ansi.Pen{}
	return s
}

func (s *Style) HyperLink(url string) string {
	if s.hyperlink == nil {
		return url
	}

	return s.hyperlink(url)
}

func (s *Style) StatusCodePen(status int) ansi.Pen {
	pen, in := s.pens[status/100]
	if !in {
		return s.no
	}

	return pen
}

func hyperlink(link string) string {
	var pen ansi.Pen
	pen.FGColor(ansi.RGBFromHex(0x4e8597))

	return pen.Sprint(ansi.HyperLinkP(link))
}
