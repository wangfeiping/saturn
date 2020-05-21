package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/wangfeiping/saturn/x/ace/types"
)

// NewQuerier creates a new querier for ace clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k)
		case types.QuerySecret:
			return querySecret(ctx, k)
		case types.QueryRounds:
			return queryRounds(ctx, k)
		case types.QueryPlayers:
			return queryPlayers(ctx, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown ace query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {
	// params := k.GetParams(ctx)

	// res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	// if err != nil {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	// }

	// return res, nil
	return nil, nil
}

func querySecret(ctx sdk.Context, k Keeper) ([]byte, error) {
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

func queryRounds(ctx sdk.Context, k Keeper) ([]byte, error) {
	// params := k.GetParams(ctx)

	rounds := []types.Round{
		types.Round{Address: "aaaaaa", Func: "draw", Args: "100chip"},
		types.Round{Address: "bbbbbb", Func: "draw", Args: "1000chip"},
		types.Round{Address: "cccccc", Func: "draw", Args: "10chip"}}
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, rounds)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func queryPlayers(ctx sdk.Context, k Keeper) ([]byte, error) {

	return nil, nil
}
