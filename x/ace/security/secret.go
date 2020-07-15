package security

import (
	"encoding/json"
)

// Secret the payload for security key transmission
type Secret struct {
	Alg    string
	Hash   string
	Height int64
}

// String implement fmt.Stringer
func (s Secret) String() string {
	bz, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(bz)
}
