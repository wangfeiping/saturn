package handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"github.com/wangfeiping/saturn/app"
	"github.com/wangfeiping/saturn/x/ace"
	acekeeper "github.com/wangfeiping/saturn/x/ace/keeper"
)

var testCdc *codec.Codec = app.MakeCodec()
var aceStoreKey *sdk.KVStoreKey = sdk.NewKVStoreKey(ace.StoreKey)
var authStoreKey *sdk.KVStoreKey = sdk.NewKVStoreKey(auth.StoreKey)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ace Handler Suite")
}

func CreateMockSdkContext() sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	cms.MountStoreWithDB(aceStoreKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(authStoreKey, sdk.StoreTypeIAVL, db)

	cms.LoadLatestVersion()
	return sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
}

func CreateMockAceKeeper() acekeeper.AceKeeper {
	return acekeeper.NewKeeper(testCdc, aceStoreKey)
}

func CreateMockParamsKeeper() params.Keeper {
	return params.NewKeeper(testCdc,
		sdk.NewKVStoreKey(params.StoreKey),
		sdk.NewTransientStoreKey(params.TStoreKey))
}

func CreateMockAccountKeeper(paramsKeeper params.Keeper) auth.AccountKeeper {
	subspace := paramsKeeper.Subspace(auth.ModuleName)

	return auth.NewAccountKeeper(
		testCdc,
		authStoreKey,
		subspace,
		auth.ProtoBaseAccount,
	)
}

func CreateMockBankKeeper(ak auth.AccountKeeper, pk params.Keeper) bank.Keeper {
	return bank.NewBaseKeeper(
		ak,
		pk.Subspace(bank.ModuleName),
		// app.ModuleAccountAddrs(),
		make(map[string]bool),
	)
}
