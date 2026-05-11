package doc

import (
	"go/doc/comment"
	"html/template"
	"io"
	"strings"
)

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

func execute(w io.Writer, text string, data any) error {
	tmpl, err := template.New("").Parse(text)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}
