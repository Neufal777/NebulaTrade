package exchanges

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/NebulaTrade/binanceaccount"
	"github.com/NebulaTrade/config"
	"github.com/NebulaTrade/utils"
	"github.com/NebulaTrade/wallet"
	"github.com/adshao/go-binance/v2"
)

var (
	apiKey    = binanceaccount.APIKEY
	secretKey = binanceaccount.APISECRET
	client    = binance.NewClient(apiKey, secretKey)
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

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)

	}
	binanceWalletFloat := utils.StringToFloat(res.Balances[4].Free)

	return binanceWalletFloat
}

//GetBinanceWalletCurrency - get current ammount of that currency
func GetBinanceWalletCurrency(currency string) float64 {

	//clean currency
	currency = strings.Replace(currency, "BNB", "", -1)
	log.Println(currency)

	res, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Println(err)

	}

	count := 0

	for _, x := range res.Balances {

		if x.Asset == currency {
			binanceWalletFloat := utils.StringToFloat(res.Balances[count].Free)
			return binanceWalletFloat
		} else {
			count++
		}
	}

	return 0
}

//ExecuteBuyOrderCURRENCY -
func ExecuteBuyOrderCURRENCY(ammountToBuy string, priceToBuy string, w *wallet.Wallet) {

	resp, err := client.NewCreateOrderService().Symbol(config.CURRENCY).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(ammountToBuy).
		Price(priceToBuy).Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("Order: ", resp)

	//Add the order in the order list, update wallet
	w.Orders = append(w.Orders, wallet.Order{
		Symbol:        resp.Symbol,
		OrderID:       resp.OrderID,
		ClientOrderID: resp.ClientOrderID,
		Price:         resp.Price,
		Ammount:       resp.OrigQuantity,
		Status:        utils.AnyTypeToString(resp.Side),
		Type:          utils.AnyTypeToString(resp.Type),
		Active:        1, //not sold yet
	})

	w.OrdNum++

	if w.OrdNum >= 3 {
		w.Status = "SELL"
	}

	w.Transactions++
	w.LastBuy = utils.StringToFloat(priceToBuy)

	//Record changes
	w.WriteInWallet()

}

//ExecuteSellOrderCURRENCY -
func ExecuteSellOrderCURRENCY(ammountToSell string, priceToSell string, w *wallet.Wallet) {

	_, err := client.NewCreateOrderService().Symbol(config.CURRENCY).
		Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(ammountToSell).
		Price(priceToSell).Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	w.LastSell = utils.StringToFloat(priceToSell)
	w.WriteInWallet()

	// openOrders, allOrders := CheckOpenOrdersBinance()

	// switch openOrders {
	// case 0:

	// 	/*
	// 		If the order is made successfuly update the wallet
	// 	*/

	// 	w.Status = "BUY"
	// 	w.Balance = w.Available
	// 	w.Transactions++
	// 	w.Ammount = 0
	// 	w.LastSell = utils.StringToFloat(priceToSell)
	// 	w.Timer = 0
	// 	w.WriteInWallet()

	// 	fmt.Println(chalk.Bold.TextStyle("SELL ORDER EXECUTED!"), chalk.Green)

	// case 1:

	// 	orders := AllOpenOrdersBinance(allOrders)

	// 	for _, ord := range orders {

	// 		if ord.Status == "SELL" && ord.Symbol == config.CURRENCY {

	// 			//That means we still have open order for selling
	// 			w.Status = "SELL ORDER"
	// 			w.Balance = w.Available
	// 			w.LastSell = utils.StringToFloat(priceToSell)
	// 			w.Ammount = GetBinanceWalletCurrency(w.Symbol)
	// 			w.Timer = 0
	// 			w.Transactions++
	// 			w.WriteInWallet()

	// 		} else {

	// 			//If there is no open orders for specific Coin and status
	// 			w.Status = "BUY"
	// 			w.LastSell = utils.StringToFloat(priceToSell)
	// 			w.Ammount = 0
	// 			w.Balance = w.Available
	// 			w.Timer = 0
	// 			w.Transactions++

	// 			//Write the update information in our register
	// 			w.WriteInWallet()

	// 			fmt.Println(chalk.Bold.TextStyle("SELL ORDER EXECUTED!"), chalk.Green)
	// 		}
	// 	}

	// }
}

//CheckOpenOrdersBinance -
func CheckOpenOrdersBinance() (int, []*binance.Order) {

	//check open orders from the user
	openOrders, err := client.NewListOpenOrdersService().Symbol(config.CURRENCY).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)

	}

	var ordersID []int64
	var orderTypes []string

	for _, o := range openOrders {

		//Show active orders
		log.Println("-----------------------")
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
func AllOpenOrdersBinance(allOrders []*binance.Order) []wallet.Order {

	var ActiveOrders []wallet.Order

	for _, o := range allOrders {

		order := wallet.Order{
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

	_, err := client.NewCancelOrderService().Symbol(config.CURRENCY).
		OrderID(orderid).Do(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}

	log.Println("Order Cancelled!")
}
