package ace

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/wangfeiping/saturn/x/ace/types"
)

// NewHandler creates an sdk.Handler for all the ace type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgAce:
			return handleMsgAce(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgAce handle request
func handleMsgAce(ctx sdk.Context, k Keeper,
	msg types.MsgAce) (*sdk.Result, error) {
	// err := k.<Action>(ctx, msg.ValidatorAddr)
	// if err != nil {
	// 	return nil, err
	// }

	fmt.Printf("handle ace msg: %s %s\n", msg.AceHash, msg.Func)

	// Define msg-ace events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule,
				types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender,
				msg.Address.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
