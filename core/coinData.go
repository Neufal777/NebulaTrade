package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"

	"github.com/NebulaTrade/utils"
	"github.com/ttacon/chalk"
)

const (
	HALFCENT = 0.0000002
	CENT     = 0.01
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
	fmt.Println(chalk.Green, "Waiting to buy..")

	updatedCoinPrice := BitstampPrice("xrpeur")
	//updatedCoinPrice := BinancePrice("XRPBTC")

	lastSellFloat := GetLastSell()
	lastPriceFloat, _ := strconv.ParseFloat(updatedCoinPrice.Last, 16)

	log.Println("Last price at:", lastPriceFloat)
	log.Println("Last sell at:", lastSellFloat)

	/*
		difference betwen last sell and actual
		to check if the price dropped
	*/
	difference := lastSellFloat - lastPriceFloat
	log.Println("Difference:", difference)

	if difference >= HALFCENT {

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

		log.Println(reflect.TypeOf(lastPriceFloat))

		lastPriceStringToWrite := utils.FloatToString(lastPriceFloat)
		utils.WriteFile(lastPriceStringToWrite, "core/lastBuy.txt")
		/*
			After buying new crypto, execute the function to sell them
		*/

		//change status from buying to selling
		_ = ioutil.WriteFile("core/status.txt", []byte("SELL"), 0)

		w.Transactions++
		log.Println("Great!, you bought this", lastSellFloat-lastPriceFloat, "cheaper ")
		log.Println("Ammount: ", w.Ammount)
		log.Println("Actual Balance:", w.Balance)
		log.Println("Transactions:", w.Transactions)
		log.Println()

	} else {
		log.Println("Actual Balance:", w.Balance)
		log.Println("Transactions:", w.Transactions)
		log.Println("Actual Ammount:", w.Ammount)
		log.Println("Waiting for the price to drop.. ")
		log.Println()

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
	log.Println("Waiting to sell, price now at:", data.Last)

	log.Print(reflect.TypeOf(data.Last))
	currentPriceFloat, _ := strconv.ParseFloat(data.Last, 16)
	/*
		Information about the last BUY
	*/
	lastBuyFloat := GetLastBuy()

	log.Println("LAST BUY:", lastBuyFloat)

	differenceToSell := currentPriceFloat - lastBuyFloat

	log.Println("Difference to sell:", differenceToSell)
	if differenceToSell >= HALFCENT {

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

		log.Println("SOLD!")
		w.Transactions++
	}

	log.Println("Actual Ammount:", w.Ammount)
	log.Println("Actual Balance:", w.Balance)
	log.Println("Transactions:", w.Transactions)
	log.Println()

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
