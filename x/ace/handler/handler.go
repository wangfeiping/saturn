package handler

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/types"
)

var CARDS = [...]string{
	"2 diamond", "2 club", "2 heart", "2 spade",
	"3 diamond", "3 club", "3 heart", "3 spade",
	"4 diamond", "4 club", "4 heart", "4 spade",
	"5 diamond", "5 club", "5 heart", "5 spade",
	"6 diamond", "6 club", "6 heart", "6 spade",
	"7 diamond", "7 club", "7 heart", "7 spade",
	"8 diamond", "8 club", "8 heart", "8 spade",
	"9 diamond", "9 club", "9 heart", "9 spade",
	"10 diamond", "10 club", "10 heart", "10 spade",
	"J diamond", "J club", "J heart", "J spade",
	"Q diamond", "Q club", "Q heart", "Q spade",
	"K diamond", "K club", "K heart", "K spade",
	"A diamond", "A club", "A heart", "A spade",
}

const DESC = "The %s of %s"

// NewHandler creates an sdk.Handler for all the ace type messages
func NewHandler(k keeper.AceKeeper, bank types.BankKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgPlay:
			return handleMsgPlay(ctx, k, bank, msg)
		// case types.MsgAce:
		// 	return handleMsgAce(ctx, k, bank, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgPlay handle request
func handleMsgPlay(ctx sdk.Context, k keeper.AceKeeper, bank types.BankKeeper,
	m types.MsgPlay) (*sdk.Result, error) {

	var err error
	h := checkGameStartHeight(ctx)
	if h <= ctx.BlockHeight()-types.GameDurationHeight {
		err = fmt.Errorf("game is over: %d current: %d", h, ctx.BlockHeight())
		ctx.Logger().With("module", "ace").
			Error("game is over", "game", h, "height", ctx.BlockHeight())
		return nil, err
	}
	// TODO How to get the TX's hash, Resolving conflicts over play TXs
	play := types.Play{
		// TxHash:  ,
		AceID:   m.AceID,
		Height:  h,
		Address: m.Address.String(),
		Seed:    m.Seed,
		Func:    m.Func,
		Args:    m.Args}
	k.SetPlay(ctx, play)
	// args := strings.Split(play.Args, ",")
	coins, err := sdk.ParseCoins("1000000chip")
	if err != nil {
		ctx.Logger().With("module", "ace").
			Error("parse coins failed", "error", err.Error())
		return nil, err
	}
	pooler, err := sdk.AccAddressFromBech32(types.PoolerAddress)
	if err != nil {
		ctx.Logger().With("module", "ace").
			Error("parse pooler's address failed", "error", err.Error())
		return nil, err
	}
	err = bank.SendCoins(ctx, m.Address, pooler, coins)
	if err != nil {
		ctx.Logger().With("module", "ace").
			Error("send coins failed", "error", err.Error())
		return nil, err
	}
	// fmt.Printf("handle play msg: %d - %s %s %s\n", ctx.BlockHeight(),
	// 	m.AceID, m.Func, m.Args)

	// Define msg-play events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule,
				types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender,
				m.Address.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
