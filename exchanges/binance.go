package exchanges

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/NebulaTrade/binanceaccount"
	"github.com/NebulaTrade/utils"
	"github.com/NebulaTrade/wallet"
	"github.com/adshao/go-binance/v2"
	"github.com/ttacon/chalk"
)

var (
	apiKey    = binanceaccount.APIKEY
	secretKey = binanceaccount.APISECRET
)

const (
	//MITHBNB exchange from binance
	MITHBNB = "MITHBNB"
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

//GetBinanceWalletBNB -
func GetBinanceWalletBNB() float64 {

	//futuresClient := binance.NewFuturesClient(apiKey, secretKey)   // USDT-M Futures
	//deliveryClient := binance.NewDeliveryClient(apiKey, secretKey) // Coin-M Futures
	client := binance.NewClient(apiKey, secretKey)

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)

	}
	binanceWalletFloat := utils.StringToFloat(res.Balances[4].Free)

	return binanceWalletFloat
}

//ExecuteBuyOrderMITHBNB -
func ExecuteBuyOrderMITHBNB(ammountToBuy string, priceToBuy string, w *wallet.Wallet) {
	client := binance.NewClient(apiKey, secretKey)

	_, err := client.NewCreateOrderService().Symbol("MITHBNB").
		Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(ammountToBuy).
		Price(priceToBuy).Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	/*
		If the order is made successfuly update the wallet
	*/

	w.Status = "SELL"
	w.Balance = GetBinanceWalletBNB()
	w.Transactions++
	w.LastBuy = utils.StringToFloat(priceToBuy)
	//w.Ammount = w.Balance / lastPriceFloat
	w.WriteInWallet()
	fmt.Println(chalk.Bold.TextStyle("BUY ORDER EXECUTED!"), chalk.Green)
}
