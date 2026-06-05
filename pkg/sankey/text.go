// Copyright 2026 Alan Lima.
// The provided current package, sankey, is licensed separately from
// the main software, Almodon. The package sankey is licensed under
// the MIT License. See the notice it the LICENSE at this directory.

package sankey

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"math"
)

var ErrIncomplete = errors.New("incomplete diagram")

func (d Diagram) MarshalSVG() ([]byte, error) {
	if d.err != nil {
		return nil, fmt.Errorf("malformed diagram: %w", d.err)
	}

	if len(d.cols) <= 1 {
		return nil, ErrIncomplete
	}
	for _, col := range d.cols {
		if len(col.Accs) == 0 {
			return nil, ErrIncomplete
		}
	}

	gap, height := d._Height()
	ascale := float64(d.Height) / height
	gap *= ascale

	var b bytes.Buffer

	b.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	fmt.Fprintf(&b, `<svg viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		d.Bar*len(d.cols)+d.FlowWidth*(len(d.cols)-1),
		d.Height,
	)

	d._Flows(&b, gap, ascale)
	d._Bars(&b, gap, ascale)

	b.WriteString(`</svg>`)
	return b.Bytes(), nil
}

func (d Diagram) _Height() (float64, float64) {
	sum := math.Inf(-1)

	for _, col := range d.cols {
		var lsum float64
		for _, src := range col.Accs {
			lsum += src.Amount
		}
		sum = max(sum, lsum)
	}

	gap := d.Gap * sum
	sum = math.Inf(-1)

	for _, col := range d.cols {
		lsum := gap * float64(len(col.Accs)-1)
		for _, src := range col.Accs {
			lsum += src.Amount
		}
		sum = max(sum, lsum)
	}

	return gap, sum
}

func (d Diagram) _Bars(b *bytes.Buffer, gap, ascale float64) {
	var x int
	for _, col := range d.cols {
		var y float64
		for _, src := range col.Accs {
			amount := ascale * src.Amount
			fmt.Fprintf(b, `<rect x="%d" y="%g" width="%d" height="%g" rx="1" fill="#%08x" />`,
				x, y, d.Bar, amount, hex(src.Color),
			)

			y += gap + amount
		}

		x += d.Bar + d.FlowWidth
	}
}

func (d Diagram) _Flows(b *bytes.Buffer, gap, ascale float64) {
	for i := range len(d.cols) {
		col := &d.cols[i]
		var sum float64

		for j := range len(col.Accs) {
			src := &col.Accs[j]

			src.mem = sum
			sum += src.Amount
		}
	}

	x0 := d.Bar
	x3 := d.Bar + d.FlowWidth

	for i := range len(d.cols) - 1 {
		pcol := &d.cols[i]
		ccol := &d.cols[i+1]

		var y3 float64

		for _, dst := range ccol.Accs {
			for _, flow := range dst.Inflows {
				src := &pcol.Accs[flow.Source]

				amount := ascale * flow.Amount
				y0 := gap*float64(flow.Source) + ascale*(src.mem)
				src.mem += flow.Amount

				d._Path(b, x0, x3, y0, y3, flow.Color, amount)
				y3 += amount
			}

			y3 += gap
		}

		x0 += d.Bar + d.FlowWidth
		x3 += d.Bar + d.FlowWidth
	}
}

func (d Diagram) _Path(b *bytes.Buffer, x0, x3 int, y0, y3 float64, c color.RGBA, amount float64) {
	alpha := float64(x0+x3) / 2
	y0 += amount / 2
	y3 += amount / 2

	fmt.Fprintf(b, `<path d="M %d,%g C %g,%g %g,%g %d,%g" fill="none" stroke="#%08x" stroke-width="%g" />`,
		x0, y0, alpha, y0, alpha, y3, x3, y3,
		hex(c), amount,
	)
}

func hex(clr color.RGBA) uint32 {
	return uint32(clr.R)<<24 | uint32(clr.G)<<16 | uint32(clr.B)<<8 | uint32(clr.A)
}
