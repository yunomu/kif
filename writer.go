package kif

import (
	"fmt"
	"io"

	"github.com/yunomu/kif/ptypes"
	"golang.org/x/text/transform"
)

type Writer struct {
	newline string

	encodingTransformer func(io.Writer) io.Writer
}

type WriterOption func(*Writer)

func SetNewline(newline string) WriterOption {
	return func(w *Writer) {
		w.newline = newline
	}
}

var sjisWriter = func(wr io.Writer) io.Writer {
	return transform.NewWriter(wr, sjisDecoder)
}

var WriteEncodingSJIS WriterOption = func(w *Writer) {
	w.encodingTransformer = sjisWriter
}

var WriteEncodingUTF8 WriterOption = func(w *Writer) {
	w.encodingTransformer = func(wr io.Writer) io.Writer {
		return wr
	}
}

func NewWriter(ops ...WriterOption) *Writer {
	w := &Writer{
		newline:             "\n",
		encodingTransformer: sjisWriter,
	}

	for _, f := range ops {
		f(w)
	}

	return w
}

func stepToLine(step *ptypes.Step) string {
	return fmt.Sprintf(
		"%4d %-12s (%s/%s)",
		step.Seq,
		PrintMove(step),
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

func (w *Writer) Write(out io.Writer, kif *ptypes.Kif) error {
	Normalize(kif)
	p := &linePrinter{
		newline: w.newline,
		w:       w.encodingTransformer(out),
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
