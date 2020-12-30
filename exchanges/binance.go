package exchanges

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//BinanceCoin -
type BinanceCoin struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

//BinancePrice - get latest data of X coin
func BinancePrice(exchange string) BinanceCoin {

	var CoinData BinanceCoin
	resp, err := http.Get("https://api.binance.com/api/v1/ticker/price?symbol=" + exchange)

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
