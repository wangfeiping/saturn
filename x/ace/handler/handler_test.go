package handler

import (
	"fmt"
	mrand "math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/wangfeiping/saturn/x/ace/types"
)

func Test_SortWinners(t *testing.T) {
	var plays []types.Play

	num := 100
	mrand.Seed(1)
	perm := mrand.Perm(num)
	fmt.Println("math/rand.Perm: ", perm)

	for i, n := range perm {
		p := types.Play{
			AceID:   "LuckyAce",
			GameID:  "1000",
			Address: fmt.Sprintf("xxx_%d", i),
			Card:    checkCardNum(n)}
		plays = append(plays, p)
	}

	winners := SortWinners(plays, 33)
	for _, w := range winners {
		i, err := strconv.Atoi(w.Args["card"])
		if err != nil {
			t.Error(err)
		}
		t.Logf("winner: %s \t card: %s\t %s", w.Address,
			w.Args["card"], CARDS[i])
	}
	t.Log("Ok")
}

func Test_SortWinners2(t *testing.T) {
	var plays []types.Play

	perm := []int{33, 51, 28, 19, 20, 3, 44, 38, 49, 31}

	for i, n := range perm {
		p := types.Play{
			AceID:   "LuckyAce",
			GameID:  "1000",
			Address: fmt.Sprintf("xxx_%d", i),
			Card:    checkCardNum(n)}
		plays = append(plays, p)
	}

	winners := SortWinners(plays, 10)
	for _, w := range winners {
		i, err := strconv.Atoi(w.Args["card"])
		if err != nil {
			t.Error(err)
		}
		t.Logf("winner: %s \t card: %s\t %s", w.Address,
			w.Args["card"], CARDS[i])
	}
	w := winners[0]
	if strings.EqualFold(w.Address, "xxx_5") {
		t.Logf("Ok, first is: %s", w.Address)
	} else {
		t.Errorf("winners sort error: %s", w.Address)
	}
	w = winners[1]
	if strings.EqualFold(w.Address, "xxx_1") {
		t.Logf("Ok, second is: %s", w.Address)
	} else {
		t.Errorf("winners sort error: %s", w.Address)
	}
	w = winners[9]
	if strings.EqualFold(w.Address, "xxx_3") {
		t.Logf("Ok, last is: %s", w.Address)
	} else {
		t.Errorf("winners sort error: %s", w.Address)
	}
}
