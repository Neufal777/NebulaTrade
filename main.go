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

	w := wallet.Wallet{
		Name:         "naoufal boubouh",
		Balance:      230,
		Ammount:      0,
		Status:       "BUY",
		LastBuy:      0.34,
		LastSell:     0.10,
		Transactions: 0,
	}

	//w.WriteInWallet()

	for {
		trading.ExecuteMarket(w)
		time.Sleep(1 * time.Second)
	}
}
