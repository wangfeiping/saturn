package handler_test

import (
	"fmt"
	"math/big"
	"time"

	// "strconv"
	// "strings"
	// "testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/wangfeiping/saturn/x/ace/handler"
	"github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/security"
	"github.com/wangfeiping/saturn/x/ace/security/paillier"
	"github.com/wangfeiping/saturn/x/ace/types"
)

var seeds []int64 = []int64{0, 1, 1, 0, 0, 0}

var _ = Describe("x/ace/handler", func() {
	var (
		denom   string = "chip"
		balance int64  = 9999999
		num     int    = 5

		pubkey        security.PublicKey
		addrs         []sdk.AccAddress
		ctx           sdk.Context
		aceKeeper     keeper.AceKeeper
		accountKeeper auth.AccountKeeper
		handle        sdk.Handler
		query         sdk.Querier
	)

	ctx = CreateMockSdkContext()
	aceKeeper = CreateMockAceKeeper()
	pk := CreateMockParamsKeeper()
	accountKeeper = CreateMockAccountKeeper(pk)
	bankKeeper := CreateMockBankKeeper(accountKeeper, pk)
	secret := security.Secret{}
	pubkey = paillier.ResumePubKey(secret)

	handle = handler.NewHandler(aceKeeper, bankKeeper)
	query = handler.NewQuerier(aceKeeper)

	for i := 0; i < num; i++ {
		addrs = append(addrs, sdk.AccAddress([]byte(
			fmt.Sprintf("addr.test.%d", i))))
		acc := accountKeeper.NewAccountWithAddress(ctx, addrs[i])
		accountKeeper.SetAccount(ctx, acc)
		bankKeeper.SetCoins(ctx, addrs[i],
			sdk.NewCoins(sdk.NewInt64Coin(denom, balance)))
	}

	BeforeEach(func() {

	})

	Describe("Create an ace handler", func() {
		Context("with mock keepers", func() {
			It("should be success", func() {
				Expect(handle).NotTo(BeZero())
			})
		})
	})

	Describe("Call ace handler", func() {
		Context(fmt.Sprintf("with %d play messages", num), func() {
			It("should be success", func() {
				ctx = ctx.WithBlockHeader(
					abci.Header{Height: 0, Time: time.Unix(10, 0)})
				for i := 0; i < num; i++ {
					hash, err := pubkey.Encrypt(
						big.NewInt(seeds[i]).Bytes())
					Expect(err).ShouldNot(HaveOccurred())
					seed := types.Seed{Secret: secret, Hash: hash}
					msg := types.NewMsgPlay(
						"LuckyAce", seed, "draw", "", addrs[i])
					_, err = handle(ctx, *msg)
					Expect(err).ShouldNot(HaveOccurred())
				}
			})
		})

		Context("Check the balance of all accounts", func() {
			It("should be success", func() {
				for i := 0; i < num; i++ {
					coins := bankKeeper.GetCoins(ctx, addrs[i])
					Expect(coins[0].Amount.Int64()).To(Equal(balance - 1000000))
				}
			})
		})

		Context("end the game", func() {
			It("should be success", func() {
				// Updates the block height
				ctx = ctx.WithBlockHeader(
					abci.Header{Height: 9, Time: time.Unix(10, 0)})
				handler.EndBlockHandle(ctx, abci.RequestEndBlock{},
					aceKeeper, bankKeeper)

				req := abci.RequestQuery{Data: []byte("0")}
				res, err := handler.QueryAllPlays(ctx, aceKeeper, &req)
				Expect(err).ShouldNot(HaveOccurred())
				var plays []*types.Play
				testCdc.MustUnmarshalJSON(res, &plays)
				// for i, p := range plays {
				// 	fmt.Printf("play: %d %s %d %s\n", i, p.Address, p.Card, handler.CARDS[p.Card])
				// }

				res, err = query(ctx, []string{types.QueryWinners}, req)
				Expect(err).ShouldNot(HaveOccurred())

				var out []*types.Winner
				testCdc.MustUnmarshalJSON(res, &out)
				// fmt.Printf("!!!winner: %s\n", out[0].Address)
				Expect(out[0].Address).To(Equal("cosmos1v9jxgu3ww3jhxapwxvj6q2am"))
			})
		})
	})
})

