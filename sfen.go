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

func StepToMove(step *ptypes.Step) string {
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

	write := func(s string) error {
		_, err := w.Write([]byte(s))
		return err
	}

	if err := write("position startpos moves"); err != nil {
		return err
	}

	for _, step := range steps {
		if step.FinishedStatus != ptypes.FinishedStatus_NOT_FINISHED {
			break
		}

		if err := write(" " + StepToMove(step)); err != nil {
			return err
		}
	}

	return nil
}
