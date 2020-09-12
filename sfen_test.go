package kif

import (
	"testing"

	"github.com/yunomu/kif/ptypes"
)

func TestSFENPos(t *testing.T) {
	if s := sfenPos(&ptypes.Pos{X: 1, Y: 1}); s != "1a" {
		t.Errorf("expected=1a actual=%v", s)
	}

	if s := sfenPos(&ptypes.Pos{X: 0, Y: 0}); s != "*" {
		t.Errorf("expected=* actual=%v", s)
	}

	if s := sfenPos(nil); s != "*" {
		t.Errorf("expected=* actual=%v", s)
	}

	if s := sfenPos(&ptypes.Pos{X: 3, Y: 4}); s != "3d" {
		t.Errorf("expected=1a actual=%v", s)
	}
}
