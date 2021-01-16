package main

import (
	"flag"
	"time"

	"github.com/NebulaTrade/config"
	"github.com/NebulaTrade/exchanges"
	"github.com/NebulaTrade/trading"
	"github.com/NebulaTrade/wallet"
)

func main() {

	/*
		Geeting all the flags and setting them to the config file
	*/

	BinanceWallet := exchanges.GetBinanceWalletBNB()

	configCurrency := flag.String("currency", "1234", "Currency")
	configProfit := flag.Float64("profit", 0.00002, "Profit to make")

	flag.Parse()

	config.CURRENCY = *configCurrency
	config.PROFIT = *configProfit

	/*
		Set all the initial config data to start initial buy
	*/
	w := wallet.ReadWallet()
	w.Symbol = config.CURRENCY
	w.LastSell = 5000
	w.Transactions = 0
	w.Status = "BUY"
	w.Available = BinanceWallet / 3
	w.Balance = BinanceWallet
	w.Profit = config.PROFIT
	w.WriteInWallet()

	for {

		trading.SellingPositions()
		w := wallet.ReadWallet()

		if w.Status == "BUY MORE" || w.Status == "BUY" {
			w.Timer++
		}
		w.WriteInWallet()
		trading.ExecuteMarket(&w)
		time.Sleep(2 * time.Second)
	}

}
