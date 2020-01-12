package ptypes

import (
	"fmt"
	"time"
)

func (s *Step) PrintPhase() string {
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

func (p Piece_Id) Print() string {
	return pieceNames[int(p)]
}

func PieceFromName(name string) Piece_Id {
	for i, s := range pieceNames {
		if s == name {
			switch i {
			case 15:
				return Piece_GYOKU
			case 16:
				return Piece_RYU
			case 17:
				return Piece_NARI_GIN
			case 18:
				return Piece_NARI_KEI
			case 19:
				return Piece_NARI_KYOU
			default:
				return Piece_Id(i)
			}
		}
	}

	return Piece_NULL
}

var (
	xstr = []rune(" １２３４５６７８９")
	ystr = []rune(" 一二三四五六七八九")
)

func (p *Pos) Print() string {
	if p == nil || p.X == 0 || p.Y == 0 {
		return ""
	}
	return fmt.Sprintf("%c%c", xstr[p.X], ystr[p.Y])
}

func (m Modifier_Id) Print() string {
	switch m {
	case Modifier_PROMOTE:
		return "成"
	case Modifier_PUTTED:
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

func (s FinishedStatus_Id) Print() string {
	return finStats[int(s)]
}

func (s *Step) PrintMove() string {
	if s.FinishedStatus != FinishedStatus_NOT_FINISHED {
		return s.PrintPhase() + s.FinishedStatus.Print()
	}

	var src string
	if s.Src != nil {
		src = fmt.Sprintf("(%d%d)", s.Src.X, s.Src.Y)
	}

	return fmt.Sprintf("%s%s%s%s%s",
		s.PrintPhase(),
		s.Dst.Print(),
		s.Piece.Print(),
		s.Modifier.Print(),
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

func (k *Kif) setTime(name string, t time.Time) {
	var idx = -1
	for i, h := range k.Headers {
		if h.Name == name {
			idx = i
			break
		}
	}
	h := &Header{
		Name:  name,
		Value: t.Format(timeFormat),
	}
	if idx != -1 {
		k.Headers = append(k.Headers[:idx], k.Headers[idx+1:]...)
	}
	k.Headers = append(k.Headers, h)
}

func (k *Kif) SetStartTime(t time.Time) {
	k.setTime(startTimeName, t)
}

func (k *Kif) SetEndTime(t time.Time) {
	k.setTime(endTimeName, t)
}
