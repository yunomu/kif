package kif

import (
	"bufio"
	"io"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/pkg/errors"

	"github.com/yunomu/kif/ptypes"
)

type lineReader struct {
	r *bufio.Reader

	unreadLine string
}

func newLineReader(r *bufio.Reader) *lineReader {
	return &lineReader{
		r: r,
	}
}

func (r *lineReader) Read() (string, error) {
	if ret := r.unreadLine; ret != "" {
		r.unreadLine = ""
		return ret, nil
	}

	bs, _, err := r.r.ReadLine()
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func (r *lineReader) Unread(line string) error {
	if r.unreadLine != "" {
		return errors.New("already unreaded")
	}

	r.unreadLine = line
	return nil
}

func dropBOM(r *bufio.Reader) error {
	bs, err := r.Peek(3)
	if err != nil {
		return err
	}

	if bs[0] == 0xEF && bs[1] == 0xBB && bs[2] == 0xBF {
		r.Discard(3)
	}

	return nil
}

type options struct {
	transformReader func(io.Reader) io.Reader
}

type ParseOption func(*options)

var (
	sjisDecoder = japanese.ShiftJIS.NewDecoder()
	sjisReader  = func(r io.Reader) io.Reader {
		return transform.NewReader(r, sjisDecoder)
	}
)

var ParseEncodingSJIS ParseOption = func(ops *options) {
	ops.transformReader = sjisReader
}

var ParseEncodingUTF8 ParseOption = func(ops *options) {
	ops.transformReader = func(r io.Reader) io.Reader {
		return r
	}
}

func Parse(in io.Reader, ops ...ParseOption) (*ptypes.Kif, error) {
	options := &options{
		transformReader: sjisReader,
	}

	for _, f := range ops {
		f(options)
	}

	var count int
	br := bufio.NewReader(options.transformReader(in))

	if err := dropBOM(br); err != nil {
		return nil, err
	}
	r := newLineReader(br)

	ret := &ptypes.Kif{}

	// read header
	for {
		count++

		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		header := strings.SplitN(line, "：", 2)
		if len(header) != 2 {
			r.Unread(line)
			break
		}

		ret.Headers = append(ret.Headers, &ptypes.Header{
			Name:  header[0],
			Value: header[1],
		})
	}
	r.Read()

	var prevStep *ptypes.Step
	for {
		count++

		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		if line[0] == '*' {
			prevStep.Notes = append(prevStep.Notes, line[1:])
			continue
		}
		if prevStep.GetFinishedStatus() != ptypes.FinishedStatus_NOT_FINISHED {
			prevStep.Notes = append(prevStep.Notes, line)
			continue
		}

		step, err := parseStep(line)
		if err != nil {
			return nil, errors.Wrapf(err, "line=%v %v", count, line)
		}

		ret.Steps = append(ret.Steps, step)
		prevStep = step
	}

	return ret, nil
}
