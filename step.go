package kif

import (
	"fmt"
	"strconv"
	"unicode"

	"golang.org/x/text/runes"
)

var (
	spaces = runes.Predicate(func(r rune) bool {
		switch r {
		case ' ', '　', '\t':
			return true
		default:
			return false
		}
	})
)

var (
	EOS         = fmt.Errorf("EOS")
	ErrMismatch = fmt.Errorf("mismatch")
)

type stepParser struct {
	line []rune
	curr int
}

func newStepParser(line string) *stepParser {
	p := &stepParser{
		line: []rune(line),
	}
	p.reset()
	return p
}

func (p *stepParser) reset() {
	p.curr = 0
}

func (p *stepParser) next() (rune, error) {
	if p.curr >= len(p.line) {
		return 0, EOS
	}

	r := p.line[p.curr]
	p.curr++

	return r, nil
}

func (p *stepParser) unread() {
	if p.curr == 0 {
		return
	}

	p.curr--
}

func (p *stepParser) skip(*Step) error {
	for {
		r, err := p.next()
		if err == EOS {
			return nil
		} else if err != nil {
			return err
		}

		if !spaces.Contains(r) {
			p.unread()
			return nil
		}
	}
}

func (p *stepParser) readRune(o rune) error {
	r, err := p.next()
	if err != nil {
		return err
	}

	if r != o {
		p.unread()
		return ErrMismatch
	}
	return nil
}

func (p *stepParser) readRunes(rs []rune) (int, error) {
	for i, r := range rs {
		if err := p.readRune(r); err == nil {
			return i, nil
		} else if err != ErrMismatch {
			return -1, err
		}
	}

	return -1, ErrMismatch
}

func (p *stepParser) readString(s string) error {
	for idx, r := range []rune(s) {
		if err := p.readRune(r); err == ErrMismatch {
			for i := 0; i < idx; i++ {
				p.unread()
			}
			return ErrMismatch
		} else if err != nil {
			return err
		}
	}

	return nil
}

func (p *stepParser) readStrings(ss []string) (int, error) {
	for i, s := range ss {
		if err := p.readString(s); err == nil {
			return i, nil
		} else if err != ErrMismatch {
			return -1, err
		}
	}

	return -1, ErrMismatch
}

func (p *stepParser) readInt() (int, error) {
	var rs []rune
	for {
		r, err := p.next()
		if err == EOS {
			break
		} else if err != nil {
			return 0, err
		}

		if !unicode.IsNumber(r) {
			p.unread()
			break
		}

		rs = append(rs, r)
	}

	if len(rs) == 0 {
		return 0, ErrMismatch
	}

	i, err := strconv.ParseInt(string(rs), 10, 32)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func (p *stepParser) readSeq(step *Step) error {
	i, err := p.readInt()
	if err != nil {
		return err
	}
	step.Seq = int32(i)
	return nil
}

func (p *stepParser) readPhase(*Step) error {
	_, err := p.readRunes([]rune("▲△"))
	if err == ErrMismatch {
		// skip
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

func (p *stepParser) readDst(step *Step) error {
	if err := p.readString("同　"); err == nil {
		return nil
	} else if err != ErrMismatch {
		return err
	}

	xidx, err := p.readRunes(xstr)
	if err != nil {
		return err
	}

	yidx, err := p.readRunes(ystr)
	if err != nil {
		return err
	}

	step.Dst = &Pos{X: int32(xidx), Y: int32(yidx)}
	return nil
}

func (p *stepParser) readPiece(step *Step) error {
	pi, err := p.readStrings(pieceNames)
	if err != nil {
		return err
	}

	var ret Piece_Id
	switch pi {
	case 15:
		ret = Piece_GYOKU
	case 16:
		ret = Piece_RYU
	case 17:
		ret = Piece_NARI_GIN
	case 18:
		ret = Piece_NARI_KEI
	case 19:
		ret = Piece_NARI_KYOU
	default:
		ret = Piece_Id(pi)
	}

	step.Piece = ret
	return nil
}

func (p *stepParser) readModifier(step *Step) error {
	rs := []rune("打成")
	i, err := p.readRunes(rs)
	if err == ErrMismatch {
		// null
		return nil
	} else if err != nil {
		return err
	}

	var ret Modifier_Id
	switch rs[i] {
	case '打':
		ret = Modifier_PUTTED
	case '成':
		ret = Modifier_PROMOTE
	default:
		return fmt.Errorf("unknown modifier")
	}
	step.Modifier = ret

	return nil
}

func (p *stepParser) readSrc(step *Step) error {
	if err := p.readRune('('); err == ErrMismatch {
		// putted
		return nil
	} else if err != nil {
		return err
	}

	var num = []rune(" 123456789")

	xi, err := p.readRunes(num)
	if err != nil {
		return err
	}

	yi, err := p.readRunes(num)
	if err != nil {
		return err
	}

	if err := p.readRune(')'); err != nil {
		return err
	}

	step.Src = &Pos{
		X: int32(xi),
		Y: int32(yi),
	}
	return nil
}

func (p *stepParser) readMove(step *Step) error {
	movei, err := p.readStrings(finStats)
	if err == nil {
		step.FinishedStatus = FinishedStatus_Id(movei)
		return nil
	} else if err != ErrMismatch {
		return err
	}

	for _, f := range []func(*Step) error{
		func(step *Step) error {
			err := p.readDst(step)
			if err != nil {
				return err
			}

			return nil
		},
		p.readPiece,
		p.readModifier,
		p.readSrc,
	} {
		if err := f(step); err != nil {
			return err
		}
	}

	return nil
}

func (p *stepParser) readTimestamp(step *Step) error {
	if err := p.readRune('('); err != nil {
		return err
	}

	if err := p.skip(step); err != nil {
		return err
	}

	var thinking int32
	tm, err := p.readInt()
	if err != nil {
		return err
	}
	thinking = int32(tm * 60)

	if err := p.readRune(':'); err != nil {
		return err
	}

	ts, err := p.readInt()
	if err != nil {
		return err
	}
	thinking += int32(ts)

	if err := p.readRune('/'); err != nil {
		return err
	}

	var elapsed int32
	eh, err := p.readInt()
	if err != nil {
		return err
	}
	elapsed += int32(eh * 60 * 60)

	if err := p.readRune(':'); err != nil {
		return err
	}

	em, err := p.readInt()
	if err != nil {
		return err
	}
	elapsed += int32(em * 60)

	if err := p.readRune(':'); err != nil {
		return err
	}

	es, err := p.readInt()
	if err != nil {
		return err
	}
	elapsed += int32(es)

	if err := p.readRune(')'); err != nil {
		return err
	}

	step.ThinkingSec = thinking
	step.ElapsedSec = elapsed
	return nil
}

func parseStep(in string) (*Step, error) {
	p := &stepParser{
		line: []rune(in),
	}

	step := &Step{}
	var prevDst *Pos
	for _, f := range []func(*Step) error{
		p.skip,
		p.readSeq,
		p.skip,
		p.readPhase,
		func(step *Step) error {
			err := p.readMove(step)
			if err != nil {
				return err
			}

			if step.Dst == nil {
				step.Dst = prevDst
			}

			prevDst = step.Dst
			return nil
		},
		p.skip,
		p.readTimestamp,
	} {
		if err := f(step); err != nil {
			return nil, err
		}
	}

	return step, nil
}
