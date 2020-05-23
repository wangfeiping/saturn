package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAce{}
var _ sdk.Msg = &MsgPlay{}

// MsgAce - struct for update request of ace
type MsgAce struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
}

// NewMsgAce creates a new MsgAce instance
func NewMsgAce(addr sdk.AccAddress) (*MsgAce, error) {
	return &MsgAce{Address: addr}, nil
}

// Route returns the name of module
func (msg MsgAce) Route() string { return RouterKey }

// Type returns action
func (msg MsgAce) Type() string { return "Ace" }

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

// MsgPlay - struct for one-step-play request of game
type MsgPlay struct {
	AceID   string         `json:"ace_id" yaml:"ace_id"`
	GameID  string         `json:"game_id" yaml:"game_id"`
	RoundID string         `json:"round_id" yaml:"round_id"`
	Seed    Seed           `json:"seed" yaml:"seed"` // random seed
	Func    string         `json:"func" yaml:"func"`
	Args    string         `json:"args" yaml:"args"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
}

// NewMsgPlay creates a new MsgPlay instance
func NewMsgPlay(
	aceID, gameID, roundID string,
	seed *Seed, function, args string,
	addr sdk.AccAddress) (*MsgPlay, error) {
	return &MsgPlay{
		AceID:   aceID,
		GameID:  gameID,
		RoundID: roundID,
		Seed:    *seed,
		Func:    function,
		Args:    args,
		Address: addr}, nil
}

// Route returns the name of module
func (msg MsgPlay) Route() string { return RouterKey }

// Type returns action
func (msg MsgPlay) Type() string { return "Play" }

// GetSigners defines whose signature is required
func (msg MsgPlay) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgPlay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgPlay) ValidateBasic() error {
	fmt.Printf("msg.Address.String(): %s\n", msg.Address.String())
	if msg.Address.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address")
	}
	return nil
}
