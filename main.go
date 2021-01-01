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
		time.Sleep(4 * time.Second)
	}

	//fmt.Println(exchanges.GetBinanceWalletBNB())
	//log.Println(exchanges.GetBinanceWalletBNB() - 0.01646)
}
