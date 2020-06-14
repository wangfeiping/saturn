package types

// Seed struct for random seed
type Seed struct {
	Seed         []byte
	SecurityHash []byte
}

// Play struct for game one-step-play
type Play struct {
	AceID   string
	GameID  string
	RoundID string
	Address string
	Seed    Seed
	Func    string
	Args    string
	Card    int
}
