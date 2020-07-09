package handler_test

import (
	"fmt"
	mrand "math/rand"

	// "strconv"
	// "strings"
	// "testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/auth"
	acehandler "github.com/wangfeiping/saturn/x/ace/handler"
	acekeeper "github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/types"
)

var _ = Describe("AceHandler", func() {
	var (
		perm []int
		num  int   = 10000
		seed int64 = 1

		ctx           sdk.Context
		keeper        acekeeper.AceKeeper
		accountKeeper auth.AccountKeeper
		handler       sdk.Handler
	)

	addr := sdk.AccAddress([]byte("addr.test.1"))

	ctx = CreateMockSdkContext()
	keeper = CreateMockAceKeeper()
	pk := CreateMockParamsKeeper()
	accountKeeper = CreateMockAccountKeeper(pk)
	bankKeeper := CreateMockBankKeeper(accountKeeper, pk)

	handler = acehandler.NewHandler(keeper, bankKeeper)

	BeforeEach(func() {
		mrand.Seed(seed)
		perm = mrand.Perm(num)

		acc := accountKeeper.NewAccountWithAddress(ctx, addr)
		accountKeeper.SetAccount(ctx, acc)
		bankKeeper.SetCoins(ctx, addr, sdk.NewCoins(sdk.NewInt64Coin("chip", 10000)))

	})

	Describe("Create an ace handler", func() {
		Context(fmt.Sprintf("with %d random int number", num), func() {
			It("should be success", func() {
				Expect(len(perm)).To(Equal(num))
				Expect(handler).NotTo(BeZero())
			})
		})
	})

	Describe("Call handler", func() {
		Context(fmt.Sprintf("with %d play messages", num), func() {
			It("should be success", func() {
				seed := types.Seed{Hash: []byte("")}
				msg, err := types.NewMsgPlay(
					"LuckyAce", "0", "",
					seed, "draw", "",
					addr)
				_, err = handler(ctx, *msg)
				Expect(err).ShouldNot(HaveOccurred())
				// coins := bankKeeper.GetCoins(ctx, addr)
				// fmt.Printf("account coin: %d %s\n", coins[0].Amount.Int64(), coins[0].Denom)
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
