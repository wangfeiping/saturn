package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAce{}

// MsgAce - struct for game play request
type MsgAce struct {
	AceHash string         `json:"ace_hash" yaml:"ace_hash"`
	Seed    *Seed          `json:"seed" yaml:"seed"` // random seed
	Func    string         `json:"func" yaml:"func"`
	Args    string         `json:"args" yaml:"args"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
}

// NewMsgAce creates a new MsgAce instance
func NewMsgAce(
	aceHash string, seed *Seed,
	function, args string,
	addr sdk.AccAddress) (*MsgAce, error) {
	return &MsgAce{
		AceHash: aceHash,
		Seed:    seed,
		Func:    function,
		Args:    args,
		Address: addr}, nil
}

// AceConst const of Ace
const AceConst = "Ace"

// Route returns the name of module
func (msg MsgAce) Route() string { return RouterKey }

// Type returns action
func (msg MsgAce) Type() string { return AceConst }

// GetSigners defines whose signature is required
func (msg MsgAce) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgAce) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgAce) ValidateBasic() error {
	fmt.Printf("msg.Address.String(): %s\n", msg.Address.String())
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}
	return nil
}
