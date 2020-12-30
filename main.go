package main

import (
	"time"

	"github.com/NebulaTrade/trading"
	"github.com/NebulaTrade/wallet"
)

func main() {

	for {
		w := wallet.ReadWallet()
		trading.ExecuteMarket(&w)
		time.Sleep(1 * time.Second)
	}

}
