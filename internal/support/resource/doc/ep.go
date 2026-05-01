package doc

import "strings"

type EndPoint struct {
	Method, Path string
	Header       string
	Body         []string
}

func NewEndPoint(doc string) (EndPoint, bool) {
	lines := split_lines(doc)
	if len(lines) == 0 {
		return EndPoint{}, false
	}

	method, path, index := find_route(lines[1:])
	if index < 0 {
		return EndPoint{}, false
	}

	ep := EndPoint{
		Method: method,
		Path:   path,
		Header: lines[0],
		Body:   make([]string, 0, len(lines)-1),
	}

	copy(ep.Body[:index], lines[1:index+1])
	copy(ep.Body[index:], lines[index+2:])

	return ep, true
}

func (ep *EndPoint) RouteID() string {
	var b strings.Builder

	b.WriteString(strings.ToLower(ep.Method))
	for _, rune := range ep.Path {
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

func split_lines(doc string) []string {
	lines := strings.Split(doc, "\n")
	final := lines[:0]

	var b strings.Builder

	for _, line := range lines {
		if len(line) == 0 {
			if b.Len() > 0 {
				final = append(final, b.String())
				b.Reset()
			}

			continue
		}

		if b.Len() > 0 {
			b.WriteRune(' ')
		}
		b.WriteString(line)
	}

	if b.Len() > 0 {
		final = append(final, b.String())
	}

	return final
}

func find_route(doc []string) (method, path string, index int) {
	for i, line := range doc {
		if len(line) > 0 && line[0] == '\t' {
			if method, path, ok := strings.Cut(line[1:], " "); ok {
				return method, path, i
			}
		}
	}

	return "", "", -1
}
