package keeper

/*
 * Define funcs about game
 */

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wangfeiping/saturn/x/ace/types"
)

// GetRound returns the total set of ace plays.
func (k AceKeeper) GetRound(ctx sdk.Context, start, end string) ([]*types.Play, error) {
	store := ctx.KVStore(k.storeKey)
	byteStart := []byte(fmt.Sprintf("%s:%s", types.QueryPlays, start))
	byteEnd := []byte(fmt.Sprintf("%s:%s", types.QueryPlays, end))
	it := store.Iterator(byteStart, byteEnd)
	// it := store.Iterator(nil, nil)
	defer it.Close()

	var round []*types.Play
	for it.Valid() {
		var play types.Play
		err := k.cdc.UnmarshalBinaryLengthPrefixed(it.Value(), &play)
		if err != nil {
			fmt.Printf("get round error: %v\n", err)
			return nil, err
		}
		round = append(round, &play)

		it.Next()
	}

	return round, nil
}

// GetPlay returns the play of ace.
func (k AceKeeper) GetPlay(ctx sdk.Context, key string) ([]types.Play, error) {
	store := ctx.KVStore(k.storeKey)

	var round []types.Play
	var play types.Play
	err := k.cdc.UnmarshalBinaryLengthPrefixed(
		store.Get([]byte(key)), &play)
	if err != nil {
		fmt.Printf("get play error: %v\n", err)
		return nil, err
	}
	round = append(round, play)
	return round, nil
}

// SetPlay sets the ace play to the param space.
func (k AceKeeper) SetPlay(ctx sdk.Context, value types.Play) {
	k.Set(ctx, value.Key(), value)
}

// GetWinners returns the total set of winners in a game.
func (k AceKeeper) GetWinners(ctx sdk.Context, start, end string) ([]*types.Winner, error) {
	store := ctx.KVStore(k.storeKey)
	byteStart := []byte(fmt.Sprintf("%s:%s", types.QueryWinners, start))
	byteEnd := []byte(fmt.Sprintf("%s:%s", types.QueryWinners, end))
	it := store.Iterator(byteStart, byteEnd)
	// it := store.Iterator(nil, nil)
	defer it.Close()

	var winners []*types.Winner
	for it.Valid() {
		var w types.Winner
		err := k.cdc.UnmarshalBinaryLengthPrefixed(it.Value(), &w)
		if err != nil {
			fmt.Printf("get winners error: %v\n", err)
			return nil, err
		}
		winners = append(winners, &w)

		it.Next()
	}

	return winners, nil
}

// SetPlay sets the ace play to the param space.
func (k AceKeeper) SetWinner(ctx sdk.Context, value types.Winner) {
	k.Set(ctx, value.Key(), value)
}
