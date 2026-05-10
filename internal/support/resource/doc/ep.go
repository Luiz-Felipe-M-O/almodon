package doc

import (
	"go/doc/comment"
	"html/template"
	"slices"
	"strings"
)

type EndPoint struct {
	Method, Path string
	RouteID      string
	Body         template.HTML
}

func NewEndPoint(text string) (EndPoint, bool) {
	var p comment.Parser
	doc := p.Parse(text).Content

	method, path, index := find_route(doc)
	if index < 0 {
		return EndPoint{}, false
	}

	doc = slices.Delete(doc, index, index+1)

	var b strings.Builder
	if err := parse_content(&b, doc); err != nil {
		return EndPoint{}, false
	}

	return EndPoint{
		Method:  method,
		Path:    path,
		RouteID: route_id(method, path),
		Body:    template.HTML(b.String()),
	}, true
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
