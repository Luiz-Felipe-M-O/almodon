// Copyright 2026 Alan Lima.
// The provided current package, sankey, is licensed separately from
// the main software, Almodon. The package sankey is licensed under
// the MIT License. See the notice it the LICENSE at this directory.

package sankey

import (
	"errors"
	"image/color"
)

type Diagram struct {
	Precision float64

	Height    int
	FlowWidth int
	Bar       int
	Gap       float64

	cols []column
	err  error
}

type column struct {
	Accs []src
}

type src struct {
	Title   string
	Inflows []flow
	Amount  float64
	Color   color.RGBA

	mem float64
}

type flow struct {
	Source int
	Amount float64

	Color color.RGBA
}

func New() Diagram {
	return Diagram{
		Precision: 1e-5,
		Height:    512,
		FlowWidth: 512,
		Bar:       4,
		Gap:       0.05,
	}
}

var (
	ErrPrevCol  = errors.New("sankey: absent or incoherent previous column")
	ErrCurrCol  = errors.New("sankey: absent or incoherent current column")
	ErrNegative = errors.New("sankey: negative amount")
)

func (d *Diagram) Source(title string, amount float64, c color.Color) {
	d.checkNonNeg(amount)
	if d.err != nil {
		return
	}

	ccol := d.currCol()
	if ccol == nil {
		return
	}

	ccol.Accs = append(ccol.Accs, src{
		Title:  title,
		Amount: amount,
		Color:  color.RGBAModel.Convert(c).(color.RGBA),
		mem:    amount,
	})
}

func (d *Diagram) Dest(title string, c color.Color) {
	if d.err != nil {
		return
	}

	ccol := d.currCol()
	if ccol == nil {
		return
	}

	ccol.Accs = append(ccol.Accs, src{
		Title:   title,
		Inflows: []flow{},
		Color:   color.RGBAModel.Convert(c).(color.RGBA),
	})
}

func (d *Diagram) Flow(from int, amount float64, c color.Color) {
	d.checkNonNeg(amount)
	if d.err != nil {
		return
	}

	dst := d.currDst()
	if dst == nil {
		return
	}

	src := d.prevSrc(from)
	if src == nil {
		return
	}

	dst.Inflows = append(dst.Inflows, flow{
		Source: from,
		Amount: amount,
		Color:  color.RGBAModel.Convert(c).(color.RGBA),
	})

	dst.Amount += amount
	dst.mem += amount

	src.mem -= amount
	d.checkNonNeg(src.mem)
}

func (d *Diagram) Col() {
	d.cols = append(d.cols, column{})
}

func (d *Diagram) Err() error {
	return d.err
}

func (d *Diagram) currDst() *src {
	ccol := d.currCol()
	if ccol == nil {
		return nil
	}

	if len(ccol.Accs) == 0 {
		d.setErr(ErrCurrCol)
		return nil
	}

	dst := &ccol.Accs[len(ccol.Accs)-1]
	if dst.Inflows == nil {
		d.setErr(ErrCurrCol)
		return nil
	}

	return dst
}

func (d *Diagram) prevSrc(i int) *src {
	pcol := d.prevCol()
	if pcol == nil {
		return nil
	}

	if i < 0 || len(pcol.Accs) <= i {
		d.setErr(ErrPrevCol)
		return nil
	}

	return &pcol.Accs[i]
}

func (d *Diagram) currCol() *column {
	if len(d.cols) == 0 {
		d.setErr(ErrCurrCol)
		return nil
	}

	return &d.cols[len(d.cols)-1]
}

func (d *Diagram) prevCol() *column {
	if len(d.cols) <= 1 {
		d.setErr(ErrPrevCol)
		return nil
	}

	return &d.cols[len(d.cols)-2]
}

func (d *Diagram) checkNonNeg(num float64) {
	if num < -d.Precision {
		d.setErr(ErrNegative)
		return
	}
}

func (d *Diagram) setErr(err error) {
	if d.err != nil || err == nil {
		return
	}

	*d = Diagram{err: err}
}
