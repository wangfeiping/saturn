package types

import (
	"fmt"
)

// Winner struct for games
type Winner struct {
	AceID   string
	GameID  string
	Address string
	Chip    Chips
	Win     Chips
	Args    map[string]string
}

// Winners struct for game winners
type Winners struct {
	plays    []Play
	exist345 bool
}

func NewWinners(plays []Play, exist345 bool) *Winners {
	return &Winners{
		plays:    plays,
		exist345: exist345}
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
			cardi += 51
		}
		if cardj < 4 {
			cardj += 51
		}
	}
	return cardi > cardj
}

func (w Winners) GetWinner(i int) Winner {
	play := w.plays[i]
	var args = map[string]string{}
	args["card"] = fmt.Sprintf("%d", play.Card)
	return Winner{
		AceID:   play.AceID,
		GameID:  play.GameID,
		Address: play.Address,
		Args:    args}
}