var _ = Describe("AceHandler", func() {
	var (
		denom   string = "chip"
		balance int64  = 9999999
		num     int    = 6

		pubkey        security.PublicKey
		addrs         []sdk.AccAddress
		ctx           sdk.Context
		aceKeeper     keeper.AceKeeper
		accountKeeper auth.AccountKeeper
		handle        sdk.Handler
		query         sdk.Querier
	)

	ctx = CreateMockSdkContext()
	aceKeeper = CreateMockAceKeeper()
	pk := CreateMockParamsKeeper()
	accountKeeper = CreateMockAccountKeeper(pk)
	bankKeeper := CreateMockBankKeeper(accountKeeper, pk)
	secret := security.Secret{}
	pubkey = paillier.ResumePubKey(secret)

	handle = handler.NewHandler(aceKeeper, bankKeeper)
	query = handler.NewQuerier(aceKeeper)

	for i := 0; i < num; i++ {
		addrs = append(addrs, sdk.AccAddress([]byte(
			fmt.Sprintf("addr.test.%d", i))))
		acc := accountKeeper.NewAccountWithAddress(ctx, addrs[i])
		accountKeeper.SetAccount(ctx, acc)
		bankKeeper.SetCoins(ctx, addrs[i],
			sdk.NewCoins(sdk.NewInt64Coin(denom, balance)))
	}

	BeforeEach(func() {

	})

	Describe("Create an ace handler", func() {
		Context("with mock keepers", func() {
			It("should be success", func() {
				Expect(handle).NotTo(BeZero())
			})
		})
	})

	Describe("Call ace handler", func() {
		Context(fmt.Sprintf("with %d play messages", num), func() {
			It("should be success", func() {
				ctx = ctx.WithBlockHeader(
					abci.Header{Height: 0, Time: time.Unix(10, 0)})
				for i := 0; i < num; i++ {
					hash, err := pubkey.Encrypt(
						big.NewInt(seeds[i]).Bytes())
					Expect(err).ShouldNot(HaveOccurred())
					seed := types.Seed{Secret: secret, Hash: hash}
					msg := types.NewMsgPlay(
						"LuckyAce", seed, "draw", "", addrs[i])
					_, err = handle(ctx, *msg)
					Expect(err).ShouldNot(HaveOccurred())
				}
			})
		})

		Context("Check the balance of all accounts", func() {
			It("should be success", func() {
				for i := 0; i < num; i++ {
					coins := bankKeeper.GetCoins(ctx, addrs[i])
					Expect(coins[0].Amount.Int64()).To(Equal(balance - 1000000))
				}
			})
		})

		Context("end the game", func() {
			It("should be success", func() {
				// Updates the block height
				ctx = ctx.WithBlockHeader(
					abci.Header{Height: 9, Time: time.Unix(10, 0)})
				handler.EndBlockHandle(ctx, abci.RequestEndBlock{},
					aceKeeper, bankKeeper)

				req := abci.RequestQuery{Data: []byte("0")}
				res, err := handler.QueryAllPlays(ctx, aceKeeper, &req)
				Expect(err).ShouldNot(HaveOccurred())
				var plays []*types.Play
				testCdc.MustUnmarshalJSON(res, &plays)
				// for i, p := range plays {
				// 	fmt.Printf("play: %d %s %d %s\n", i, p.Address, p.Card, handler.CARDS[p.Card])
				// }

				res, err = query(ctx, []string{types.QueryWinners}, req)
				Expect(err).ShouldNot(HaveOccurred())

				var out []types.Winner
				testCdc.MustUnmarshalJSON(res, &out)
				// fmt.Printf("!!!winner: %s\n", out[0].Address)
				Expect(out[0].Address).To(Equal("cosmos1v9jxgu3ww3jhxapwxyvwyxrx"))
			})
		})
	})
})

// func Test_SortWinners(t *testing.T) {
// 	var plays []types.Play
// 	num := 100
// 	mrand.Seed(1)
// 	perm := mrand.Perm(num)
// 	fmt.Println("math/rand.Perm: ", perm)

// 	for i, n := range perm {
// 		p := types.Play{
// 			AceID:   "LuckyAce",
// 			GameID:  "1000",
// 			Address: fmt.Sprintf("xxx_%d", i),
// 			Card:    checkCardNum(n)}
// 		plays = append(plays, p)
// 	}

// 	winners := SortWinners(plays, 33)
// 	for _, w := range winners {
// 		i, err := strconv.Atoi(w.Args["card"])
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		t.Logf("winner: %s \t card: %s\t %s", w.Address,
// 			w.Args["card"], CARDS[i])
// 	}
// 	t.Log("Ok")
// }

// func Test_SortWinners2(t *testing.T) {
// 	var plays []types.Play

// 	perm := []int{33, 51, 28, 19, 20, 3, 44, 38, 49, 31}

// 	for i, n := range perm {
// 		p := types.Play{
// 			AceID:   "LuckyAce",
// 			GameID:  "1000",
// 			Address: fmt.Sprintf("xxx_%d", i),
// 			Card:    checkCardNum(n)}
// 		plays = append(plays, p)
// 	}

// 	winners := SortWinners(plays, 10)
// 	for _, w := range winners {
// 		i, err := strconv.Atoi(w.Args["card"])
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		t.Logf("winner: %s \t card: %s\t %s", w.Address,
// 			w.Args["card"], CARDS[i])
// 	}
// 	w := winners[0]
// 	if strings.EqualFold(w.Address, "xxx_5") {
// 		t.Logf("Ok, first is: %s", w.Address)
// 	} else {
// 		t.Errorf("winners sort error: %s", w.Address)
// 	}
// 	w = winners[1]
// 	if strings.EqualFold(w.Address, "xxx_1") {
// 		t.Logf("Ok, second is: %s", w.Address)
// 	} else {
// 		t.Errorf("winners sort error: %s", w.Address)
// 	}
// 	w = winners[9]
// 	if strings.EqualFold(w.Address, "xxx_3") {
// 		t.Logf("Ok, last is: %s", w.Address)
// 	} else {
// 		t.Errorf("winners sort error: %s", w.Address)
// 	}
// }
