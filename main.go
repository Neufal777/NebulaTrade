package main

import (
	"time"

	"github.com/NebulaTrade/exchanges"
	"github.com/NebulaTrade/trading"
	"github.com/NebulaTrade/wallet"
)

func main() {

	for {

		w := wallet.ReadWallet()
		w.Balance = exchanges.GetBinanceWalletBNB()
		w.WriteInWallet()
		trading.ExecuteMarket(&w)
		time.Sleep(2 * time.Second)
	}

}
