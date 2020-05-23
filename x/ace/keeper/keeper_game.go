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
func (k Keeper) GetRound(ctx sdk.Context, start, end string) ([]types.Play, error) {
	store := ctx.KVStore(k.storeKey)
	// byteStart := []byte(start)
	// byteEnd := []byte(end)
	// it := store.Iterator(byteStart, byteEnd)
	// defer it.Close()

	var round []types.Play
	// for it.Valid() {
	// 	it.Next()
	// 	var play types.Play
	// 	err := k.cdc.UnmarshalBinaryLengthPrefixed(it.Value(), &play)
	// 	if err != nil {
	// 		fmt.Printf("get round error: %v\n", err)
	// 		return nil, err
	// 	}
	// 	round = append(round, play)
	// }

	var play types.Play
	err := k.cdc.UnmarshalBinaryLengthPrefixed(
		store.Get([]byte("LuckyAce-LuckyAce-30-:cosmos1ah0pfkn6kwuj24ks2uqqu8m387ffwnm7fvsugm")), &play)
	if err != nil {
		fmt.Printf("get round error: %v\n", err)
		return nil, err
	}
	round = append(round, play)
	return round, nil
}

// SetPlay sets the ace play to the param space.
func (k Keeper) SetPlay(ctx sdk.Context, key string, value types.Play) {
	fmt.Printf("keeper set: %s %v\n", key, value)
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set([]byte(key), bz)
}
