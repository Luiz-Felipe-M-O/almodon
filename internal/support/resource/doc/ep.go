package doc

import (
	"go/doc/comment"
	"html/template"
	"io"
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

func parse_content(b *strings.Builder, doc []comment.Block) error {
	for _, block := range doc {
		switch block := block.(type) {
		case *comment.Code:
			if err := execute(b, `<pre><code>{{ . }}</code></pre>`, block.Text); err != nil {
				return err
			}

		case *comment.Heading:
			b.WriteString(`<h3>`)
			if err := parse_text(b, block.Text); err != nil {
				return err
			}
			b.WriteString(`</h3>`)

		case *comment.Paragraph:
			b.WriteString(`<p>`)
			if err := parse_text(b, block.Text); err != nil {
				return err
			}
			b.WriteString(`</p>`)

		case *comment.List:
			if err := parse_list(b, block); err != nil {
				return err
			}
		}
	}

	return nil
}

func parse_text(b *strings.Builder, text []comment.Text) error {
	for _, text := range text {
		switch text := text.(type) {
		case comment.Plain:
			b.WriteString(string(text))

		case comment.Italic:
			if err := execute(b, `<em>{{ . }}</em>`, string(text)); err != nil {
				return err
			}

		case *comment.Link:
			if err := execute(b, `<a href="{{ . }}">`, text.URL); err != nil {
				return err
			}
			if err := parse_text(b, text.Text); err != nil {
				return err
			}
			b.WriteString(`</a>`)

		case *comment.DocLink:
			if err := parse_text(b, text.Text); err != nil {
				return err
			}
		}
	}

	return nil
}

func parse_list(b *strings.Builder, list *comment.List) error {
	if len(list.Items) == 0 {
		return nil
	}

	ordered := list.Items[0].Number != ""

	if ordered {
		b.WriteString(`<ol>`)
	} else {
		b.WriteString(`<ul>`)
	}

	for _, item := range list.Items {
		b.WriteString(`<li>`)
		if err := parse_content(b, item.Content); err != nil {
			return err
		}
		b.WriteString(`</li>`)
	}

	if ordered {
		b.WriteString(`</ol>`)
	} else {
		b.WriteString(`</ul>`)
	}

	return nil
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

func execute(w io.Writer, text string, data any) error {
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}
