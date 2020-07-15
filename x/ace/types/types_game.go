package types

import (
	"fmt"

	"github.com/wangfeiping/saturn/x/ace/security"
)

// Seed struct for random seed
type Seed struct {
	Secret security.Secret
	Hash   []byte
}

// Game struct for game info
type Game struct {
	AceID       string
	GameID      int64
	Type        string
	Info        string
	IsGroupGame bool
}

// Play struct for game one-step-play
type Play struct {
	TxHash  string
	AceID   string
	Height  int64
	RoundID int
	Address string
	Seed    Seed
	Func    string
	Args    string
	Card    int
}

type Chips struct {
	Amount int
	Denom  string
}

func (p Play) Key() string {
	return fmt.Sprintf("%s:%s:%d:%s",
		QueryPlays, CreateGameID(p.AceID, p.Height), p.RoundID, p.Address)
}
func CreateGameID(aceID string, height int64) string {
	return fmt.Sprintf("%s:%d", aceID, height)
}
