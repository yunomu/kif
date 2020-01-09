package parser

import (
	"testing"

	"github.com/yunomu/kif/ptypes"
)

func TestStepParser_skip(t *testing.T) {
	var p *stepParser

	p = newStepParser("   ")
	if err := p.skip(nil); err != nil {
		t.Fatalf("skip EOS: %v", err)
	}
	if p.curr != len(p.line) {
		t.Fatalf("line is not empty: remain=%v", len(p.line)-p.curr)
	}
	if r, err := p.next(); err != EOS {
		t.Fatalf("empty but not occur error: err=`%v`, rune=%v", err, r)
	}

	p = newStepParser(" 12  ")
	if err := p.skip(nil); err != nil {
		t.Fatalf("skip before `12`: %v", err)
	}
	if p.curr != 1 {
		t.Fatalf("p.curr=%v expected=1", p.curr)
	}
	r, err := p.next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r != '1' {
		t.Fatalf("next is not '1' actual='%v'", r)
	}
}

func TestStepParser_readInt(t *testing.T) {
	var p *stepParser

	p = newStepParser("123")
	i, err := p.readInt()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if i != 123 {
		t.Fatalf("expected=123 actual=%v", i)
	}

	p = newStepParser("45abc")
	i2, err := p.readInt()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if i2 != 45 {
		t.Fatalf("expected=45 actual=%v", i)
	}

	p = newStepParser("abc")
	if _, err := p.readInt(); err != ErrMismatch {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepParser_readPiece(t *testing.T) {
	var p *stepParser
	var s *ptypes.Step

	p = newStepParser("歩")
	s = &ptypes.Step{}
	if err := p.readPiece(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Piece != ptypes.Piece_FU {
		t.Fatalf("expected=FU actual=%v", s.Piece)
	}

	p = newStepParser("成銀")
	s = &ptypes.Step{}
	if err := p.readPiece(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Piece != ptypes.Piece_NARI_GIN {
		t.Fatalf("expected=NARI_GIN actual=%v", s.Piece)
	}

	p = newStepParser("全")
	s = &ptypes.Step{}
	if err := p.readPiece(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Piece != ptypes.Piece_NARI_GIN {
		t.Fatalf("expected=NARI_GIN actual=%v", s.Piece)
	}
}

func TestStepParser_readPhase(t *testing.T) {
	var p *stepParser
	s := &ptypes.Step{}

	p = newStepParser("▲")
	if err := p.readPhase(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p = newStepParser("△")
	if err := p.readPhase(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p = &stepParser{
		line: []rune("三"),
	}
	if err := p.readPhase(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStepParser_readDst(t *testing.T) {
	var p *stepParser
	s := &ptypes.Step{}

	p = &stepParser{
		line: []rune("７六"),
	}
	if err := p.readDst(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Dst.X != 7 || s.Dst.Y != 6 {
		t.Fatalf("expected dst (7, 6) but <%v>", s.Dst)
	}

	p = &stepParser{
		line: []rune("同　"),
	}
	s = &ptypes.Step{}
	if err := p.readDst(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Dst != nil {
		t.Fatalf("pos is not nil: <%v>", s.Dst)
	}
}

func TestStepParser_readSrc(t *testing.T) {
	var p *stepParser
	var s *ptypes.Step

	p = &stepParser{
		line: []rune("(83)"),
	}
	s = &ptypes.Step{}
	if err := p.readSrc(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Src.X != 8 || s.Src.Y != 3 {
		t.Fatalf("expected src (8, 3) but <%v>", s.Src)
	}

	p = &stepParser{
		line: []rune("　"),
	}
	s = &ptypes.Step{}
	if err := p.readSrc(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Src != nil {
		t.Fatalf("pos is not nil: <%v>", s.Src)
	}
}

func TestStepParser_readTimestamp(t *testing.T) {
	var p *stepParser

	p = newStepParser("( 1:23/01:23:45)")
	s := &ptypes.Step{}
	if err := p.readTimestamp(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.ThinkingSec != 60+23 {
		t.Fatalf("unexpected thinking sec: %v", s.ThinkingSec)
	}
	if s.ElapsedSec != 3600+23*60+45 {
		t.Fatalf("unexpected thinking sec: %v", s.ElapsedSec)
	}
}

func TestStepParser_readString(t *testing.T) {
	var p *stepParser

	p = newStepParser("すもももももももものうち")
	if err := p.readString("すもも"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := p.readString("ももい"); err != ErrMismatch {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.curr != 3 {
		t.Fatalf("curr expected=3 actual=%v", p.curr)
	}
}

func TestStepParser_readMove(t *testing.T) {
	var p *stepParser
	var s *ptypes.Step

	p = newStepParser("７六歩(77)")
	s = &ptypes.Step{}
	if err := p.readMove(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.FinishedStatus != ptypes.FinishedStatus_NOT_FINISHED {
		t.Fatalf("unexpected sp_move: %v", s.FinishedStatus)
	}

	p = newStepParser("投了")
	s = &ptypes.Step{}
	if err := p.readMove(s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.FinishedStatus != ptypes.FinishedStatus_SURRENDER {
		t.Fatalf("unexpected sp_move: %v", s.FinishedStatus)
	}
}

func TestParseStep(t *testing.T) {
	step, err := ParseStep("  73 ７八金打   ( 0:00/00:00:00)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if step.Seq != 73 {
		t.Fatalf("unexpected seq: %v", step.Seq)
	}

	step2, err := ParseStep("  1 ７八成銀(69)   ( 0:10/00:00:00)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if step2.Seq != 1 {
		t.Fatalf("unexpected seq: %v", step2.Seq)
	}
}
