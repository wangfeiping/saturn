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
	AceID   string         `json:"ace_id" yaml:"ace_id"`
	GameID  int64          `json:"game_id" yaml:"game_id"`
	Action  string         `json:"action" yaml:"action"` // start, pause, cancel, end
	Address sdk.AccAddress `json:"address" yaml:"address"`
}

// NewMsgAce creates a new MsgAce instance
func NewMsgAce(aceID string, gameID int64, action string,
	addr sdk.AccAddress) *MsgAce {
	return &MsgAce{
		AceID:   aceID,
		GameID:  gameID,
		Action:  action,
		Address: addr}
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
	Seed    Seed           `json:"seed" yaml:"seed"` // random seed
	Func    string         `json:"func" yaml:"func"`
	Args    string         `json:"args" yaml:"args"`
	Address sdk.AccAddress `json:"address" yaml:"address"`
}

// NewMsgPlay creates a new MsgPlay instance
func NewMsgPlay(
	aceID string, seed Seed, function, args string,
	addr sdk.AccAddress) *MsgPlay {
	return &MsgPlay{
		AceID:   aceID,
		Seed:    seed,
		Func:    function,
		Args:    args,
		Address: addr}
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
