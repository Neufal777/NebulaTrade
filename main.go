package main

import "github.com/NebulaTrade/core"

const (
	XRP = "xrpeur"
	BTC = "btceur"
)

func main() {

	// w := core.Wallet{
	// 	Name:         "Naoufal dahouli",
	// 	Balance:      100,
	// 	Ammount:      0,
	// 	Transactions: 0,
	// }
	// for x := 1; x < 14400; x++ {

	// 	w.ExecuteMarket()
	// 	time.Sleep(2 * time.Second)
	// }

	core.BitfinexPrice("tXRPUSD")
}
