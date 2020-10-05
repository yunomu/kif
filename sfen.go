package kif

import (
	"fmt"
	"io"

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

func writeSFEN(w io.Writer, steps []*ptypes.Step) error {
	if len(steps) == 0 || steps[0].FinishedStatus != ptypes.FinishedStatus_NOT_FINISHED {
		return nil
	}

	if _, err := w.Write([]byte("position startpos moves")); err != nil {
		return err
	}

	sp := []byte(" ")
	for _, step := range steps {
		if step.FinishedStatus != ptypes.FinishedStatus_NOT_FINISHED {
			break
		}

		if _, err := w.Write(append(sp, []byte(StepToSFEN(step))...)); err != nil {
			return err
		}
	}

	return nil
}
