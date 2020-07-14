package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAce{}, "Saturn/MsgAce", nil)
	cdc.RegisterConcrete(MsgPlay{}, "Saturn/MsgPlay", nil)
	cdc.RegisterConcrete(Secret{}, "Saturn/Secret", nil)
	cdc.RegisterConcrete(Game{}, "Saturn/Game", nil)
	cdc.RegisterConcrete(Play{}, "Saturn/Play", nil)
	cdc.RegisterConcrete(Winner{}, "Saturn/Winner", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
