package kif

import (
	"testing"

	"time"

	"github.com/yunomu/kif/ptypes"
)

func TestPieceFronName(t *testing.T) {
	for i := 0; i <= 14; i++ {
		e := pieceNames[i]
		a := PrintPiece(PieceFromName(e))
		if e != a {
			t.Errorf("expected=%s actual=%s", e, a)
		}
	}

	if s := PrintPiece(PieceFromName("王")); s != "玉" {
		t.Errorf("expected=玉 actual=%s", s)
	}
	if s := PrintPiece(PieceFromName("竜")); s != "龍" {
		t.Errorf("expected=竜 actual=%s", s)
	}
	if s := PrintPiece(PieceFromName("全")); s != "成銀" {
		t.Errorf("expected=全 actual=%s", s)
	}
	if s := PrintPiece(PieceFromName("圭")); s != "成桂" {
		t.Errorf("expected=圭 actual=%s", s)
	}
	if s := PrintPiece(PieceFromName("杏")); s != "成香" {
		t.Errorf("expected=杏 actual=%s", s)
	}
}

func TestPos_Print(t *testing.T) {
	pos := &ptypes.Pos{X: 7, Y: 6}

	if v := PrintPos(pos); v != "７六" {
		t.Errorf("expected=７六 actual=%v", v)
	}
}

func TestKif_SetStartTime(t *testing.T) {
	kif := &ptypes.Kif{
		Headers: []*ptypes.Header{
			{Name: "棋戦", Value: ""},
			{Name: "開始日時", Value: "test"},
			{Name: "先手", Value: "宮尾美也"},
		},
	}
	SetStartTime(kif, time.Now())
	if l := len(kif.Headers); l != 3 {
		for i, h := range kif.Headers {
			t.Logf("header[%d]=%v:%v", i, h.Name, h.Value)
		}
		t.Errorf("header num: expected=3 actual=%d", len(kif.Headers))
	}
}
