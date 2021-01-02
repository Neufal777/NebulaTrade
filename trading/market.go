package trading

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/ttacon/chalk"

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
	buyActualPrice := exchanges.BinancePrice(exchanges.MITHBNB)
	lastSellFloat := wallet.GetLastSell()
	lastPriceFloat := utils.StringToFloat(buyActualPrice.Price)

	/*
		difference betwen last sell and actual
		to check if the price dropped
	*/
	difference := lastSellFloat - lastPriceFloat
	fmt.Println(chalk.Bold.TextStyle("Waiting for price drop to buy..."), chalk.Red)

	/*
		Showing information about the account
		 - Ammount of crypto
		 - Actual Balance
		 - Transaction...
	*/
	console.InformationDisplayConsole()

	//RandomProfitPerTrans(0.000001, 0.000002)
	if difference >= 0.000001 {

		/*
			EXECUTE BUY ORDER
		*/

		currentWallet := exchanges.GetBinanceWalletBNB() //- 0.01646
		ammountToBuy := currentWallet / lastPriceFloat

		truncateAmmountToBuy := mathnebula.ToFixed((ammountToBuy), 7)

		ammountString := utils.FloatToString(truncateAmmountToBuy)

		//Execute Buy
		exchanges.ExecuteBuyOrderMITHBNB(ammountString[:len(ammountString)-13], buyActualPrice.Price, w)

		//we reset the counter
		w.Timer = 0
		w.WriteInWallet()
		/*
			change last sell file with updated info
			with the lastPriceFloat
		*/

		/*
			- Displaying information
			- Details about the *wallet.Wallet
		*/

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
	data := exchanges.BinancePrice(exchanges.MITHBNB)

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

	//RandomProfitPerTrans(0.000001, 0.000002)
	if differenceToSell >= 0.000001 {

		/*
			EXECUTE SELL ORDER
		*/

		/*
			Change 2 files:
				- Last Sell
				- Ststus to BUY
		*/
		w.Ammount = w.Ammount - 2
		truncatedAmmountToSell := mathnebula.ToFixed((w.Ammount), 7)
		ammountStringSell := utils.FloatToString(truncatedAmmountToSell)
		w.Timer = 0
		exchanges.ExecuteSellOrderMITHBNB(ammountStringSell[:len(ammountStringSell)-13], data.Price, &w)

		w.WriteInWallet()

		fmt.Println(chalk.Bold.TextStyle("SOLD!"), chalk.Green)

	}

}

//ExecuteMarket - workflow of the program
func ExecuteMarket(w *wallet.Wallet) {

	/*
		Check the status (BUY OR SELL)
	*/

	if w.Timer >= 2400 {

		//If we didnt bought anything in X time buy at current price

		orders, allOrders := exchanges.CheckOpenOrdersBinance()

		if orders == 1 {

			for _, ot := range allOrders {

				orderType := utils.AnyTypeToString(ot.Side)
				if orderType == "BUY" {
					exchanges.CancelOrderBinance(ot.OrderID)
					log.Println("Deleted Buy order")
				}
			}

		}

		w.Status = "BUY"
		w.LastSell = 2.9
		w.Timer = 0
		w.WriteInWallet()

	}

	actualStatusString := wallet.GetStatus()
	opened, allorders := exchanges.CheckOpenOrdersBinance()

	switch actualStatusString {
	case "SELL ORDER":
		if opened == 0 {
			w.Status = "BUY"
			w.WriteInWallet()
			DecisionMakeBuy(w)
		} else {

			for _, o := range allorders {

				/*
					Check if the opened orders are sell or buy
				*/

				orderType := utils.AnyTypeToString(o.Side)

				if orderType != "SELL" {

					w.Status = "BUY"
					w.WriteInWallet()
					DecisionMakeBuy(w)
				}
			}
		}
	case "BUY ORDER":
		if opened == 0 {
			w.Status = "SELL"
			w.Timer = 0
			w.WriteInWallet()
			DecisionMakeSell()

		} else {
			for _, o := range allorders {

				/*
					Check if the opened orders are sell or buy
				*/

				orderType := utils.AnyTypeToString(o.Side)

				if orderType != "BUY" {

					w.Status = "SELL"
					w.WriteInWallet()
					DecisionMakeBuy(w)
				}
			}
		}
	case "BUY":
		DecisionMakeBuy(w)
	case "SELL":
		DecisionMakeSell()

	default:
		log.Println("Waiting to close an order..")
	}

}
