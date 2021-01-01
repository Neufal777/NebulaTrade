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
	XRPBNB  = "XRPBNB"
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
	opened := CheckOpenOrdersBinance()

	switch opened {
	case 0:

		//If the order is closed, update the file with the information
		w.Status = "SELL"
		w.Balance = GetBinanceWalletBNB()
		w.LastBuy = utils.StringToFloat(priceToBuy)
		w.Ammount = utils.StringToFloat(ammountToBuy)
		w.Balance = GetBinanceWalletBNB()
		w.Transactions++

		//Write the update information in our register
		w.WriteInWallet()

		fmt.Println(chalk.Bold.TextStyle("BUY ORDER EXECUTED!"), chalk.Green)

	case 1:
		w.Status = "ORDER"
		w.WriteInWallet()
	}

}

//ExecuteSellOrderMITHBNB -
func ExecuteSellOrderMITHBNB(ammountToSell string, priceToSell string, w *wallet.Wallet) {

	client := binance.NewClient(apiKey, secretKey)

	_, err := client.NewCreateOrderService().Symbol("MITHBNB").
		Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(ammountToSell).
		Price(priceToSell).Do(context.Background())

	log.Println("PRICE TO SELL", priceToSell)
	log.Println("AMMOUNT TO SELL", ammountToSell)
	if err != nil {
		fmt.Println(err)
		return
	}

	openOrders := CheckOpenOrdersBinance()

	switch openOrders {
	case 0:

		/*
			If the order is made successfuly update the wallet
		*/

		w.Status = "BUY"
		w.Balance = GetBinanceWalletBNB()
		w.Transactions++
		w.LastSell = utils.StringToFloat(priceToSell)

		w.WriteInWallet()
		fmt.Println(chalk.Bold.TextStyle("SELL ORDER EXECUTED!"), chalk.Green)

	case 1:
		w.Status = "ORDER"
		w.WriteInWallet()
	}
}

//CheckOpenOrdersBinance -
func CheckOpenOrdersBinance() int {

	client := binance.NewClient(apiKey, secretKey)

	//check open orders from the user
	openOrders, err := client.NewListOpenOrdersService().Symbol("MITHBNB").
		Do(context.Background())
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(openOrders)

	for _, o := range openOrders {

		//Show active orders
		log.Println("ORDER Id:", o.OrderID)
		log.Println("ORDER Status:", o.Side)
		log.Println("ORDER Price:", o.Price)
	}

	if len(openOrders) >= 1 {
		return 1
	}

	return 0
}
