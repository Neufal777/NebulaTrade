package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/ttacon/chalk"

	"github.com/NebulaTrade/utils"
)

const (
	PROFITPERTRANS = 0.0000002
)

//Wallet - contains wallet & info about it
type Wallet struct {
	Name         string
	Balance      float64
	Ammount      float64
	Transactions int
}

//GetLastSell - gets the price of your last sell
func GetLastSell() float64 {

	//Get the data from the last sell
	lastSellPrice, _ := ioutil.ReadFile("core/lastSell.txt")
	lastSellString := string(lastSellPrice)
	lastSellFloat, _ := strconv.ParseFloat(lastSellString, 8)

	return lastSellFloat
}

//GetLastBuy - gets the price of your last buy
func GetLastBuy() float64 {

	//Get the data from the last sell
	lastBuy, _ := ioutil.ReadFile("core/lastBuy.txt")
	lastBuyString := string(lastBuy)
	lastBuyFloat, _ := strconv.ParseFloat(lastBuyString, 8)

	return lastBuyFloat
}

//DecisionMakeBuy - where the decisions of buying or selling is made
func (w *Wallet) DecisionMakeBuy() {

	//Check at how much we sold and if the actual price is lower
	buyActualPrice := BitstampPrice("xrpeur")
	lastSellFloat := GetLastSell()
	lastPriceFloat := utils.StringToFloat(buyActualPrice.Last)

	/*
		difference betwen last sell and actual
		to check if the price dropped
	*/
	difference := lastSellFloat - lastPriceFloat
	fmt.Println("Waiting for price drop to buy...", chalk.Red)
	fmt.Println("LIVE Price at: "+strconv.FormatFloat(lastPriceFloat, 'f', 12, 64), chalk.Green)
	fmt.Println("LAST Sell at: "+strconv.FormatFloat(lastSellFloat, 'f', 12, 64), chalk.Green)

	/*
		Showing information about the account
		 - Ammount of crypto
		 - Actual Balance
		 - Transaction
	*/

	fmt.Println("Ammount: ", w.Ammount, chalk.Green)
	fmt.Println("Actual Balance:", w.Balance, chalk.Green)
	fmt.Println("Transactions:", w.Transactions, chalk.Green)
	fmt.Println("-----------------------------", chalk.Red)

	if difference >= PROFITPERTRANS {

		/*
			EXECUTE BUY ORDER
		*/

		/*
			change last sell file with updated info
			with the lastPriceFloat
		*/

		w.Ammount = w.Balance / lastPriceFloat

		/*
			Displaying information
		*/

		lastPriceStringToWrite := utils.FloatToString(lastPriceFloat)
		utils.WriteFile(lastPriceStringToWrite, "core/lastBuy.txt")
		/*
			After buying new crypto, execute the function to sell them
		*/

		//change status from buying to selling
		_ = ioutil.WriteFile("core/status.txt", []byte("SELL"), 0)

		w.Transactions++

		fmt.Println("Buy order executed!", chalk.Green)

		/*
			Details about the wallet
		*/
		fmt.Println("Ammount: ", w.Ammount, chalk.Green)
		fmt.Println("Actual Balance:", w.Balance, chalk.Green)
		fmt.Println("Transactions:", w.Transactions, "\n", chalk.Green)
		fmt.Println("-----------------------------", chalk.Red)

	} else {

		// fmt.Println("Actual Balance:", w.Balance, chalk.Green)
		// fmt.Println("Transactions:", w.Transactions, chalk.Green)
		// fmt.Println("Actual Ammount:", w.Ammount, chalk.Green)
		// fmt.Println("Waiting for the price to drop.. ", chalk.Green)
		// fmt.Println(, chalk.Green)

	}
}

//DecisionMakeSell - once we've bought new crypto, we wait to sell them
func (w *Wallet) DecisionMakeSell() {

	/*
		Get latest data from the coin
		To evaluate when to sell the coin
		and change status from buying to selling
	*/

	data := BitstampPrice("xrpeur")
	//data := BinancePrice("TRXBTC")

	currentPriceFloat, _ := strconv.ParseFloat(data.Last, 16)
	/*
		Information about the last BUY
	*/
	lastBuyFloat := GetLastBuy()

	differenceToSell := currentPriceFloat - lastBuyFloat

	/*
		Displaying information in the console
	*/
	fmt.Println("Waiting to sell, price now at:", data.Last, chalk.Green)
	fmt.Println("Difference to sell:", differenceToSell, chalk.Green)
	fmt.Println("LAST BUY:", lastBuyFloat, chalk.Green)
	fmt.Println("Ammount:", w.Ammount, chalk.Green)
	fmt.Println("Balance:", w.Balance, chalk.Green)
	fmt.Println("Transactions:", w.Transactions, chalk.Green)
	fmt.Println("-----------------------------", chalk.Red)

	if differenceToSell >= PROFITPERTRANS {

		/*
			If the price is greater, then sell
			Execute sell order
		*/

		/*
			Change 2 files:
				- Last Sell
				- Ststus to BUY
		*/
		w.Balance = w.Ammount * currentPriceFloat
		w.Ammount = 0

		lastPriceStringToWrite := utils.FloatToString(currentPriceFloat)

		err := ioutil.WriteFile("core/lastSell.txt",
			[]byte(lastPriceStringToWrite), 0)

		_ = ioutil.WriteFile("core/status.txt",
			[]byte("BUY"), 0)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("SOLD!", chalk.Green)
		w.Transactions++
	}

}

//ExecuteMarket - workflow of the program
func (w *Wallet) ExecuteMarket() {

	/*
		Check the status (BUY OR SELL)
	*/

	actualStatusString := utils.ReadFile("core/status.txt")

	/*
		depending on the status,
		we execute buy or sell orders
	*/
	switch actualStatusString {
	case "BUY":
		w.DecisionMakeBuy()
	case "SELL":
		w.DecisionMakeSell()
	}

}
