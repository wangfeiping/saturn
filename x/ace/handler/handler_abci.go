package handler

import (
	"fmt"
	mrand "math/rand"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/types"
)

// BeginBlockHandle check for infraction evidence or downtime of validators
// on every begin block
func BeginBlockHandle(ctx sdk.Context, req abci.RequestBeginBlock,
	k keeper.AceKeeper) {
}

// EndBlockHandle called every block, process inflation, update validator set.
func EndBlockHandle(ctx sdk.Context, req abci.RequestEndBlock,
	k keeper.AceKeeper, bk types.BankKeeper) (vus []abci.ValidatorUpdate) {

	endGameID := CheckGameID(ctx)
	if endGameID != ctx.BlockHeight()+1-types.GameDurationHeight {
		return
	}

	ctx.Logger().Info("end game: %d at height %d\n", endGameID, ctx.BlockHeight())
	// TODO A large number of TXs processing optimizations
	plays, err := k.GetRound(ctx,
		types.CreateGameID(types.AceID, endGameID),
		types.CreateGameID(types.AceID, endGameID+1))
	if err != nil {
		fmt.Printf("query round error: %v\n", err)
		return
	}

	// TODO Random Seed function to be implemented
	err = drawCards(plays, ctx, k)
	if err != nil {
		fmt.Printf("drawing cards error: %v\n", err)
		return
	}
	// for i, p := range plays {
	// 	c := CARDS[p.Card]
	// 	fmt.Printf("draw %d play: %d; %s by %s\n", i, p.Card, c, p.Address)
	// }

	winners := checkWinner(plays)
	for _, w := range winners {
		k.SetWinner(ctx, w)
	}

	err = AwardToWinners(winners, len(plays), bk)
	if err != nil {
		fmt.Printf("award to winners error: %v\n", err)
	}
	return
}

func drawCards(plays []*types.Play, ctx sdk.Context, k keeper.AceKeeper) error {
	num := len(plays)
	if num < 51 {
		num = 51
	}
	// b := make([]byte, 51)
	// n, err := rand.Read(b)
	// if err != nil {
	// 	fmt.Println("gen crypto/rand error: ", err.Error())
	// 	return nil, err
	// }
	// fmt.Println("crypto/rand.Read：", b[:n])
	// n, err := rand.Int(rand.Reader, big.NewInt(10000000))
	// if err != nil {
	// 	fmt.Println("gen rand error: ", err.Error())
	// 	return err
	// }
	// fmt.Println("crypto/rand.Int: ", n, n.BitLen())

	// mrand.Seed(n.Int64())
	mrand.Seed(2)
	perm := mrand.Perm(num)
	fmt.Println("math/rand.Perm: ", perm)

	for i, p := range plays {
		num := perm[i]
		p.Card = checkCardNum(num)
		k.SetPlay(ctx, *p)
	}
	return nil
}

func checkCardNum(num int) int {
	if num > 51 {
		num = num % 51
	}
	return num
}

func checkWinner(plays []*types.Play) (winners []types.Winner) {
	num := len(plays)
	if num <= 1 {
		return make([]types.Winner, 0)
	}
	var winNum int
	if num < 10 {
		winNum = 1
	} else {
		winNum = num / 10
	}
	fmt.Printf("total: %d; winners: %d\n", num, winNum)
	return SortWinners(plays, winNum)
}

func SortWinners(plays []*types.Play, num int) (winners []types.Winner) {
	winners = make([]types.Winner, num)
	exist345 := isExist345(plays)
	w := types.NewWinners(plays, exist345)
	sort.Sort(w)
	for i := 0; i < num; i++ {
		winners[i] = w.GetWinner(i)
	}
	return
}

func isExist345(plays []*types.Play) bool {
	for _, p := range plays {
		if 3 < p.Card && p.Card < 16 {
			return true
		}
	}
	return false
}

// AwardToWinners launches transfers to all winners
func AwardToWinners(winners []types.Winner, playsNum int, bank types.BankKeeper) error {
	winNum := len(winners)
	if winNum == 0 {
		return nil
	}

	pooler, err := sdk.AccAddressFromBech32(types.PoolerAddress)
	if err != nil {
		fmt.Println("parse pooler's address error: " + err.Error())
		return err
	}

	if winNum == 1 {
		// only one winner
		// gold := 10
		return awardsTo(winners, playsNum, pooler, bank)
	}

	if winNum == 2 {
		// gold: 1, silver: 1 when 2 winners in all.
		// gold, silver := 6, 4
		goldAwards := playsNum * 6 / 10
		err = awardsTo(winners[0:1], goldAwards, pooler, bank)
		if err != nil {
			return err
		}
		silverAwards := playsNum - goldAwards
		return awardsTo(winners[1:], silverAwards, pooler, bank)
	}

	winNum = winNum - 1
	silverIndex := winNum / 3
	if silverIndex < 1 {
		silverIndex = 1
	}
	silverIndex++
	// gold, silver, copper := 5, 3, 2
	goldAwards := playsNum * 5 / 10
	err = awardsTo(winners[0:1], goldAwards, pooler, bank)
	if err != nil {
		return err
	}
	silverAwards := playsNum * 3 / 10
	err = awardsTo(winners[1:silverIndex], silverAwards, pooler, bank)
	if err != nil {
		return err
	}
	copperAwards := playsNum - goldAwards - silverAwards
	return awardsTo(winners[silverIndex:], copperAwards, pooler, bank)
}

func awardsTo(winners []types.Winner, chips int,
	pooler sdk.AccAddress, bank types.BankKeeper) error {
	pcs := chips * 1000000
	fmt.Printf("awards to: winners %d; awards chips %d, pcs %d\n",
		len(winners), chips, pcs)

	// allCoins := make(sdk.Coins, 1)
	// chips, err := sdk.NewCoin("chip", playsNum)
	// if err != nil {
	// 	fmt.Println("parse coins error: " + err.Error())
	// 	return nil, err
	// }
	// allCoins[0] = chips

	// goldAddr, err := sdk.AccAddressFromBech32(winners[0].Address)
	// if err != nil {
	// 	fmt.Println("parse winner's address error: " + err.Error())
	// 	return nil, err
	// }
	// err = bank.SendCoins(ctx, pooler, goldAddr, allCoins)
	// if err != nil {
	// 	fmt.Println("award to winner error: " + err.Error())
	// 	return err
	// }
	return nil
}
