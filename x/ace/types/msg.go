package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAce{}

// MsgAce - struct for game play request
type MsgAce struct {
	Seed    *Seed          `json:"seed" yaml:"seed"` // random seed
	Func    string         `json:"func" yaml:"func"`
	Args    string         `json:"args" yaml:"args"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
}

// NewMsgAce creates a new MsgAce instance
func NewMsgAce(seed *Seed, function, args string) MsgAce {
	return MsgAce{
		Seed: seed, Func: function, Args: args}
}

// AceConst const of Ace
const AceConst = "Ace"

// Route returns route key
func (msg MsgAce) Route() string { return RouterKey }

// Type returns message type
func (msg MsgAce) Type() string { return AceConst }

// GetSigners returns message signers
func (msg MsgAce) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Address)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgAce) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgAce) ValidateBasic() error {
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing validator address")
	}
	return nil
}
