package web

import (
	"fmt"
	"go/doc/comment"
	"html/template"
	"strings"
)

func GoComment(text string) (template.HTML, error) {
	var parser comment.Parser
	doc := parser.Parse(text)

	var printer printer
	return printer.html(doc), nil
}

type printer struct {
	comment.Printer
}

func (p *printer) html(doc *comment.Doc) template.HTML {
	var b strings.Builder
	p.content(&b, doc.Content)
	return template.HTML(b.String())
}

func (p *printer) content(b *strings.Builder, content []comment.Block) {
	for _, block := range content {
		p.block(b, block)
	}
}

func (p *printer) block(b *strings.Builder, block comment.Block) {
	switch block := block.(type) {
	default:
		p.escape(b, fmt.Sprintf("?%T", block))

	case *comment.Paragraph:
		b.WriteString(`<p>`)
		p.text(b, block.Text)
		b.WriteString(`</p>`)

	case *comment.Code:
		b.WriteString(`<pre><code>`)
		p.escape(b, block.Text)
		b.WriteString(`</code></pre>`)

	case *comment.Heading:
		h := '0' + p.headingLevel()

		b.WriteString(`<h`)
		b.WriteByte(h)
		if id := p.headingID(block); id != "" {
			b.WriteString(` id="`)
			p.escape(b, id)
			b.WriteString(`"`)
		}
		b.WriteString(`>`)

		p.text(b, block.Text)

		b.WriteString(`</h`)
		b.WriteByte(h)
		b.WriteString(`>`)

	case *comment.List:
		if len(block.Items) == 0 {
			break
		}

		kind := `ol>`
		if block.Items[0].Number == "" {
			kind = `ul>`
		}

		b.WriteString(`<`)
		b.WriteString(kind)

		for _, item := range block.Items {
			b.WriteString(`<li>`)
			p.content(b, item.Content)
			b.WriteString(`</li>`)
		}

		b.WriteString(`</`)
		b.WriteString(kind)
	}
}

func (p *printer) text(b *strings.Builder, text []comment.Text) {
	for _, text := range text {
		switch text := text.(type) {
		case comment.Plain:
			p.escape(b, string(text))

		case comment.Italic:
			b.WriteString(`<i>`)
			p.escape(b, string(text))
			b.WriteString(`</i>`)

		case *comment.Link:
			b.WriteString(`<a href="`)
			p.escape(b, text.URL)
			b.WriteString(`">`)
			p.text(b, text.Text)
			b.WriteString(`</a>`)

		case *comment.DocLink:
			p.text(b, text.Text)
		}
	}
}

func (p *printer) headingLevel() byte {
	if p.HeadingLevel <= 0 || 6 < p.HeadingLevel {
		return 3
	}
	return byte(p.HeadingLevel)
}

func (p *printer) headingID(h *comment.Heading) string {
	if p.HeadingID == nil {
		return h.DefaultID()
	}
	return p.HeadingID(h)
}

var escaper = strings.NewReplacer(
	`&`, `&amp;`,
	`'`, `&apos;`,
	`"`, `&quot;`,
	`<`, `&lt;`,
	`>`, `&gt;`,
)

func (p *printer) escape(b *strings.Builder, s string) {
	escaper.WriteString(b, s)
}
