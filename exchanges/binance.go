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

//Order - defines struct of an order
type Order struct {
	Symbol        string
	OrderID       int64
	ClientOrderID string
	Price         string
	Ammount       string
	Status        string //BUY OR SELL ORDER
	Type          string //LIMIT, STOP, ETC
}

const (
	//MITHBNB exchange from binance
	MITHBNB = "MITHBNB"
	//XRPBNB  = "XRPBNB"
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

	_, err := client.NewCreateOrderService().Symbol(MITHBNB).
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
	opened, allOrders := CheckOpenOrdersBinance()

	switch opened {
	case 0:

		//If there is no open orders
		w.Status = "SELL"
		w.Balance = GetBinanceWalletBNB()
		w.LastBuy = utils.StringToFloat(priceToBuy)
		w.Ammount = utils.StringToFloat(ammountToBuy)
		w.Balance = GetBinanceWalletBNB()
		w.Timer = 0
		w.Transactions++

		//Write the update information in our register
		w.WriteInWallet()

		fmt.Println(chalk.Bold.TextStyle("BUY ORDER EXECUTED!"), chalk.Green)

	case 1:

		orders := AllOpenOrdersBinance(allOrders)

		for _, ord := range orders {

			if ord.Status == "BUY" && ord.Symbol == MITHBNB {

				//That means we still have open order for buying
				w.Status = "BUY ORDER"
				w.Balance = GetBinanceWalletBNB()
				w.LastBuy = utils.StringToFloat(priceToBuy)
				w.Ammount = utils.StringToFloat(ammountToBuy)
				w.Timer = 0
				w.Transactions++
				w.WriteInWallet()

			} else {

				//If there is no open orders for specific Coin and status
				w.Status = "SELL"
				w.Balance = GetBinanceWalletBNB()
				w.LastBuy = utils.StringToFloat(priceToBuy)
				w.Ammount = utils.StringToFloat(ammountToBuy)
				w.Balance = GetBinanceWalletBNB()
				w.Timer = 0
				w.Transactions++

				//Write the update information in our register
				w.WriteInWallet()

				fmt.Println(chalk.Bold.TextStyle("BUY ORDER EXECUTED!"), chalk.Green)
			}
		}

	}

}

//ExecuteSellOrderMITHBNB -
func ExecuteSellOrderMITHBNB(ammountToSell string, priceToSell string, w *wallet.Wallet) {

	client := binance.NewClient(apiKey, secretKey)

	_, err := client.NewCreateOrderService().Symbol(MITHBNB).
		Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(ammountToSell).
		Price(priceToSell).Do(context.Background())

	log.Println("PRICE TO SELL", priceToSell)
	log.Println("AMMOUNT TO SELL", ammountToSell)

	if err != nil {
		fmt.Println(err)
		return
	}

	openOrders, allOrders := CheckOpenOrdersBinance()

	switch openOrders {
	case 0:

		/*
			If the order is made successfuly update the wallet
		*/

		w.Status = "BUY"
		w.Balance = GetBinanceWalletBNB()
		w.Transactions++
		w.LastSell = utils.StringToFloat(priceToSell)
		w.Timer = 0
		w.WriteInWallet()

		fmt.Println(chalk.Bold.TextStyle("SELL ORDER EXECUTED!"), chalk.Green)

	case 1:

		orders := AllOpenOrdersBinance(allOrders)

		for _, ord := range orders {

			if ord.Status == "SELL" && ord.Symbol == MITHBNB {

				//That means we still have open order for selling
				w.Status = "SELL ORDER"
				w.Balance = GetBinanceWalletBNB()
				w.LastSell = utils.StringToFloat(priceToSell)
				w.Ammount = utils.StringToFloat(ammountToSell)
				w.Timer = 0
				w.Transactions++
				w.WriteInWallet()

			} else {

				//If there is no open orders for specific Coin and status
				w.Status = "BUY"
				w.Balance = GetBinanceWalletBNB()
				w.LastSell = utils.StringToFloat(priceToSell)
				w.Ammount = utils.StringToFloat(ammountToSell)
				w.Balance = GetBinanceWalletBNB()
				w.Timer = 0
				w.Transactions++

				//Write the update information in our register
				w.WriteInWallet()

				fmt.Println(chalk.Bold.TextStyle("SELL ORDER EXECUTED!"), chalk.Green)
			}
		}

	}
}

//CheckOpenOrdersBinance -
func CheckOpenOrdersBinance() (int, []*binance.Order) {

	client := binance.NewClient(apiKey, secretKey)

	//check open orders from the user
	openOrders, err := client.NewListOpenOrdersService().Symbol(MITHBNB).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)

	}

	var ordersID []int64
	var orderTypes []string

	for _, o := range openOrders {

		//Show active orders
		log.Println("ORDER Id:", o.OrderID)
		log.Println("ORDER Status:", o.Side)
		log.Println("ORDER Price:", o.Price)
		log.Println("-----------------------")
		ordersID = append(ordersID, o.OrderID)
		orderType := utils.AnyTypeToString(o.Side)
		orderTypes = append(orderTypes, orderType)
	}

	if len(openOrders) >= 1 {
		return 1, openOrders
	}

	return 0, nil
}

//AllOpenOrdersBinance - show all open orders in binance
func AllOpenOrdersBinance(allOrders []*binance.Order) []Order {

	var ActiveOrders []Order

	for _, o := range allOrders {

		order := Order{
			Symbol:        o.Symbol,
			OrderID:       o.OrderID,
			ClientOrderID: o.ClientOrderID,
			Price:         o.Price,
			Ammount:       o.OrigQuantity,
			Status:        utils.AnyTypeToString(o.Side),
			Type:          utils.AnyTypeToString(o.Type),
		}

		ActiveOrders = append(ActiveOrders, order)
	}

	return ActiveOrders
}

//CancelOrderBinance -
func CancelOrderBinance(orderid int64) {

	client := binance.NewClient(apiKey, secretKey)

	_, err := client.NewCancelOrderService().Symbol(MITHBNB).
		OrderID(orderid).Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("Order Cancelled!")
}
