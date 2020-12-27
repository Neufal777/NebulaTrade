package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/NebulaTrade/utils"
	"github.com/ttacon/chalk"
)

const (
	HALFCENT = 0.000005
	CENT     = 0.01
)

//Coin - latest info about the coin
type Coin struct {
	High      string `json:"high"`
	Last      string `json:"last"`
	Timestamp string `json:"timestamp"`
	Bid       string `json:"bid"`
	Vwap      string `json:"vwap"`
	Volume    string `json:"volume"`
	Low       string `json:"low"`
	Ask       string `json:"ask"`
	Open      string `json:"open"`
}

//GetLatestData - get latest data of X coin
func GetLatestData(exchange string) Coin {

	var CoinData Coin
	resp, err := http.Get("https://www.bitstamp.net/api/v2/ticker/" + exchange)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)

	}
	err = json.Unmarshal(data, &CoinData)

	if err != nil {
		log.Fatal("Decoding error: ", err)
	}

	//time.Sleep(1 * time.Second)
	//DecisionMakeBuy(&CoinData)
	return CoinData

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
func DecisionMakeBuy() {

	//Check at how much we sold and if the actual price is lower
	fmt.Println(chalk.Green, "Waiting to buy..")

	updatedCoinPrice := GetLatestData("xrpeur")

	lastSellFloat := GetLastSell()
	lastPriceFloat, _ := strconv.ParseFloat(updatedCoinPrice.Last, 8)

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
		log.Println("Great!, you bought this", lastSellFloat-lastPriceFloat, "cheaper \n")

		/*
			change last sell file with updated info
			with the lastPriceFloat
		*/

		utils.WriteFile(lastPriceFloat, "core/lastBuy.txt")
		/*
			After buying new crypto, execute the function to sell them
		*/

		//change status from buying to selling
		_ = ioutil.WriteFile("core/status.txt", []byte("SELL"), 0)

	} else {
		log.Println("Waiting for the price to drop.. \n")

	}
}

//DecisionMakeSell - once we've bought new crypto, we wait to sell them
func DecisionMakeSell() {

	/*
		Get latest data from the coin
		To evaluate when to sell the coin
		and change status from buying to selling
	*/

	data := GetLatestData("xrpeur")
	log.Println("Waiting to sell, price now at:", data.Last)

	currentPriceFloat, _ := strconv.ParseFloat(data.Last, 8)
	/*
		Information about the last BUY
	*/
	lastBuyFloat := GetLastBuy()

	log.Println("LAST BUY:", lastBuyFloat)

	differenceToSell := currentPriceFloat - lastBuyFloat

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

		err := ioutil.WriteFile("core/lastSell.txt",
			[]byte(fmt.Sprintf("%f", currentPriceFloat)), 0)

		_ = ioutil.WriteFile("core/status.txt",
			[]byte("BUY"), 0)

		if err != nil {
			log.Fatal(err)
		}

	}

}

//ExecuteMarket - workflow of the program
func ExecuteMarket() {

	/*
		Check the status (BUY OR SELL)
	*/

	actualStatusString := utils.ReadFile("core/status.txt")

	log.Println("status!", actualStatusString)
	/*
		depending on the status,
		we execute buy or sell orders
	*/
	switch actualStatusString {
	case "BUY":
		DecisionMakeBuy()
	case "SELL":
		DecisionMakeSell()
	}

}
