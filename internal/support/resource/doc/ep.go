package doc

import (
	"errors"
	"go/doc/comment"
	"html/template"
	"strings"

	"github.com/alan-b-lima/almodon/ui/web"
)

type EndPoint struct {
	Method, Path string
	RouteID      string
	Body         template.HTML
}

var ErrNotRoute = errors.New("route not found")

func NewEndPoint(text string) (EndPoint, error) {
	var p comment.Parser
	doc := p.Parse(text).Content

	method, path, index := find_route(doc)
	if index < 0 {
		return EndPoint{}, ErrNotRoute
	}

	html, err := web.GoComment(text)
	if err != nil {
		return EndPoint{}, err
	}

	return EndPoint{
		Method:  method,
		Path:    path,
		RouteID: route_id(method, path),
		Body:    html,
	}, nil
}

func route_id(method, path string) string {
	var b strings.Builder

	b.WriteString(strings.ToLower(method))
	for _, rune := range path {
		switch rune {
		case ':', '{', '}':
			// ignore
		case '/':
			b.WriteByte('-')
		default:
			b.WriteRune(rune)
		}
	}

	id := b.String()

	if id[len(id)-1] == '-' {
		return id[:len(id)-1]
	}
	return id
}

func find_route(doc []comment.Block) (method, path string, index int) {
	for i, block := range doc {
		block, ok := block.(*comment.Code)
		if !ok {
			continue
		}

		line := block.Text
		if index := strings.IndexRune(block.Text, '\n'); index >= 0 {
			line = line[:index]
		}

		line = strings.TrimSpace(line)
		if method, path, ok := strings.Cut(line, " "); ok {
			return method, path, i
		}
	}

	return "", "", -1
}
