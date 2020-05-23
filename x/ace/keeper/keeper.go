package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wangfeiping/saturn/x/ace/types"
)

// Keeper of the ace store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        *codec.Codec
	paramspace types.ParamSubspace
}

// NewKeeper creates a ace keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	keeper := Keeper{
		storeKey: key,
		cdc:      cdc,
		// paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Get the entire Ace metadata struct for a key
func (k Keeper) Get(ctx sdk.Context, key string) (types.Secret, error) {
	store := ctx.KVStore(k.storeKey)
	var item types.Secret
	byteKey := []byte(key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &item)
	return item, err
}

func (k Keeper) set(ctx sdk.Context, key string, value types.Secret) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set([]byte(key), bz)
}

func (k Keeper) delete(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(key))
}
