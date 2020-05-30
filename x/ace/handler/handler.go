package handler

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/types"
)

// NewHandler creates an sdk.Handler for all the ace type messages
func NewHandler(k keeper.Keeper, bank types.BankKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgPlay:
			fmt.Println("handle play msg")
			return handleMsgPlay(ctx, k, bank, msg)
		case types.MsgAce:
			fmt.Println("handle ace msg")
			return nil, fmt.Errorf("not support")
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgPlay handle request
func handleMsgPlay(ctx sdk.Context, k keeper.Keeper, bank types.BankKeeper,
	m types.MsgPlay) (*sdk.Result, error) {

	h, err := strconv.ParseInt(m.GameID, 10, 64)
	if err != nil {
		fmt.Printf("wrong game id: %s %v\n", m.GameID, err)
		return nil, err
	}
	if h <= ctx.BlockHeight()-types.GameDurationHeight {
		err = fmt.Errorf("game is over: %d current: %d", h, ctx.BlockHeight())
		fmt.Println(err.Error())
		return nil, err
	}
	play := types.Play{
		AceID:   m.AceID,
		GameID:  m.GameID,
		RoundID: m.RoundID,
		Address: m.Address.String(),
		Seed:    m.Seed,
		Func:    m.Func,
		Args:    m.Args}
	k.SetPlay(ctx, fmt.Sprintf("%s-%s-%s:%s",
		play.AceID, play.GameID, play.RoundID, play.Address), play)
	args := strings.Split(play.Args, ",")
	coins, err := sdk.ParseCoins(args[0])
	if err != nil {
		fmt.Println("parse coins error: " + err.Error())
		return nil, err
	}
	dealerAddr := "cosmos1kl86mq7264f8x4pumdk7la7w5svm8lep6626vz"
	dealer, err := sdk.AccAddressFromBech32(dealerAddr)
	if err != nil {
		fmt.Println("parse dealer's address error: " + err.Error())
		return nil, err
	}
	err = bank.SendCoins(ctx, m.Address, dealer, coins)
	if err != nil {
		fmt.Println("send coins error: " + err.Error())
		return nil, err
	}
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
