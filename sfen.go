package kif

import (
	"fmt"
	"strings"

	"github.com/yunomu/kif/ptypes"
)

func sfenPos(pos *ptypes.Pos) string {
	if pos == nil || pos.X == 0 || pos.Y == 0 {
		return "*"
	}

	return fmt.Sprintf("%d%c", pos.X, 'a'+pos.Y-1)
}

var sfenPiece = []string{
	" ",
	"K",
	"R",
	"+R",
	"B",
	"+B",
	"G",
	"S",
	"+S",
	"N",
	"+N",
	"L",
	"+L",
	"P",
	"+P",
}

func StepToSFEN(step *ptypes.Step) string {
	var drop string
	if step.Modifier == ptypes.Modifier_PUTTED {
		drop = sfenPiece[int(step.Piece)]
	}

	var prom string
	if step.Modifier == ptypes.Modifier_PROMOTE {
		prom = "+"
	}

	return drop + sfenPos(step.Src) + sfenPos(step.Dst) + prom
}

func ToSFEN(steps []*ptypes.Step) string {
	var ss []string
	for _, step := range steps {
		ss = append(ss, StepToSFEN(step))
	}

	return strings.Join(ss, " ")
}
