package main

import (
	"fmt"

	db "github.com/tendermint/tm-db"
)

func main() {
	wr1, err := db.NewGoLevelDB("testing", "./test-db/")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	defer wr1.Close()
	// err = wr1.Set([]byte("key1"), []byte("value2"))
	// if err != nil {
	// 	fmt.Println("err: ", err)
	//  return
	// }
	v, err := wr1.Get([]byte("key1"))
	if err != nil {
		fmt.Println("err: ", err)
		return
	}
	fmt.Println("value: ", string(v))
}
