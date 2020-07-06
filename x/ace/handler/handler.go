package handler

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mrand "math/rand"
	"sort"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/wangfeiping/saturn/x/ace/keeper"
	"github.com/wangfeiping/saturn/x/ace/types"
)

var CARDS = [...]string{
	"2 diamond", "2 club", "2 heart", "2 spade",
	"3 diamond", "3 club", "3 heart", "3 spade",
	"4 diamond", "4 club", "4 heart", "4 spade",
	"5 diamond", "5 club", "5 heart", "5 spade",
	"6 diamond", "6 club", "6 heart", "6 spade",
	"7 diamond", "7 club", "7 heart", "7 spade",
	"8 diamond", "8 club", "8 heart", "8 spade",
	"9 diamond", "9 club", "9 heart", "9 spade",
	"10 diamond", "10 club", "10 heart", "10 spade",
	"J diamond", "J club", "J heart", "J spade",
	"Q diamond", "Q club", "Q heart", "Q spade",
	"K diamond", "K club", "K heart", "K spade",
	"A diamond", "A club", "A heart", "A spade",
}

const DESC = "The %s of %s"

// NewHandler creates an sdk.Handler for all the ace type messages
func NewHandler(k keeper.Keeper, bank types.BankKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgPlay:
			fmt.Println("handle play msg")
			return handleMsgPlay(ctx, k, bank, msg)
		case types.MsgAce:
			fmt.Printf("handle ace msg: %s %s %s\n", msg.AceID, msg.GameID, msg.Action)
			// return nil, fmt.Errorf("not support")
			return handleMsgAce(ctx, k, bank, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgAce handle request
func handleMsgAce(ctx sdk.Context, k keeper.Keeper, bank types.BankKeeper,
	m types.MsgAce) (*sdk.Result, error) {
	if !strings.EqualFold("end", m.Action) {
		return nil, fmt.Errorf("not support action: %s", m.Action)
	}
	h, err := strconv.ParseInt(m.GameID, 10, 64)
	if err != nil {
		fmt.Printf("wrong game id: %s %v\n", m.GameID, err)
		return nil, err
	}
	if h >= ctx.BlockHeight()-types.GameDurationHeight {
		err = fmt.Errorf("game is not over: %d current: %d", h, ctx.BlockHeight())
		fmt.Println(err.Error())
		return nil, err
	}

	// TODO A large number of TXs processing optimizations
	plays, err := k.GetRound(ctx,
		fmt.Sprintf("%s-%d", types.AceID, h), fmt.Sprintf("%s-%d", types.AceID, h+1))
	if err != nil {
		fmt.Printf("query round error: %v\n", err)
		return nil, err
	}

	// TODO Random Seed function to be implemented
	err = drawCards(plays, ctx, k)
	if err != nil {
		fmt.Printf("drawing cards error: %v\n", err)
		return nil, err
	}

	winners := checkWinner(plays)
	for _, w := range winners {
		fmt.Printf("winner: %s\n", w.Address)
	}

	err = AwardToWinners(winners, len(plays), bank)
	if err != nil {
		fmt.Printf("award to winners error: %v\n", err)
		return nil, err
	}

	// Define msg-ace events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule,
				types.AttributeValueCategory),
			sdk.NewAttribute("ace_id", m.AceID),
			sdk.NewAttribute("game_id", m.GameID),
			sdk.NewAttribute("action", m.Action),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgPlay handle request
func handleMsgPlay(ctx sdk.Context, k keeper.Keeper, bank types.BankKeeper,
	m types.MsgPlay) (*sdk.Result, error) {

	h, err := strconv.ParseInt(m.GameID, 10, 64)
	if err != nil {
		fmt.Printf("wrong game id: %s %v\n", m.GameID, err)
		return nil, err
	}
	if h <= ctx.BlockHeight()-types.GameDurationHeight {
		err = fmt.Errorf("game is over: %d current: %d", h, ctx.BlockHeight())
		fmt.Println(err.Error())
		return nil, err
	}
	// TODO How to get the TX's hash, Resolving conflicts over play TXs
	play := types.Play{
		// TxHash:  ,
		AceID:   m.AceID,
		GameID:  m.GameID,
		RoundID: m.RoundID,
		Address: m.Address.String(),
		Seed:    m.Seed,
		Func:    m.Func,
		Args:    m.Args}
	k.SetPlay(ctx, fmt.Sprintf("%s-%s-%s:%s",
		play.AceID, play.GameID, play.RoundID, play.Address), play)
	// args := strings.Split(play.Args, ",")
	coins, err := sdk.ParseCoins("1chip")
	if err != nil {
		fmt.Println("parse coins error: " + err.Error())
		return nil, err
	}
	pooler, err := sdk.AccAddressFromBech32(types.PoolerAddress)
	if err != nil {
		fmt.Println("parse pooler's address error: " + err.Error())
		return nil, err
	}
	err = bank.SendCoins(ctx, m.Address, pooler, coins)
	if err != nil {
		fmt.Println("send coins error: " + err.Error())
		return nil, err
	}
	fmt.Printf("handle play msg: %d - %s %s %s\n", ctx.BlockHeight(),
		m.AceID, m.Func, m.Args)

	// Define msg-play events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule,
				types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender,
				m.Address.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
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

func drawCards(plays []types.Play, ctx sdk.Context, k keeper.Keeper) error {
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
	n, err := rand.Int(rand.Reader, big.NewInt(10000000))
	if err != nil {
		fmt.Println("gen rand error: ", err.Error())
		return err
	}
	fmt.Println("crypto/rand.Int: ", n, n.BitLen())

	mrand.Seed(n.Int64())
	perm := mrand.Perm(num)
	fmt.Println("math/rand.Perm: ", perm)

	for i, p := range plays {
		num := perm[i]
		p.Card = checkCardNum(num)
		c := CARDS[p.Card]
		k.SetPlay(ctx, fmt.Sprintf("%s-%s-%s:%s",
			p.AceID, p.GameID, p.RoundID, p.Address), p)
		fmt.Printf("%d play: %s by %s\n", i, c, p.Address)
	}
	return nil
}

func checkWinner(plays []types.Play) (winners []types.Winner) {
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

func SortWinners(plays []types.Play, num int) (winners []types.Winner) {
	winners = make([]types.Winner, num)
	exist345 := isExist345(plays)
	w := types.NewWinners(plays, exist345)
	sort.Sort(w)
	for i := 0; i < num; i++ {
		winners[i] = w.GetWinner(i)
	}
	return
}

func checkCardNum(num int) int {
	if num > 51 {
		num = num % 51
	}
	return num
}

func isExist345(plays []types.Play) bool {
	for _, p := range plays {
		if 3 < p.Card && p.Card < 16 {
			return true
		}
	}
	return false
}
