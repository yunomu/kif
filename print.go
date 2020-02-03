package kif

import (
	"fmt"
	"time"

	"github.com/yunomu/kif/ptypes"
)

func PrintPhase(s *ptypes.Step) string {
	if s.Seq == 0 {
		return ""
	}
	if s.Seq%2 == 0 {
		return "△"
	} else {
		return "▲"
	}
}

var pieceNames = []string{
	" ",
	"玉",
	"飛",
	"龍",
	"角",
	"馬",
	"金",
	"銀",
	"成銀",
	"桂",
	"成桂",
	"香",
	"成香",
	"歩",
	"と",
	// synonym
	"王",
	"竜",
	"全",
	"圭",
	"杏",
}

func PrintPiece(p ptypes.Piece_Id) string {
	return pieceNames[int(p)]
}

func PieceFromName(name string) ptypes.Piece_Id {
	for i, s := range pieceNames {
		if s == name {
			switch i {
			case 15:
				return ptypes.Piece_GYOKU
			case 16:
				return ptypes.Piece_RYU
			case 17:
				return ptypes.Piece_NARI_GIN
			case 18:
				return ptypes.Piece_NARI_KEI
			case 19:
				return ptypes.Piece_NARI_KYOU
			default:
				return ptypes.Piece_Id(i)
			}
		}
	}

	return ptypes.Piece_NULL
}

var (
	xstr = []rune(" １２３４５６７８９")
	ystr = []rune(" 一二三四五六七八九")
)

func PrintPos(p *ptypes.Pos) string {
	if p == nil || p.X == 0 || p.Y == 0 {
		return ""
	}
	return fmt.Sprintf("%c%c", xstr[p.X], ystr[p.Y])
}

func PrintModifier(m ptypes.Modifier_Id) string {
	switch m {
	case ptypes.Modifier_PROMOTE:
		return "成"
	case ptypes.Modifier_PUTTED:
		return "打"
	default:
		return ""
	}
}

var finStats = []string{
	" ",
	"中断",
	"投了",
	"持将棋",
	"千日手",
	"詰み",
	"切れ負け",
	"反則勝ち",
	"反則負け",
	"入玉勝ち",
}

func PrintFinishedStatus(s ptypes.FinishedStatus_Id) string {
	return finStats[int(s)]
}

func PrintMove(s *ptypes.Step) string {
	if s.FinishedStatus != ptypes.FinishedStatus_NOT_FINISHED {
		return PrintPhase(s) + PrintFinishedStatus(s.FinishedStatus)
	}

	var src string
	if s.Src != nil {
		src = fmt.Sprintf("(%d%d)", s.Src.X, s.Src.Y)
	}

	return fmt.Sprintf("%s%s%s%s%s",
		PrintPhase(s),
		PrintPos(s.Dst),
		PrintPiece(s.Piece),
		PrintModifier(s.Modifier),
		src,
	)
}

func PrintThinking(sec int32) string {
	m := sec / 60
	s := sec % 60
	return fmt.Sprintf("%2d:%02d", m, s)
}

func PrintElapsed(sec int32) string {
	h := sec / 3600
	m := sec / 60
	s := sec % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

const (
	timeFormat    = "2006/01/02 15:04:05"
	startTimeName = "開始日時"
	endTimeName   = "終了日時"
)

func setTime(k *ptypes.Kif, name string, t time.Time) {
	var idx = -1
	for i, h := range k.Headers {
		if h.Name == name {
			idx = i
			break
		}
	}
	h := &ptypes.Header{
		Name:  name,
		Value: t.Format(timeFormat),
	}
	if idx != -1 {
		k.Headers = append(k.Headers[:idx], k.Headers[idx+1:]...)
	}
	k.Headers = append(k.Headers, h)
}

func SetStartTime(k *ptypes.Kif, t time.Time) {
	setTime(k, startTimeName, t)
}

func SetEndTime(k *ptypes.Kif, t time.Time) {
	setTime(k, endTimeName, t)
}
