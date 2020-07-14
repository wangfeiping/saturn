package types

import (
	"fmt"
)

// Winner struct for games
type Winner struct {
	AceID   string
	Height  int64
	Rank    int
	Address string
	Chip    Chips
	Win     Chips
}

// Winners struct for game winners
type Winners struct {
	plays    []*Play
	exist345 bool
}

func NewWinners(plays []*Play, exist345 bool) *Winners {
	return &Winners{
		plays:    plays,
		exist345: exist345}
}

func (w Winner) Key() string {
	return fmt.Sprintf("%s:%s:%d:%s",
		QueryWinners, CreateGameID(w.AceID, w.Height), w.Rank, w.Address)
}

// Len implements sort.Interface
func (w Winners) Len() int {
	return len(w.plays)
}

// Swap implements sort.Interface
func (w Winners) Swap(i, j int) {
	w.plays[i], w.plays[j] = w.plays[j], w.plays[i]
}

// Less implements sort.Interface
func (w Winners) Less(i, j int) bool {
	cardi := w.plays[i].Card
	cardj := w.plays[j].Card
	if !w.exist345 {
		if cardi < 4 {
			// 0: "2 diamond", 1: "2 club", 2: "2 heart", 3: "2 spade",
			cardi += 51
		}
		if cardj < 4 {
			// same with above
			cardj += 51
		}
	}
	// fmt.Printf("%d:%d %t %d:%d\n", i, cardi, cardi > cardj, j, cardj)
	return cardi > cardj
}

func (w Winners) GetWinner(i int) Winner {
	play := w.plays[i]
	// var args = map[string]string{}
	// args["card"] = fmt.Sprintf("%d", play.Card)
	// fmt.Printf("winner sort %d: %s; %d\n", i, play.Address, play.Card)
	return Winner{
		AceID:   play.AceID,
		Height:  play.Height,
		Address: play.Address}
}
