package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/wangfeiping/saturn/x/ace/types"
)

// AceKeeper of the ace store
type AceKeeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	// paramspace types.ParamSubspace
}

// NewKeeper creates a ace keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) AceKeeper {
	keeper := AceKeeper{
		storeKey: key,
		cdc:      cdc,
		// paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k AceKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k AceKeeper) Set(ctx sdk.Context, key string, value interface{}) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(value)
	store.Set([]byte(key), bz)
}

func (k AceKeeper) Has(ctx sdk.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(key))
}

// Get the entire Ace metadata struct for a key
func (k AceKeeper) Get(ctx sdk.Context, key string, ptr interface{}) error {
	store := ctx.KVStore(k.storeKey)
	byteKey := []byte(key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), ptr)
	return err
}

func (k AceKeeper) Delete(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(key))
}
