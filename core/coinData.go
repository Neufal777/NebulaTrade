package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	HALFCENT = 0.005
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

	log.Println(lastSellFloat)

	return lastSellFloat
}

//GetLastBuy - gets the price of your last buy
func GetLastBuy() float64 {

	//Get the data from the last sell
	lastBuy, _ := ioutil.ReadFile("core/lastBuy.txt")
	lastBuyString := string(lastBuy)
	lastBuyFloat, _ := strconv.ParseFloat(lastBuyString, 8)

	log.Println(lastBuyFloat)
	return lastBuyFloat
}

//DecisionMakeBuy - where the decisions of buying or selling is made
func DecisionMakeBuy() {

	//Check at how much we sold and if the actual price is lower

	updatedCoinPrice := GetLatestData("xrpeur")

	lastSellFloat := GetLastSell()
	lastPriceFloat, _ := strconv.ParseFloat(updatedCoinPrice.Last, 8)

	log.Println("Last price at:", lastPriceFloat)

	/*
		difference betwen last sell and actual
		to check if the price dropped
	*/
	difference := lastSellFloat - lastPriceFloat
	log.Println("Difference:", difference)

	if difference >= HALFCENT {

		//Execute buy order..
		log.Println("Executing buy order... ")
		log.Println("Great!, you bought this", lastSellFloat-lastPriceFloat, "cheaper ")
		log.Println()

		/*
			change last sell file with updated info
			with the lastPriceFloat
		*/

		updatedPriceBuy := []byte(fmt.Sprintf("%f", lastPriceFloat))

		err := ioutil.WriteFile("core/lastBuy.txt", updatedPriceBuy, 0)

		if err != nil {
			log.Fatal(err)
		}

		/*
			After buying new crypto, execute the function to sell them
		*/

		//change status from buying to selling
		_ = ioutil.WriteFile("core/status.txt", []byte("SELL"), 0)

	} else {
		log.Println("Waiting for the price to drop.. ")
		log.Println()

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

	log.Println("Difference to sell", differenceToSell)
	if differenceToSell >= HALFCENT {

		/*
			If the price is greater, then sell
			Execute sell order
		*/
	}

}

//ExecuteMarket - workflow of the program
func ExecuteMarket() {

	/*
		Check the status (BUY OR SELL)
	*/

	actualStatus, _ := ioutil.ReadFile("core/status.txt")
	actualStatusString := string(actualStatus)

	log.Println(actualStatusString)

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
