package kif

import (
	"sort"

	"github.com/yunomu/kif/ptypes"
)

type stepSlice []*ptypes.Step

func (s stepSlice) Len() int               { return len(s) }
func (s stepSlice) Less(i int, j int) bool { return s[i].GetSeq() < s[j].GetSeq() }
func (s stepSlice) Swap(i int, j int)      { s[i], s[j] = s[j], s[i] }

var stdHeaders = []string{
	"対局日",
	"開始日時",
	"終了日時",
	"棋戦",
	"手合割",
	"先手",
	"後手",
	"戦型",
	"表題",
	"持ち時間",
	"消費時間",
	"場所",
	"掲載",
	"備考",
	"先手省略名",
	"後手省略名",
}

type headerSlice []*ptypes.Header

func (h headerSlice) Len() int { return len(h) }

func findHeader(name string) int {
	for i, s := range stdHeaders {
		if s == name {
			return i
		}
	}
	return -1
}

func (h headerSlice) Less(i, j int) bool {
	ji := findHeader(h[j].Name)
	if ji == -1 {
		return true
	}
	return findHeader(h[i].Name) < ji
}

func (h headerSlice) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func Normalize(k *ptypes.Kif) {
	sort.Sort(headerSlice(k.Headers))
	sort.Sort(stepSlice(k.Steps))
}
