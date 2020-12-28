package exchanges

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//BitstampCoin - latest info about the coin
type BitstampCoin struct {
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

//BitstampPrice - get latest data of X coin
func BitstampPrice(exchange string) BitstampCoin {

	var CoinData BitstampCoin
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

	return CoinData

}
