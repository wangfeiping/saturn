package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/types"
)

// NewQuerier creates a new querier for ace clients.
func NewQuerier(k keeper.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		fmt.Println("query: " + strings.Join(path, "\n"))
		fmt.Println("data: " + string(req.Data))
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k)
		case types.QuerySecret:
			return querySecret(ctx, k)
		case types.QueryGames:
			return queryGames(ctx, k, &req)
		case types.QueryRounds:
			return queryRounds(ctx, k, &req)
		case types.QueryPlayers:
			return queryPlayers(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown ace query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k keeper.Keeper) ([]byte, error) {
	// params := k.GetParams(ctx)

	// res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	// if err != nil {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	// }

	// return res, nil
	return nil, nil
}

func querySecret(ctx sdk.Context, k keeper.Keeper) ([]byte, error) {
	// params := k.GetParams(ctx)

	secret := types.Secret{
		Alg:    "paillier",
		Pub:    "******",
		Height: "0xHHHHHH"}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, secret)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryGames(ctx sdk.Context, k keeper.Keeper, req *abci.RequestQuery) ([]byte, error) {
	// params := k.GetParams(ctx)
	seq := strconv.FormatInt(ctx.BlockHeight(), 10)
	if len(seq) > 0 {
		seq = seq[:len(seq)-1] + "0"
	} else {
		seq = "0"
	}
	lkGame := types.Game{
		AceID:       "LuckyAce",
		GameID:      seq,
		Type:        "melee",
		IsGroupGame: false,
		Info: `
############################################################################
# Welcome! Wish you get the lucky ace!                                     #
# Game rules:                                                              #
#     1) A > K > Q > J > ... > 2                                           #
#     2) 2 > A > K ... > 6 when no 3, 4 or 5 in cards                      #
#     3) Spade > Heart > Club > Diamond when the cards are the same number #
#     3) The biggest one is winner.                                        #
############################################################################`}
	if len(req.Data) > 0 {
		name := string(req.Data)
		if !strings.EqualFold(lkGame.AceID, name) {
			lkGame.AceID = ""
			lkGame.Type = ""
			lkGame.IsGroupGame = false
			lkGame.GameID = ""
			lkGame.Info = ""
		}
		res, err := codec.MarshalJSONIndent(types.ModuleCdc, lkGame)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
		}
		return res, nil
	}
	games := []types.Game{
		lkGame,
		types.Game{
			AceID:       "Texas",
			Type:        "not_ready",
			IsGroupGame: true,
			Info:        `not ready`}}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, games)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryRounds(ctx sdk.Context, k keeper.Keeper, req *abci.RequestQuery) ([]byte, error) {
	// round := []types.Play{
	// 	types.Play{Address: "aaaaaa", Func: "draw", Args: "100chip"},
	// 	types.Play{Address: "bbbbbb", Func: "draw", Args: "1000chip"},
	// 	types.Play{Address: "cccccc", Func: "draw", Args: "10chip"}}

	var err error
	if len(req.Data) == 0 {
		err = errors.New("need game id")
		fmt.Println(err.Error())
		return nil, err
	}
	var h int64
	h, err = strconv.ParseInt(string(req.Data), 10, 64)
	if err != nil {
		fmt.Printf("wrong game id: %s %v\n", string(req.Data), err)
		return nil, err
	}
	round, err := k.GetRound(ctx,
		fmt.Sprintf("LuckyAce-%d", h), fmt.Sprintf("LuckyAce-%d", h+1))
	if err != nil {
		fmt.Printf("query round error: %v\n", err)
		return nil, err
	}

	res, err := codec.MarshalJSONIndent(types.ModuleCdc, round)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryPlayers(ctx sdk.Context, k keeper.Keeper) ([]byte, error) {

	return nil, nil
}
