package main

import (
	"time"

	"github.com/NebulaTrade/trading"
	"github.com/NebulaTrade/wallet"
)

const (
	XRP = "xrpeur"
	BTC = "btceur"
)

func main() {

	// w := trading.SimulationWallet{
	// 	Name:         "Naoufal dahouli",
	// 	Balance:      250,
	// 	Ammount:      0,
	// 	Transactions: 0,
	// }

	// for x := 1; x < 14400; x++ {

	// 	w.ExecuteMarket()
	// 	time.Sleep(2 * time.Second)
	// }

	//w.WriteInWallet()

	for {
		w := wallet.ReadWallet()
		trading.ExecuteMarket(&w)
		time.Sleep(3 * time.Second)
	}
}
