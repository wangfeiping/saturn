package types

import (
	"fmt"
)

// Secret Queries Result Payload for a query
type Secret struct {
	Alg    string
	Pub    string
	Height string
}

// String implement fmt.Stringer
func (s Secret) String() string {
	return fmt.Sprintf(`{"alg":"%s","pub":"%s","height":"%s"}`,
		s.Alg, s.Pub, s.Height)
}
