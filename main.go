package main

import (
	"log"

	"github.com/NebulaTrade/exchanges"
)

func main() {

	// for {

	// 	w := wallet.ReadWallet()
	// 	w.Balance = exchanges.GetBinanceWalletBNB()
	// 	w.WriteInWallet()
	// 	trading.ExecuteMarket(&w)
	// 	time.Sleep(2 * time.Second)
	// }

	openOrders := exchanges.CheckOpenOrdersBinance()

	if openOrders == 1 {

		log.Println("OPENED ORDERS")

	} else if openOrders == 0 {
		log.Println("NO OPEN ORDERS")
	}

}
