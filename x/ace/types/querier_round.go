package types

// Seed struct for random seed
type Seed struct {
	Seed         []byte
	SecurityHash []byte
}

// Round struct for game play
type Round struct {
	Address string
	Seed    Seed
	Func    string
	Args    string
}
