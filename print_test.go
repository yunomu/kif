package kif

import (
	"testing"

	"time"
)

func TestPieceFronName(t *testing.T) {
	for i := 0; i <= 14; i++ {
		e := pieceNames[i]
		a := PieceFromName(e).Print()
		if e != a {
			t.Errorf("expected=%s actual=%s", e, a)
		}
	}

	if s := PieceFromName("王").Print(); s != "玉" {
		t.Errorf("expected=玉 actual=%s", s)
	}
	if s := PieceFromName("竜").Print(); s != "龍" {
		t.Errorf("expected=竜 actual=%s", s)
	}
	if s := PieceFromName("全").Print(); s != "成銀" {
		t.Errorf("expected=全 actual=%s", s)
	}
	if s := PieceFromName("圭").Print(); s != "成桂" {
		t.Errorf("expected=圭 actual=%s", s)
	}
	if s := PieceFromName("杏").Print(); s != "成香" {
		t.Errorf("expected=杏 actual=%s", s)
	}
}

func TestPos_Print(t *testing.T) {
	pos := &Pos{X: 7, Y: 6}

	if v := pos.Print(); v != "７六" {
		t.Errorf("expected=７六 actual=%v", v)
	}
}

func TestKif_SetStartTime(t *testing.T) {
	kif := &Kif{
		Headers: []*Header{
			{Name: "棋戦", Value: ""},
			{Name: "開始日時", Value: "test"},
			{Name: "先手", Value: "宮尾美也"},
		},
	}
	kif.SetStartTime(time.Now())
	if l := len(kif.Headers); l != 3 {
		for i, h := range kif.Headers {
			t.Logf("header[%d]=%v:%v", i, h.Name, h.Value)
		}
		t.Errorf("header num: expected=3 actual=%d", len(kif.Headers))
	}
}
