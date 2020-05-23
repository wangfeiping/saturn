package handler

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/types"
)

// NewHandler creates an sdk.Handler for all the ace type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgPlay:
			return handleMsgPlay(ctx, k, msg)
		// case types.MsgAce:

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgPlay handle request
func handleMsgPlay(ctx sdk.Context, k keeper.Keeper,
	m types.MsgPlay) (*sdk.Result, error) {
	// err := k.<Action>(ctx, msg.ValidatorAddr)
	// if err != nil {
	// 	return nil, err
	// }

	fmt.Printf("handle play msg: %d - %s %s %s\n", ctx.BlockHeight(),
		m.AceID, m.Func, m.Args)

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
