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
func (k AceKeeper) GetRound(ctx sdk.Context, start, end string) ([]types.Play, error) {
	store := ctx.KVStore(k.storeKey)
	// byteStart := []byte("LuckyAce-LuckyAce-30-:cosmos1afaq874drmzn8lg00jmzlsn2lwdkk5qwhgdjnx")
	// byteEnd := []byte("LuckyAce-LuckyAce-30-:cosmos1ah0pfkn6kwuj24ks2uqqu8m387ffwnm7fvsugm")
	byteStart := []byte(start)
	byteEnd := []byte(end)
	it := store.Iterator(byteStart, byteEnd)
	// it := store.Iterator(nil, nil)
	defer it.Close()

	var round []types.Play
	for it.Valid() {
		// it.Next()
		// fmt.Printf("get key: %s\n", string(it.Key()))
		// fmt.Printf("get play: %s\n", string(it.Value()))
		var play types.Play
		err := k.cdc.UnmarshalBinaryLengthPrefixed(it.Value(), &play)
		if err != nil {
			fmt.Printf("get round error: %v\n", err)
			return nil, err
		}
		round = append(round, play)

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
func (k AceKeeper) SetPlay(ctx sdk.Context, key string, value types.Play) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set([]byte(key), bz)
}
