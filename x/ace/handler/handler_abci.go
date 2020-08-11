package handler

import (
	"math/big"
	mrand "math/rand"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/security/paillier"
	"github.com/wangfeiping/saturn/x/ace/types"
)

// BeginBlockHandle check for infraction evidence or downtime of validators
// on every begin block
func BeginBlockHandle(ctx sdk.Context, req abci.RequestBeginBlock,
	k keeper.AceKeeper) {
	ctx.Logger().With("module", "ace").
		Debug("Begin block handle", "height", ctx.BlockHeight())
}

// EndBlockHandle called every block, process inflation, update validator set.
func EndBlockHandle(ctx sdk.Context, req abci.RequestEndBlock,
	k keeper.AceKeeper, bk types.BankKeeper) (vus []abci.ValidatorUpdate) {

	ctx.Logger().With("module", "ace").
		Debug("End block handle", "height", ctx.BlockHeight())

	gameStartHeight := checkGameStartHeight(ctx)
	if !isGameOver(gameStartHeight, ctx.BlockHeight()) {
		return
	}
	gameOverHeight := ctx.BlockHeight()

	// TODO A large number of TXs processing optimizations
	plays, err := k.GetRound(ctx,
		types.CreateGameID(types.AceID, gameStartHeight),
		types.CreateGameID(types.AceID, gameOverHeight))
	if err != nil {
		ctx.Logger().With("module", "ace").
			Error("Query plays error", "error", err)
		return
	}

	if len(plays) == 0 {
		ctx.Logger().With("module", "ace").
			Debug("No game", "height", gameOverHeight)
		return
	}

	k.Set(ctx, types.KeyLastGameOverHeight, gameOverHeight)
	ctx.Logger().With("module", "ace").
		Debug(types.KeyLastGameOverHeight,
			"keeper", "set", "height", gameOverHeight)
	ctx.Logger().With("module", "ace").
		Info("Game over",
			"game", gameOverHeight, "height", ctx.BlockHeight())

	// TODO Random Seed function to be implemented
	err = drawCards(plays, ctx, k)
	if err != nil {
		ctx.Logger().With("module", "ace").
			Error("Drawing card error", "error", err)
		return
	}

	winners := checkWinner(plays)
	for _, w := range winners {
		k.SetWinner(ctx, w)
	}

	err = AwardToWinners(ctx, bk, winners, len(plays))
	if err != nil {
		ctx.Logger().With("module", "ace").
			Error("Award to winners error", "error", err)
	}
	return
}

func isGameOver(gameStartHeight, height int64) bool {
	return gameStartHeight == height+1-types.GameDurationHeight
}

func drawCards(plays []*types.Play,
	ctx sdk.Context, k keeper.AceKeeper) (err error) {
	num := len(plays)
	if num < 51 {
		num = 51
	}

	// mrand.Seed(n.Int64())
	priv := paillier.Create()
	pub := priv.PublicKey()
	seedBytes := big.NewInt(0).Bytes()
	for _, p := range plays {
		seedBytes, err = pub.Add(seedBytes, p.Seed.Hash)
	}
	seedBytes, err = priv.Decrypt(seedBytes)
	if err != nil {
		return
	}
	i := new(big.Int).SetBytes(seedBytes)
	ctx.Logger().With("module", "ace").
		Info("Open verifiable random seeds", "seed", i)

	mrand.Seed(i.Int64())
	perm := mrand.Perm(num)

	for i, p := range plays {
		num := perm[i]
		p.Card = checkCardNum(num)
		k.SetPlay(ctx, *p)
	}
	return
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
	// fmt.Printf("total: %d; winners: %d\n", num, winNum)
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
func AwardToWinners(ctx sdk.Context, bank types.BankKeeper,
	winners []types.Winner, playsNum int) error {
	winNum := len(winners)
	if winNum == 0 {
		return nil
	}

	pooler, err := sdk.AccAddressFromBech32(types.PoolerAddress)
	if err != nil {
		ctx.Logger().With("module", "ace").
			Error("parse pooler's address", "error", err.Error())
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
	// pcs := chips * 1000000
	// fmt.Printf("awards to: winners %d; awards chips %d, pcs %d\n",
	// 	len(winners), chips, pcs)

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
