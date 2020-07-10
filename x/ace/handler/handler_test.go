package handler_test

import (
	"fmt"
	"time"

	// "strconv"
	// "strings"
	// "testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"

	acehandler "github.com/wangfeiping/saturn/x/ace/handler"
	acekeeper "github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/types"
)

var _ = Describe("AceHandler", func() {
	var (
		denom  string = "chip"
		amount int64  = 100
		num    int    = 100

		addrs         []sdk.AccAddress
		ctx           sdk.Context
		keeper        acekeeper.AceKeeper
		accountKeeper auth.AccountKeeper
		handler       sdk.Handler
	)

	ctx = CreateMockSdkContext()
	keeper = CreateMockAceKeeper()
	pk := CreateMockParamsKeeper()
	accountKeeper = CreateMockAccountKeeper(pk)
	bankKeeper := CreateMockBankKeeper(accountKeeper, pk)

	handler = acehandler.NewHandler(keeper, bankKeeper)

	for i := 0; i < num; i++ {
		addrs = append(addrs, sdk.AccAddress([]byte(
			fmt.Sprintf("addr.test.%d", i))))
		acc := accountKeeper.NewAccountWithAddress(ctx, addrs[i])
		accountKeeper.SetAccount(ctx, acc)
		bankKeeper.SetCoins(ctx, addrs[i],
			sdk.NewCoins(sdk.NewInt64Coin(denom, amount)))
	}

	BeforeEach(func() {

	})

	Describe("Create an ace handler", func() {
		Context("with mock keepers", func() {
			It("should be success", func() {
				Expect(handler).NotTo(BeZero())
			})
		})
	})

	Describe("Call ace handler", func() {
		Context(fmt.Sprintf("with %d play messages", num), func() {
			It("should be success", func() {
				for i := 0; i < num; i++ {
					seed := types.Seed{Hash: []byte("0")}
					msg := types.NewMsgPlay(
						"LuckyAce", "0", "",
						seed, "draw", "", addrs[i])
					_, err := handler(ctx, *msg)
					Expect(err).ShouldNot(HaveOccurred())
					// coins := bankKeeper.GetCoins(ctx, addrs[i])
					// Expect(coins[0].Amount.Int64()).To(Equal(amount - 1))
				}
			})
		})

		Context("Check the balance of all accounts", func() {
			It("should be success", func() {
				for i := 0; i < num; i++ {
					coins := bankKeeper.GetCoins(ctx, addrs[i])
					Expect(coins[0].Amount.Int64()).To(Equal(amount - 1))
				}
			})
		})

		Context("end the game", func() {
			It("should be success", func() {
				// Updates the block height
				ctx = ctx.WithBlockHeader(
					abci.Header{Height: 11, Time: time.Unix(10, 0)})
				msg := types.NewMsgAce(
					"LuckyAce", "0", "end", addrs[0])
				_, err := handler(ctx, *msg)
				Expect(err).ShouldNot(HaveOccurred())
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
