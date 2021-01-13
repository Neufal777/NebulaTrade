package trading

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/ttacon/chalk"

	"github.com/NebulaTrade/config"
	"github.com/NebulaTrade/console"
	"github.com/NebulaTrade/exchanges"
	"github.com/NebulaTrade/mathnebula"
	"github.com/NebulaTrade/utils"
	"github.com/NebulaTrade/wallet"
)

//RandomProfitPerTrans -
func RandomProfitPerTrans(min, max float64) float64 {

	x := (rand.Float64() * (max - min)) + min
	cutFloat := mathnebula.ToFixed((x), 7)
	fmt.Println(cutFloat)

	return cutFloat
}

//DecisionMakeBuy - where the decisions of buying or selling is made
func DecisionMakeBuy(w *wallet.Wallet) {

	//Check at how much we sold and if the actual price is lower
	buyActualPrice := exchanges.BinancePrice(config.CURRENCY)
	lastSellFloat := wallet.GetLastSell()
	//lastBuyfloat := w.LastBuy
	lastPriceFloat := utils.StringToFloat(buyActualPrice.Price)

	/*
		difference betwen last sell and actual
		to check if the price dropped
	*/
	difference := lastSellFloat - lastPriceFloat

	fmt.Println(chalk.Bold.TextStyle("Waiting for price drop to buy..."), chalk.Red)
	console.InformationDisplayConsole()

	if difference >= config.PROFIT {

		/*
			EXECUTE BUY ORDER
		*/

		ammountToBuy := w.Available / lastPriceFloat
		truncateAmmountToBuy := mathnebula.ToFixed((ammountToBuy), 7)
		ammountString := utils.FloatToString(truncateAmmountToBuy)

		exchanges.ExecuteBuyOrderCURRENCY(ammountString[:len(ammountString)-13], utils.AnyTypeToString(lastPriceFloat), w)

		w.Timer = 0
		w.FirstBuy = lastPriceFloat
		w.Status = "BUY MORE"
		w.WriteInWallet()

		console.InformationDisplayConsole()

	}
}

//DecisionMakeSell - once we've bought new crypto, we wait to sell them
func DecisionMakeSell() {

	/*
		Get latest data from the coin
		To evaluate when to sell the coin
		and change status from buying to selling
	*/

	w := wallet.ReadWallet()
	data := exchanges.BinancePrice(config.CURRENCY)

	currentPriceFloat, _ := strconv.ParseFloat(data.Price, 32)
	/*
		Information about the last BUY
	*/

	differenceToSell := currentPriceFloat - w.LastBuy

	/*
		Displaying information in the console
	*/
	fmt.Println(chalk.Bold.TextStyle("Waiting to sell.."), chalk.Green)
	console.InformationDisplayConsole()

	if differenceToSell >= config.PROFIT {

		/*
			EXECUTE SELL ORDER
		*/

		/*
			Change 2 files:
				- Last Sell
				- Ststus to BUY
		*/
		truncatedAmmountToSell := mathnebula.ToFixed(
			+(exchanges.GetBinanceWalletCurrency(w.Symbol)), 7)
		ammountStringSell := utils.FloatToString(truncatedAmmountToSell)
		w.Timer = 0
		exchanges.ExecuteSellOrderCURRENCY(ammountStringSell[:len(ammountStringSell)-13], data.Price, &w)

		w.WriteInWallet()

		fmt.Println(chalk.Bold.TextStyle("SOLD!"), chalk.Green)

	}

}

//ExecuteMarket - workflow of the program
func ExecuteMarket(w *wallet.Wallet) {

	/*
		Check the status (BUY OR SELL)
	*/

	if w.Timer >= config.COUNTER {

		//If we didnt bought anything in X time buy at current price

		orders, allOrders := exchanges.CheckOpenOrdersBinance()
		ord := exchanges.AllOpenOrdersBinance(allOrders)

		if orders == 1 {

			for _, o := range ord {

				if o.Status == "BUY" && o.Symbol == config.CURRENCY {

					//We delete the actual buy order [EXPIRED]
					exchanges.CancelOrderBinance(o.OrderID)
					log.Println("Deleted Buy orders")

				} else {
					w.Timer = 0
					w.WriteInWallet()
				}
			}

		} else {
			w.Timer = 0
			w.WriteInWallet()
		}

		w.Status = "BUY"
		w.LastSell = 200.0 //
		w.Timer = 0
		w.OrdNum = 0
		w.WriteInWallet()

	}

	actualStatusString := wallet.GetStatus()

	switch actualStatusString {

	case "BUY":

		//Check if there is available funds to buy
		availablebnb := exchanges.GetBinanceWalletBNB()

		if availablebnb >= 0.02 {

			w.WriteInWallet()
			DecisionMakeBuy(w)
		} else {
			log.Println("No funds available")
		}
	case "SELL":

		log.Println("Ready to sell..")
		w.Timer = 0
		w.WriteInWallet()
		SellingPositions()

	case "BUY MORE":
		SellingPositions()
		RecurrentBuy()
	default:
		log.Println("Waiting to close an order..")
	}

}
