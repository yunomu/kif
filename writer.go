package kif

import (
	"fmt"
	"io"
)

type Writer struct {
	newline string
}

type WriterOption func(*Writer)

func SetNewline(newline string) WriterOption {
	return func(w *Writer) {
		w.newline = newline
	}
}

func NewWriter(ops ...WriterOption) *Writer {
	w := &Writer{
		newline: "\n",
	}

	for _, f := range ops {
		f(w)
	}

	return w
}

func stepToLine(step *Step) string {
	return fmt.Sprintf(
		"%4d %-12s (%s/%s)",
		step.Seq,
		step.PrintMove(),
		PrintThinking(step.ThinkingSec),
		PrintElapsed(step.ElapsedSec),
	)
}

type linePrinter struct {
	newline string
	w       io.Writer
}

func (p *linePrinter) Print(str string) error {
	if _, err := fmt.Fprint(p.w, str); err != nil {
		return err
	}
	if _, err := fmt.Fprint(p.w, p.newline); err != nil {
		return err
	}
	return nil
}

func (w *Writer) Write(out io.Writer, kif *Kif) error {
	kif.Normalize()
	p := &linePrinter{
		newline: w.newline,
		w:       out,
	}

	for _, h := range kif.Headers {
		if err := p.Print(fmt.Sprintf("%s：%s", h.Name, h.Value)); err != nil {
			return err
		}
	}

	if err := p.Print("手数----指手---------消費時間--"); err != nil {
		return err
	}

	for _, step := range kif.Steps {
		if err := p.Print(stepToLine(step)); err != nil {
			return err
		}
		for _, note := range step.Notes {
			if err := p.Print("*" + note); err != nil {
				return err
			}
		}
	}

	return nil
}
