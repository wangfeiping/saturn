package types

// Game struct for game info
type Game struct {
	AceID       string
	GameID      string
	Type        string
	Info        string
	IsGroupGame bool
}

// Play struct for game one-step-play
type Play struct {
	TxHash  string
	AceID   string
	GameID  string
	RoundID string
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
