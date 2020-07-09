package keeper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"github.com/wangfeiping/saturn/app"
	"github.com/wangfeiping/saturn/x/ace"
	acekeeper "github.com/wangfeiping/saturn/x/ace/keeper"
)

func TestKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ace Keeper Suite")
}

func CreateMockAceKeeper() (sdk.Context, acekeeper.AceKeeper) {
	cdc := app.MakeCodec()
	aceStoreKey := sdk.NewKVStoreKey(ace.StoreKey)
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(aceStoreKey, sdk.StoreTypeIAVL, db)
	cms.LoadLatestVersion()
	return sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger()),
		acekeeper.NewKeeper(cdc, aceStoreKey)
}
